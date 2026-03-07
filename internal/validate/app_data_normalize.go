package validate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

func collectDiagnostics(diagnostics []domain.Diagnostic) []domain.Diagnostic {
	if diagnostics == nil {
		diagnostics = []domain.Diagnostic{}
	}
	return ordering.Diagnostics(diagnostics)
}

func normalizeValidatedApp(raw domain.RawApp, registry domain.ArgumentRegistry, datasetLabels []string, graphPolicy domain.GraphIntegrityPolicy) domain.ValidatedApp {
	reports := make([]domain.Report, 0, len(raw.Reports))
	for _, report := range raw.Reports {
		reports = append(reports, domain.Report{
			ID:    reportIDFromFilepath(report.Filepath),
			Title: report.Title,
		})
	}

	notes := make([]domain.Note, 0, len(raw.Notes))
	for _, note := range raw.Notes {
		notes = append(notes, domain.Note{
			ID:           note.Name,
			Label:        note.Title,
			Title:        note.Title,
			Markdown:     note.Markdown,
			Filepath:     strings.TrimSpace(note.Filepath),
			LabelsCSV:    strings.Join(ordering.Strings(note.Labels), ","),
			ArgumentsCSV: strings.Join(ordering.Strings(note.Arguments), "\n"),
		})
	}

	relationships := make([]domain.Relationship, 0, len(raw.Relationships))
	for _, relationship := range raw.Relationships {
		labels := append([]string{}, relationship.Labels...)
		if relationship.Label != "" {
			labels = append(labels, relationship.Label)
		}
		relationships = append(relationships, domain.Relationship{
			FromID:    relationship.FromID,
			ToID:      relationship.ToID,
			Label:     relationship.Label,
			LabelsCSV: strings.Join(ordering.Strings(labels), ","),
		})
	}

	configDir := "."
	if raw.ConfigPath != "" {
		configDir = filepath.Dir(raw.ConfigPath)
	}

	return domain.ValidatedApp{
		Name:                 raw.Name,
		ConfigDir:            configDir,
		Modules:              raw.Modules,
		Reports:              ordering.Reports(reports),
		MarkdownReports:      normalizeMarkdownReports(raw),
		Notes:                ordering.Notes(notes),
		Relationships:        ordering.Relationships(relationships),
		Registry:             registry,
		DatasetLabels:        datasetLabels,
		GraphIntegrityPolicy: graphPolicy,
	}
}

func normalizeMarkdownReports(raw domain.RawApp) []domain.MarkdownReport {
	reports := make([]domain.MarkdownReport, 0, len(raw.Reports))
	for reportIndex, rawReport := range raw.Reports {
		report := domain.MarkdownReport{
			Title:       rawReport.Title,
			Filepath:    rawReport.Filepath,
			Description: rawReport.Description,
			Sections:    make([]domain.MarkdownH2Section, 0, len(rawReport.Sections)),
		}
		for h2Index, rawH2 := range rawReport.Sections {
			h2 := domain.MarkdownH2Section{
				Title:       rawH2.Title,
				Description: rawH2.Description,
				Sections:    make([]domain.MarkdownH3Section, 0, len(rawH2.Sections)),
			}
			for h3Index, rawH3 := range rawH2.Sections {
				h2.Sections = append(h2.Sections, domain.MarkdownH3Section{
					Title:       rawH3.Title,
					Description: rawH3.Description,
					NoteIDs:     ordering.Strings(rawH3.Notes),
					Arguments:   ordering.Strings(rawH3.Arguments),
					Path:        fmt.Sprintf("reports[%d].sections[%d].sections[%d]", reportIndex, h2Index, h3Index),
					H2Title:     rawH2.Title,
					ReportTitle: rawReport.Title,
				})
			}
			h2.Sections = ordering.MarkdownH3Sections(h2.Sections)
			report.Sections = append(report.Sections, h2)
		}
		report.Sections = ordering.MarkdownH2Sections(report.Sections)
		reports = append(reports, report)
	}
	return ordering.MarkdownReports(reports)
}
