package validate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

const schemaValidationSource = "validate.app.data.schema"

type AppDataValidator struct {
	stepHook func(string)
}

func (v AppDataValidator) Validate(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
	v.step("raw_model_normalization_precheck")
	rawModel := normalizeRaw(raw)

	v.step("schema_structure_validation_placeholder")
	diagnostics := validateStructure(rawModel)

	v.step("diagnostics_collection")
	diagnostics = collectDiagnostics(diagnostics)

	v.step("validated_app_normalization")
	validated := normalizeValidatedApp(rawModel)

	return validated, domain.ValidationReport{Diagnostics: diagnostics}.Canonical(), nil
}

func (v AppDataValidator) step(name string) {
	if v.stepHook != nil {
		v.stepHook(name)
	}
}

func normalizeRaw(raw domain.RawApp) domain.RawApp {
	if raw.Modules == nil {
		raw.Modules = []string{}
	}
	return raw
}

func validateStructure(raw domain.RawApp) []domain.Diagnostic {
	diagnostics := make([]domain.Diagnostic, 0)

	if raw.Source == "" {
		diagnostics = append(diagnostics, newDiagnostic(
			"FBV000",
			"missing required field: source",
			"source",
		))
	}

	if raw.Reports == nil {
		diagnostics = append(diagnostics, newDiagnostic(
			"FBV100",
			"missing required collection: reports",
			"reports",
		))
	}
	if raw.Notes == nil {
		diagnostics = append(diagnostics, newDiagnostic(
			"FBV200",
			"missing required collection: notes",
			"notes",
		))
	}
	if raw.Relationships == nil {
		diagnostics = append(diagnostics, newDiagnostic(
			"FBV300",
			"missing required collection: relationships",
			"relationships",
		))
	}

	for i, report := range raw.Reports {
		reportCtx := report.Title
		if report.Title == "" {
			diagnostics = append(diagnostics, newDiagnosticWithContext(
				"FBV101",
				"missing required field: report title",
				reportLocation(i, "title"),
				reportCtx,
				"",
				"",
			))
		}
		if report.Filepath == "" {
			diagnostics = append(diagnostics, newDiagnosticWithContext(
				"FBV102",
				"missing required field: report filepath",
				reportLocation(i, "filepath"),
				reportCtx,
				"",
				"",
			))
		}
		if report.Sections == nil {
			diagnostics = append(diagnostics, newDiagnosticWithContext(
				"FBV103",
				"missing required field: report sections",
				reportLocation(i, "sections"),
				reportCtx,
				"",
				"",
			))
		}
		for j, section := range report.Sections {
			if strings.TrimSpace(section.Title) == "" {
				diagnostics = append(diagnostics, newDiagnosticWithContext(
					"FBV104",
					"missing required field: section title",
					reportSectionLocation(i, j, "title"),
					reportCtx,
					section.Title,
					"",
				))
			}
		}
	}

	for i, note := range raw.Notes {
		if note.Name == "" {
			diagnostics = append(diagnostics, newDiagnosticWithContext(
				"FBV201",
				"missing required field: note name",
				noteLocation(i, "name"),
				"",
				"",
				note.Name,
			))
		}
		if note.Title == "" {
			diagnostics = append(diagnostics, newDiagnosticWithContext(
				"FBV202",
				"missing required field: note title",
				noteLocation(i, "title"),
				"",
				"",
				note.Name,
			))
		}
	}

	for i, relationship := range raw.Relationships {
		if relationship.FromID == "" {
			diagnostics = append(diagnostics, newDiagnostic(
				"FBV301",
				"missing required field: relationship from",
				relationshipLocation(i, "from"),
			))
		}
		if relationship.ToID == "" {
			diagnostics = append(diagnostics, newDiagnostic(
				"FBV302",
				"missing required field: relationship to",
				relationshipLocation(i, "to"),
			))
		}
	}

	return diagnostics
}

func collectDiagnostics(diagnostics []domain.Diagnostic) []domain.Diagnostic {
	if diagnostics == nil {
		diagnostics = []domain.Diagnostic{}
	}
	return ordering.Diagnostics(diagnostics)
}

func normalizeValidatedApp(raw domain.RawApp) domain.ValidatedApp {
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
			ID:    note.Name,
			Label: note.Title,
		})
	}

	relationships := make([]domain.Relationship, 0, len(raw.Relationships))
	for _, relationship := range raw.Relationships {
		relationships = append(relationships, domain.Relationship(relationship))
	}

	return domain.ValidatedApp{
		Name:          raw.Name,
		Modules:       raw.Modules,
		Reports:       ordering.Reports(reports),
		Notes:         ordering.Notes(notes),
		Relationships: ordering.Relationships(relationships),
	}
}

func reportIDFromFilepath(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	id := strings.TrimSuffix(base, ext)
	return id
}

func newDiagnostic(code, message, location string) domain.Diagnostic {
	return domain.Diagnostic{
		Code:     code,
		Severity: domain.SeverityError,
		Source:   schemaValidationSource,
		Message:  message,
		Location: location,
		Path:     location,
	}
}

func newDiagnosticWithContext(code, message, location, reportTitle, sectionTitle, noteName string) domain.Diagnostic {
	d := newDiagnostic(code, message, location)
	d.ReportTitle = reportTitle
	d.SectionTitle = sectionTitle
	d.NoteName = noteName
	return d
}

func reportLocation(i int, field string) string {
	return fmt.Sprintf("reports[%d].%s", i, field)
}

func reportSectionLocation(reportIndex, sectionIndex int, field string) string {
	return fmt.Sprintf("reports[%d].sections[%d].%s", reportIndex, sectionIndex, field)
}

func noteLocation(i int, field string) string {
	return fmt.Sprintf("notes[%d].%s", i, field)
}

func relationshipLocation(i int, field string) string {
	return fmt.Sprintf("relationships[%d].%s", i, field)
}
