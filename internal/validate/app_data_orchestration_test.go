package validate

import (
	"context"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestAppDataValidatorRunsSubstepsInOrder(t *testing.T) {
	t.Parallel()

	steps := make([]string, 0, 9)
	validator := AppDataValidator{stepHook: func(step string) {
		steps = append(steps, step)
	}}

	_, _, err := validator.Validate(context.Background(), domain.RawApp{Source: "app"})
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	want := []string{
		"raw_model_normalization_precheck",
		"validate_cue_schema",
		"args_registry_resolve",
		"args_registry_validate",
		"args_validate_config",
		"labels_dataset_collect",
		"labels_reference_validate",
		"graph_integrity_policy_resolve",
		"graph_integrity_validate",
		"diagnostics_collection",
		"validated_app_normalization",
	}
	if !reflect.DeepEqual(steps, want) {
		t.Fatalf("step order mismatch: got %#v want %#v", steps, want)
	}
}
func TestAppDataValidatorDeterministicAcrossRuns(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source: "app",
		Name:   "stub-app",
		Reports: []domain.RawReport{
			{Title: "Memory Health", Filepath: "reports/memory-health.md", Sections: []domain.RawReportSection{{Title: "S"}}},
			{Title: "CPU Overview", Filepath: "reports/cpu-overview.md", Sections: []domain.RawReportSection{{Title: "S"}}},
		},
		Notes:         []domain.RawNote{{Name: "n2", Title: "service.db"}, {Name: "n1", Title: "service.api"}},
		Relationships: []domain.RawRelationship{{FromID: "n1", ToID: "n2", Label: "owns"}, {FromID: "n1", ToID: "n2", Label: "depends_on"}},
		Registry:      domain.RawArgumentRegistry{Arguments: []domain.RawArgumentDefinition{{Name: "beta", ValueType: "string", Scopes: []string{"note"}}, {Name: "alpha", ValueType: "string", Scopes: []string{"renderer"}}}},
	}

	app1, report1, err1 := AppDataValidator{}.Validate(context.Background(), raw)
	if err1 != nil {
		t.Fatalf("first validate failed: %v", err1)
	}
	app2, report2, err2 := AppDataValidator{}.Validate(context.Background(), raw)
	if err2 != nil {
		t.Fatalf("second validate failed: %v", err2)
	}
	if !reflect.DeepEqual(app1, app2) {
		t.Fatalf("app mismatch: %#v vs %#v", app1, app2)
	}
	if !reflect.DeepEqual(report1.Diagnostics, report2.Diagnostics) {
		t.Fatalf("diagnostics mismatch: %#v vs %#v", report1.Diagnostics, report2.Diagnostics)
	}
}
