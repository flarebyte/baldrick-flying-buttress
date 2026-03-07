package validate

import (
	"fmt"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

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
	if raw.Reports == nil {
		raw.Reports = []domain.RawReport{}
	}
	if raw.Notes == nil {
		raw.Notes = []domain.RawNote{}
	}
	if raw.Relationships == nil {
		raw.Relationships = []domain.RawRelationship{}
	}
	return raw
}

func validateStructure(raw domain.RawApp) []domain.Diagnostic {
	diagnostics := make([]domain.Diagnostic, 0)

	if raw.Source == "" {
		diagnostics = append(diagnostics, domain.Diagnostic{
			Code:     "FBV000",
			Severity: domain.SeverityError,
			Message:  "missing required field: source",
			Path:     "source",
		})
	}

	for i, report := range raw.Reports {
		if report.ID == "" {
			diagnostics = append(diagnostics, domain.Diagnostic{
				Code:     "FBV101",
				Severity: domain.SeverityError,
				Message:  "missing required field: report id",
				Path:     fmt.Sprintf("reports[%d].id", i),
			})
		}
		if report.Title == "" {
			diagnostics = append(diagnostics, domain.Diagnostic{
				Code:     "FBV102",
				Severity: domain.SeverityError,
				Message:  "missing required field: report title",
				Path:     fmt.Sprintf("reports[%d].title", i),
			})
		}
	}

	for i, note := range raw.Notes {
		if note.ID == "" {
			diagnostics = append(diagnostics, domain.Diagnostic{
				Code:     "FBV201",
				Severity: domain.SeverityError,
				Message:  "missing required field: note id",
				Path:     fmt.Sprintf("notes[%d].id", i),
			})
		}
		if note.Label == "" {
			diagnostics = append(diagnostics, domain.Diagnostic{
				Code:     "FBV202",
				Severity: domain.SeverityError,
				Message:  "missing required field: note label",
				Path:     fmt.Sprintf("notes[%d].label", i),
			})
		}
	}

	for i, relationship := range raw.Relationships {
		if relationship.FromID == "" {
			diagnostics = append(diagnostics, domain.Diagnostic{
				Code:     "FBV301",
				Severity: domain.SeverityError,
				Message:  "missing required field: relationship from",
				Path:     fmt.Sprintf("relationships[%d].from", i),
			})
		}
		if relationship.ToID == "" {
			diagnostics = append(diagnostics, domain.Diagnostic{
				Code:     "FBV302",
				Severity: domain.SeverityError,
				Message:  "missing required field: relationship to",
				Path:     fmt.Sprintf("relationships[%d].to", i),
			})
		}
		if relationship.Label == "" {
			diagnostics = append(diagnostics, domain.Diagnostic{
				Code:     "FBV303",
				Severity: domain.SeverityError,
				Message:  "missing required field: relationship label",
				Path:     fmt.Sprintf("relationships[%d].label", i),
			})
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
		reports = append(reports, domain.Report(report))
	}

	notes := make([]domain.Note, 0, len(raw.Notes))
	for _, note := range raw.Notes {
		notes = append(notes, domain.Note(note))
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
