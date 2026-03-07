package cli

import (
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

func (a generateMarkdownAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = a.out
	_ = report
	diagnostics, err := writeMarkdownReports(validated)
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

func writeMarkdownReports(app domain.ValidatedApp) ([]domain.Diagnostic, error) {
	noteByID := map[string]domain.Note{}
	for _, note := range ordering.Notes(app.Notes) {
		noteByID[note.ID] = note
	}
	diagnostics := make([]domain.Diagnostic, 0)
	registry := renderer.ResolveRegistry()

	for _, report := range ordering.MarkdownReports(app.MarkdownReports) {
		destination := filepath.Join(app.ConfigDir, report.Filepath)
		content, sectionDiagnostics, err := renderMarkdownReport(report, noteByID, app, registry)
		if err != nil {
			return nil, err
		}
		diagnostics = append(diagnostics, sectionDiagnostics...)
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return nil, fmt.Errorf("create report directory %s: %w", filepath.Dir(destination), err)
		}
		if err := os.WriteFile(destination, []byte(content), 0o644); err != nil {
			return nil, fmt.Errorf("write report %s: %w", destination, err)
		}
	}
	return ordering.Diagnostics(diagnostics), nil
}
