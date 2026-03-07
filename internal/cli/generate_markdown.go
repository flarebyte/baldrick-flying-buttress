package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	clioutput "github.com/flarebyte/baldrick-flying-buttress/internal/output"
	"github.com/flarebyte/baldrick-flying-buttress/internal/renderer"
)

type generateMarkdownAction struct {
	out io.Writer
}

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

	for _, report := range ordering.MarkdownReports(app.MarkdownReports) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		destination := filepath.Join(app.ConfigDir, report.Filepath)
		content, sectionDiagnostics, err := renderMarkdownReport(ctx, report, noteByID, app, registry)
		if err != nil {
			return nil, err
		}
		diagnostics = append(diagnostics, sectionDiagnostics...)
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return nil, fmt.Errorf("create report directory %s: %w", filepath.Dir(destination), err)
		}
		if err := writeFileAtomically(ctx, destination, []byte(content), 0o644); err != nil {
			return nil, fmt.Errorf("write report %s: %w", destination, err)
		}
	}
	return ordering.Diagnostics(diagnostics), nil
}

func writeFileAtomically(ctx context.Context, destination string, data []byte, perm os.FileMode) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	tmp := destination + ".tmp"
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, destination)
}
