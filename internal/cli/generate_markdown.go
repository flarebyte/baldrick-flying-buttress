package cli

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/fsio"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	clioutput "github.com/flarebyte/baldrick-flying-buttress/internal/output"
	"github.com/flarebyte/baldrick-flying-buttress/internal/renderer"
)

type generateMarkdownAction struct {
	out io.Writer
}

var markdownReportWorkers int32 = int32(defaultMarkdownReportWorkers())

func (a generateMarkdownAction) Execute(ctx context.Context, validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = a.out
	_ = report
	diagnostics, err := writeMarkdownReports(ctx, validated)
	if err != nil {
		return err
	}
	if len(diagnostics) == 0 {
		return nil
	}
	outReport := domain.ValidationReport{Diagnostics: diagnostics}
	if err := clioutput.EmitDiagnostics(a.out, outReport); err != nil {
		return err
	}
	if outReport.HasErrors() {
		return outcome.ValidationBlockedError()
	}
	return nil
}

func (generateMarkdownAction) AllowOnValidationErrors() bool {
	return false
}

func writeMarkdownReports(ctx context.Context, app domain.ValidatedApp) ([]domain.Diagnostic, error) {
	return writeMarkdownReportsWithWorkers(ctx, app, getMarkdownReportWorkers())
}

func writeMarkdownReportsWithWorkers(ctx context.Context, app domain.ValidatedApp, workers int) ([]domain.Diagnostic, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	noteByID := map[string]domain.Note{}
	for _, note := range ordering.Notes(app.Notes) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		noteByID[note.ID] = note
	}
	diagnostics := make([]domain.Diagnostic, 0)
	registry := renderer.ResolveRegistry()
	orderedReports := ordering.MarkdownReports(app.MarkdownReports)
	rendered, renderDiagnostics, err := renderReportsConcurrently(ctx, app, noteByID, registry, orderedReports, normalizeWorkerCount(workers))
	if err != nil {
		return nil, err
	}
	diagnostics = append(diagnostics, renderDiagnostics...)

	for _, output := range rendered {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if err := fsio.WriteFileAtomic(ctx, output.destination, []byte(output.content), 0o644); err != nil {
			return nil, fmt.Errorf("write report %s: %w", output.destination, err)
		}
	}
	return ordering.Diagnostics(diagnostics), nil
}

type markdownReportJob struct {
	index  int
	report domain.MarkdownReport
}

type markdownRenderedReport struct {
	index       int
	destination string
	content     string
	diagnostics []domain.Diagnostic
	err         error
}

func renderReportsConcurrently(
	ctx context.Context,
	app domain.ValidatedApp,
	noteByID map[string]domain.Note,
	registry renderer.Registry,
	reports []domain.MarkdownReport,
	workers int,
) ([]markdownRenderedReport, []domain.Diagnostic, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan markdownReportJob)
	results := make(chan markdownRenderedReport, len(reports))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					destination := filepath.Join(app.ConfigDir, job.report.Filepath)
					content, sectionDiagnostics, err := renderMarkdownReport(ctx, job.report, noteByID, app, registry)
					if err != nil {
						cancel()
						results <- markdownRenderedReport{index: job.index, destination: destination, err: err}
						continue
					}
					results <- markdownRenderedReport{
						index:       job.index,
						destination: destination,
						content:     content,
						diagnostics: sectionDiagnostics,
					}
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for i, report := range reports {
			select {
			case <-ctx.Done():
				return
			case jobs <- markdownReportJob{index: i, report: report}:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	renderedByIndex := make([]markdownRenderedReport, len(reports))
	renderedFlags := make([]bool, len(reports))
	diagsByIndex := make([][]domain.Diagnostic, len(reports))
	errByIndex := make([]error, len(reports))

	for result := range results {
		if result.err != nil {
			if errByIndex[result.index] == nil {
				errByIndex[result.index] = result.err
			}
			continue
		}
		renderedByIndex[result.index] = result
		renderedFlags[result.index] = true
		diagsByIndex[result.index] = append(diagsByIndex[result.index], result.diagnostics...)
	}

	for _, err := range errByIndex {
		if err != nil {
			return nil, nil, err
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	orderedRendered := make([]markdownRenderedReport, 0, len(reports))
	orderedDiagnostics := make([]domain.Diagnostic, 0)
	for i := range reports {
		if !renderedFlags[i] {
			return nil, nil, fmt.Errorf("render report %d: canceled before completion", i)
		}
		orderedRendered = append(orderedRendered, renderedByIndex[i])
		orderedDiagnostics = append(orderedDiagnostics, diagsByIndex[i]...)
	}
	return orderedRendered, orderedDiagnostics, nil
}

func defaultMarkdownReportWorkers() int {
	workers := runtime.GOMAXPROCS(0)
	if workers > 4 {
		workers = 4
	}
	if workers < 1 {
		return 1
	}
	return workers
}

func normalizeWorkerCount(workers int) int {
	if workers < 1 {
		return 1
	}
	return workers
}

func getMarkdownReportWorkers() int {
	return int(atomic.LoadInt32(&markdownReportWorkers))
}
