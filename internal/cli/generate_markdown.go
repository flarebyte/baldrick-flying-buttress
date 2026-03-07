package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

type generateMarkdownAction struct {
	out io.Writer
}

func (a generateMarkdownAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = a.out
	_ = report
	return writeMarkdownReports(validated)
}

func (generateMarkdownAction) AllowOnValidationErrors() bool {
	return false
}

func writeMarkdownReports(app domain.ValidatedApp) error {
	noteByID := map[string]domain.Note{}
	for _, note := range ordering.Notes(app.Notes) {
		noteByID[note.ID] = note
	}

	for _, report := range ordering.MarkdownReports(app.MarkdownReports) {
		destination := filepath.Join(app.ConfigDir, report.Filepath)
		content := renderMarkdownReport(report, noteByID)
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return fmt.Errorf("create report directory %s: %w", filepath.Dir(destination), err)
		}
		if err := os.WriteFile(destination, []byte(content), 0o644); err != nil {
			return fmt.Errorf("write report %s: %w", destination, err)
		}
	}
	return nil
}

func renderMarkdownReport(report domain.MarkdownReport, noteByID map[string]domain.Note) string {
	var b strings.Builder
	writeMarkdownHeading(&b, 1, report.Title)
	writeMarkdownParagraph(&b, report.Description)

	for _, h2 := range ordering.MarkdownH2Sections(report.Sections) {
		writeMarkdownHeading(&b, 2, h2.Title)
		writeMarkdownParagraph(&b, h2.Description)
		for _, h3 := range ordering.MarkdownH3Sections(h2.Sections) {
			writeMarkdownHeading(&b, 3, h3.Title)
			writeMarkdownParagraph(&b, h3.Description)
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
	return b.String()
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
