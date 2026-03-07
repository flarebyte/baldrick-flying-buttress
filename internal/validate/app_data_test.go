package validate

import (
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestAppDataValidatorRunsSubstepsInOrder(t *testing.T) {
	t.Parallel()

	steps := make([]string, 0, 4)
	validator := AppDataValidator{stepHook: func(step string) {
		steps = append(steps, step)
	}}

	_, _, err := validator.Validate(domain.RawApp{Source: "app"})
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	want := []string{
		"raw_model_normalization_precheck",
		"schema_structure_validation_placeholder",
		"diagnostics_collection",
		"validated_app_normalization",
	}
	if len(steps) != len(want) {
		t.Fatalf("step count mismatch: got %d want %d", len(steps), len(want))
	}
	for i := range want {
		if steps[i] != want[i] {
			t.Fatalf("step %d mismatch: got %q want %q", i, steps[i], want[i])
		}
	}
}

func TestAppDataValidatorCollectsDiagnosticsDeterministically(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source:        "app",
		Reports:       []domain.RawReport{{ID: "", Title: ""}},
		Notes:         []domain.RawNote{{ID: "", Label: ""}},
		Relationships: []domain.RawRelationship{{FromID: "", ToID: "", Label: ""}},
	}

	validator := AppDataValidator{}
	_, report1, err1 := validator.Validate(raw)
	if err1 != nil {
		t.Fatalf("first validate failed: %v", err1)
	}
	_, report2, err2 := validator.Validate(raw)
	if err2 != nil {
		t.Fatalf("second validate failed: %v", err2)
	}

	if len(report1.Diagnostics) != len(report2.Diagnostics) {
		t.Fatalf("diagnostic length mismatch: %d vs %d", len(report1.Diagnostics), len(report2.Diagnostics))
	}
	for i := range report1.Diagnostics {
		if report1.Diagnostics[i] != report2.Diagnostics[i] {
			t.Fatalf("diagnostic %d mismatch: %#v vs %#v", i, report1.Diagnostics[i], report2.Diagnostics[i])
		}
	}
}

func TestAppDataValidatorReturnsNormalizedValidatedAppWithDiagnostics(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source:        "app",
		Reports:       []domain.RawReport{{ID: "z", Title: "Z"}, {ID: "", Title: "Missing"}, {ID: "a", Title: "A"}},
		Notes:         []domain.RawNote{{ID: "n2", Label: "service.db"}, {ID: "", Label: ""}, {ID: "n1", Label: "service.api"}},
		Relationships: []domain.RawRelationship{{FromID: "n1", ToID: "n2", Label: "owns"}, {FromID: "", ToID: "", Label: ""}, {FromID: "n1", ToID: "n2", Label: "depends_on"}},
	}

	validator := AppDataValidator{}
	app, report, err := validator.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if len(report.Diagnostics) == 0 {
		t.Fatal("expected diagnostics")
	}
	if len(app.Reports) != 3 {
		t.Fatalf("expected reports to be returned, got %d", len(app.Reports))
	}
	if app.Reports[0].ID != "" {
		t.Fatalf("expected canonical sorted reports with empty id first, got %#v", app.Reports)
	}
	if len(app.Notes) != 3 {
		t.Fatalf("expected notes to be returned, got %d", len(app.Notes))
	}
	if len(app.Relationships) != 3 {
		t.Fatalf("expected relationships to be returned, got %d", len(app.Relationships))
	}
}

func TestAppDataValidatorDeterministicAcrossRuns(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source: "app",
		Name:   "stub-app",
		Reports: []domain.RawReport{
			{ID: "memory-health", Title: "Memory Health"},
			{ID: "cpu-overview", Title: "CPU Overview"},
		},
		Notes: []domain.RawNote{
			{ID: "n2", Label: "service.db"},
			{ID: "n1", Label: "service.api"},
		},
		Relationships: []domain.RawRelationship{
			{FromID: "n1", ToID: "n2", Label: "owns"},
			{FromID: "n1", ToID: "n2", Label: "depends_on"},
		},
	}

	validator := AppDataValidator{}
	app1, report1, err1 := validator.Validate(raw)
	if err1 != nil {
		t.Fatalf("first validate failed: %v", err1)
	}
	app2, report2, err2 := validator.Validate(raw)
	if err2 != nil {
		t.Fatalf("second validate failed: %v", err2)
	}

	if !reflect.DeepEqual(app1, app2) {
		t.Fatalf("non-deterministic app: %#v vs %#v", app1, app2)
	}
	if len(report1.Diagnostics) != len(report2.Diagnostics) {
		t.Fatalf("non-deterministic diagnostics length: %d vs %d", len(report1.Diagnostics), len(report2.Diagnostics))
	}
	for i := range report1.Diagnostics {
		if report1.Diagnostics[i] != report2.Diagnostics[i] {
			t.Fatalf("non-deterministic diagnostic at %d: %#v vs %#v", i, report1.Diagnostics[i], report2.Diagnostics[i])
		}
	}
}
