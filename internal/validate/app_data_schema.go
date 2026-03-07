package validate

import (
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func validateStructure(raw domain.RawApp) []domain.Diagnostic {
	diagnostics := make([]domain.Diagnostic, 0)

	if raw.Source == "" {
		diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV000", "missing required field: source", "source"))
	}

	if raw.Reports == nil {
		diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV100", "missing required collection: reports", "reports"))
	}
	if raw.Notes == nil {
		diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV200", "missing required collection: notes", "notes"))
	}
	if raw.Relationships == nil {
		diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV300", "missing required collection: relationships", "relationships"))
	}

	for i, report := range raw.Reports {
		if report.Title == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV101", "missing required field: report title", reportLocation(i, "title")))
		}
		if report.Filepath == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV102", "missing required field: report filepath", reportLocation(i, "filepath")))
		}
		if report.Sections == nil {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV103", "missing required field: report sections", reportLocation(i, "sections")))
		}
		for j, section := range report.Sections {
			if strings.TrimSpace(section.Title) == "" {
				diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV104", "missing required field: section title", reportSectionLocation(i, j, "title")))
			}
		}
	}

	for i, note := range raw.Notes {
		if note.Name == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV201", "missing required field: note name", noteLocation(i, "name")))
		}
		if note.Title == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV202", "missing required field: note title", noteLocation(i, "title")))
		}
	}

	for i, relationship := range raw.Relationships {
		if relationship.FromID == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV301", "missing required field: relationship from", relationshipLocation(i, "from")))
		}
		if relationship.ToID == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV302", "missing required field: relationship to", relationshipLocation(i, "to")))
		}
	}

	return diagnostics
}
