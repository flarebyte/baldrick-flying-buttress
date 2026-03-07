package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/graph"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	clioutput "github.com/flarebyte/baldrick-flying-buttress/internal/output"
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

	for _, report := range ordering.MarkdownReports(app.MarkdownReports) {
		destination := filepath.Join(app.ConfigDir, report.Filepath)
		content, sectionDiagnostics := renderMarkdownReport(report, noteByID, app)
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

func renderMarkdownReport(report domain.MarkdownReport, noteByID map[string]domain.Note, app domain.ValidatedApp) (string, []domain.Diagnostic) {
	var b strings.Builder
	diagnostics := make([]domain.Diagnostic, 0)
	writeMarkdownHeading(&b, 1, report.Title)
	writeMarkdownParagraph(&b, report.Description)

	for _, h2 := range ordering.MarkdownH2Sections(report.Sections) {
		writeMarkdownHeading(&b, 2, h2.Title)
		writeMarkdownParagraph(&b, h2.Description)
		for _, h3 := range ordering.MarkdownH3Sections(h2.Sections) {
			writeMarkdownHeading(&b, 3, h3.Title)
			writeMarkdownParagraph(&b, h3.Description)
			if graph.HasGraphArgs(h3.Arguments) {
				query := graph.QueryFromArgs(h3.Arguments)
				selected := graph.Select(app, query)
				shape := graph.DetectShape(selected)
				cyclePolicy, err := graph.ResolveCyclePolicy(h3.Arguments)
				if err != nil {
					diagnostics = append(diagnostics, domain.Diagnostic{
						Code:         "GRAPH_CYCLE_POLICY_INVALID",
						Severity:     domain.SeverityError,
						Source:       "graph.policy.cycle",
						Message:      err.Error(),
						Location:     h3.Path,
						Path:         h3.Path,
						ReportTitle:  h3.ReportTitle,
						SectionTitle: h3.H2Title,
					})
					continue
				}
				if shape == graph.ShapeCyclic && cyclePolicy == graph.CyclePolicyDisallow {
					diagnostics = append(diagnostics, domain.Diagnostic{
						Code:         "GRAPH_CYCLE_DISALLOWED",
						Severity:     domain.SeverityWarning,
						Source:       "graph.policy.cycle",
						Message:      "cycle detected while cycle-policy=disallow; skipping graph section",
						Location:     h3.Path,
						Path:         h3.Path,
						ReportTitle:  h3.ReportTitle,
						SectionTitle: h3.H2Title,
					})
					continue
				}
				graphText := graph.RenderMarkdownText(selected, shape, cyclePolicy)
				if graphText != "" {
					b.WriteString(graphText)
					if !strings.HasSuffix(graphText, "\n") {
						b.WriteByte('\n')
					}
					b.WriteByte('\n')
				}
				continue
			}
			for _, noteID := range ordering.Strings(h3.NoteIDs) {
				note, ok := noteByID[noteID]
				if !ok {
					continue
				}
				writeMarkdownHeading(&b, 4, note.Title)
				writeMarkdownParagraph(&b, note.Markdown)
			}
		}
	}

	if b.Len() == 0 || !strings.HasSuffix(b.String(), "\n") {
		b.WriteByte('\n')
	}
	return b.String(), diagnostics
}

func writeMarkdownHeading(b *strings.Builder, level int, title string) {
	if strings.TrimSpace(title) == "" {
		return
	}
	b.WriteString(strings.Repeat("#", level))
	b.WriteByte(' ')
	b.WriteString(title)
	b.WriteString("\n\n")
}

func writeMarkdownParagraph(b *strings.Builder, text string) {
	if strings.TrimSpace(text) == "" {
		return
	}
	b.WriteString(text)
	b.WriteString("\n\n")
}
