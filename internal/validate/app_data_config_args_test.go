package validate

import (
	"context"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestAppDataValidatorConfiguredArgumentsDiagnostics(t *testing.T) {
	t.Parallel()

	raw := rawAppWithMinimalShape(domain.RawArgumentRegistry{Arguments: []domain.RawArgumentDefinition{
		{Name: "fmt", ValueType: "boolean", Scopes: []string{"renderer"}},
		{Name: "mode", ValueType: "enum", Scopes: []string{"h3-section"}, AllowedValues: []string{"a", "b"}},
		{Name: "verbose", ValueType: "boolean", Scopes: []string{"note"}},
	}})
	raw.Reports = []domain.RawReport{{
		Title:    "R",
		Filepath: "reports/r.md",
		Sections: []domain.RawReportSection{{
			Title:     "S",
			Arguments: []string{"unknown=x", "fmt=true", "mode=z", "badarg", "=x", "k="},
		}},
	}}
	raw.Notes = []domain.RawNote{{
		Name:      "n1",
		Title:     "N1",
		Arguments: []string{"fmt=x", "verbose=maybe"},
	}}

	_, report := validateRaw(t, raw)

	want := []diagExpectation{
		{code: "FBC001", path: "reports[0].sections[0].arguments[3]"},
		{code: "FBC001", path: "reports[0].sections[0].arguments[4]"},
		{code: "FBC001", path: "reports[0].sections[0].arguments[5]"},
		{code: "FBC002", path: "reports[0].sections[0].arguments[0]"},
		{code: "FBC003", path: "notes[0].arguments[0]"},
		{code: "FBC004", path: "notes[0].arguments[0]"},
		{code: "FBC004", path: "notes[0].arguments[1]"},
		{code: "FBC004", path: "reports[0].sections[0].arguments[2]"},
	}
	assertHasDiagnostics(t, report.Diagnostics, want, "configured-args")
}

func TestAppDataValidatorConfiguredArgumentsContextFields(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source: "app",
		Reports: []domain.RawReport{{
			Title:    "CPU Overview",
			Filepath: "reports/cpu-overview.md",
			Sections: []domain.RawReportSection{{Title: "Overview", Arguments: []string{"unknown=x"}}},
		}},
		Notes:         []domain.RawNote{{Name: "n1", Title: "Service API", Arguments: []string{"unknown=x"}}},
		Relationships: []domain.RawRelationship{{FromID: "n1", ToID: "n2", Label: "L"}},
		Registry:      domain.RawArgumentRegistry{Arguments: []domain.RawArgumentDefinition{}},
	}

	_, report, err := AppDataValidator{}.Validate(context.Background(), raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	var sectionDiag *domain.Diagnostic
	var noteDiag *domain.Diagnostic
	for i := range report.Diagnostics {
		d := &report.Diagnostics[i]
		if d.Path == "reports[0].sections[0].arguments[0]" {
			sectionDiag = d
		}
		if d.Path == "notes[0].arguments[0]" {
			noteDiag = d
		}
	}
	if sectionDiag == nil || noteDiag == nil {
		t.Fatalf("expected diagnostics for section and note arguments, got %#v", report.Diagnostics)
	}
	if sectionDiag.ArgumentName != "unknown" || sectionDiag.SectionTitle != "Overview" || sectionDiag.ReportTitle != "CPU Overview" {
		t.Fatalf("unexpected section diagnostic context: %#v", sectionDiag)
	}
	if noteDiag.ArgumentName != "unknown" || noteDiag.NoteName != "n1" {
		t.Fatalf("unexpected note diagnostic context: %#v", noteDiag)
	}
}
