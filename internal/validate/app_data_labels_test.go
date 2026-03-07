package validate

import (
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestAppDataValidatorDatasetLabelsCollectedDeterministically(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source:  "app",
		Reports: []domain.RawReport{{Title: "R", Filepath: "reports/r.md", Sections: []domain.RawReportSection{{Title: "S"}}}},
		Notes: []domain.RawNote{
			{Name: "n1", Title: "N1", Labels: []string{"beta", "alpha"}},
			{Name: "n2", Title: "N2", Labels: []string{"alpha", "gamma"}},
		},
		Relationships: []domain.RawRelationship{
			{FromID: "n1", ToID: "n2", Label: "depends_on", Labels: []string{"delta", "beta"}},
		},
	}

	app, _, err := AppDataValidator{}.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	want := []string{"alpha", "beta", "delta", "depends_on", "gamma"}
	if !reflect.DeepEqual(app.DatasetLabels, want) {
		t.Fatalf("dataset labels mismatch: got %#v want %#v", app.DatasetLabels, want)
	}
}

func TestAppDataValidatorLabelReferenceValidation(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source: "app",
		Reports: []domain.RawReport{{
			Title: "R", Filepath: "reports/r.md",
			Sections: []domain.RawReportSection{{Title: "S", Arguments: []string{"orphan-subject-label=known", "orphan-edge-label=missing-a", "orphan-counterpart-label=missing-b"}}},
		}},
		Notes: []domain.RawNote{
			{Name: "n1", Title: "N1", Labels: []string{"known"}, Arguments: []string{"orphan-subject-label=missing-c"}},
		},
		Relationships: []domain.RawRelationship{{FromID: "n1", ToID: "n2", Label: "depends_on", Labels: []string{"known-rel"}}},
		Registry: domain.RawArgumentRegistry{Arguments: []domain.RawArgumentDefinition{
			{Name: "orphan-subject-label", ValueType: "string", Scopes: []string{"h3-section", "note"}},
			{Name: "orphan-edge-label", ValueType: "string", Scopes: []string{"h3-section"}},
			{Name: "orphan-counterpart-label", ValueType: "string", Scopes: []string{"h3-section"}},
		}},
	}

	_, report, err := AppDataValidator{}.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	want := []struct {
		path  string
		value string
	}{
		{path: "notes[0].arguments[0]", value: "missing-c"},
		{path: "reports[0].sections[0].arguments[1]", value: "missing-a"},
		{path: "reports[0].sections[0].arguments[2]", value: "missing-b"},
	}
	for _, item := range want {
		found := false
		for _, d := range report.Diagnostics {
			if d.Code == "LABEL_REF_UNKNOWN" && d.Path == item.path && d.Source == labelReferenceValidationSource && d.LabelValue == item.value && d.Severity == domain.SeverityWarning {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected label reference warning at %s for %s, got %#v", item.path, item.value, report.Diagnostics)
		}
	}
}
