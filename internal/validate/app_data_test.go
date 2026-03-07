package validate

import (
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

	_, _, err := validator.Validate(domain.RawApp{Source: "app"})
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	want := []string{
		"raw_model_normalization_precheck",
		"schema_structure_validation",
		"args_registry_resolve",
		"args_registry_validate",
		"args_validate_config",
		"labels_dataset_collect",
		"labels_reference_validate",
		"diagnostics_collection",
		"validated_app_normalization",
	}
	if !reflect.DeepEqual(steps, want) {
		t.Fatalf("step order mismatch: got %#v want %#v", steps, want)
	}
}

func TestAppDataValidatorSchemaDiagnostics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		raw       domain.RawApp
		wantCode  string
		wantPath  string
		wantCount int
	}{
		{
			name:     "missing report title",
			raw:      domain.RawApp{Source: "app", Reports: []domain.RawReport{{Title: "", Filepath: "reports/r1.md", Sections: []domain.RawReportSection{{Title: "S"}}}}},
			wantCode: "FBV101", wantPath: "reports[0].title", wantCount: 3,
		},
		{
			name:     "missing report filepath",
			raw:      domain.RawApp{Source: "app", Reports: []domain.RawReport{{Title: "R1", Filepath: "", Sections: []domain.RawReportSection{{Title: "S"}}}}},
			wantCode: "FBV102", wantPath: "reports[0].filepath", wantCount: 3,
		},
		{
			name:     "missing note name",
			raw:      domain.RawApp{Source: "app", Notes: []domain.RawNote{{Name: "", Title: "N"}}},
			wantCode: "FBV201", wantPath: "notes[0].name", wantCount: 3,
		},
		{
			name:     "missing note title",
			raw:      domain.RawApp{Source: "app", Notes: []domain.RawNote{{Name: "n1", Title: ""}}},
			wantCode: "FBV202", wantPath: "notes[0].title", wantCount: 3,
		},
		{
			name:     "missing relationship from",
			raw:      domain.RawApp{Source: "app", Relationships: []domain.RawRelationship{{FromID: "", ToID: "n2", Label: "depends_on"}}},
			wantCode: "FBV301", wantPath: "relationships[0].from", wantCount: 3,
		},
		{
			name:     "missing relationship to",
			raw:      domain.RawApp{Source: "app", Relationships: []domain.RawRelationship{{FromID: "n1", ToID: "", Label: "depends_on"}}},
			wantCode: "FBV302", wantPath: "relationships[0].to", wantCount: 3,
		},
		{
			name:     "missing report sections shape",
			raw:      domain.RawApp{Source: "app", Reports: []domain.RawReport{{Title: "R1", Filepath: "reports/r1.md", Sections: nil}}},
			wantCode: "FBV103", wantPath: "reports[0].sections", wantCount: 3,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, report, err := AppDataValidator{}.Validate(tc.raw)
			if err != nil {
				t.Fatalf("validate failed: %v", err)
			}
			if len(report.Diagnostics) != tc.wantCount {
				t.Fatalf("diagnostic count mismatch: got %d want %d", len(report.Diagnostics), tc.wantCount)
			}
			found := false
			for _, d := range report.Diagnostics {
				if d.Code == tc.wantCode && d.Path == tc.wantPath {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected diagnostic %s at %s, got %#v", tc.wantCode, tc.wantPath, report.Diagnostics)
			}
		})
	}
}

func TestAppDataValidatorRegistryDiagnostics(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source:        "app",
		Reports:       []domain.RawReport{{Title: "R", Filepath: "reports/r.md", Sections: []domain.RawReportSection{{Title: "S"}}}},
		Notes:         []domain.RawNote{{Name: "n1", Title: "N1"}},
		Relationships: []domain.RawRelationship{{FromID: "n1", ToID: "n2", Label: "L"}},
		Registry: domain.RawArgumentRegistry{Arguments: []domain.RawArgumentDefinition{
			{Name: "", ValueType: "", Scopes: nil},
			{Name: "fmt", ValueType: "bogus", Scopes: []string{"bad-scope"}},
			{Name: "mode", ValueType: "enum", Scopes: []string{"note"}, AllowedValues: []string{"a", "a"}, DefaultValue: "z"},
			{Name: "dup", ValueType: "string", Scopes: []string{"note"}},
			{Name: "dup", ValueType: "string", Scopes: []string{"note"}},
		}},
	}

	_, report, err := AppDataValidator{}.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	want := []struct {
		code string
		path string
	}{
		{code: "FBR000", path: "argumentRegistry.arguments[name=\"dup\"]"},
		{code: "FBR001", path: "argumentRegistry.arguments[0].name"},
		{code: "FBR002", path: "argumentRegistry.arguments[0].valueType"},
		{code: "FBR003", path: "argumentRegistry.arguments[name=\"fmt\"].valueType"},
		{code: "FBR004", path: "argumentRegistry.arguments[0].scopes"},
		{code: "FBR005", path: "argumentRegistry.arguments[name=\"fmt\"].scopes"},
		{code: "FBR006", path: "argumentRegistry.arguments[name=\"mode\"].defaultValue"},
		{code: "FBR007", path: "argumentRegistry.arguments[name=\"mode\"].allowedValues"},
	}
	for _, item := range want {
		found := false
		for _, d := range report.Diagnostics {
			if d.Code == item.code && d.Path == item.path {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected registry diagnostic %s at %s, got %#v", item.code, item.path, report.Diagnostics)
		}
	}
}

func TestAppDataValidatorRegistryNormalizedOrder(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source:        "app",
		Reports:       []domain.RawReport{{Title: "R", Filepath: "reports/r.md", Sections: []domain.RawReportSection{{Title: "S"}}}},
		Notes:         []domain.RawNote{{Name: "n1", Title: "N1"}},
		Relationships: []domain.RawRelationship{{FromID: "n1", ToID: "n2", Label: "L"}},
		Registry: domain.RawArgumentRegistry{Arguments: []domain.RawArgumentDefinition{
			{Name: "zeta", ValueType: "enum", Scopes: []string{"renderer", "note", "note"}, AllowedValues: []string{"b", "a", "a"}, DefaultValue: "a"},
			{Name: "alpha", ValueType: "string", Scopes: []string{"note"}},
		}},
	}

	app, _, err := AppDataValidator{}.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if len(app.Registry.Arguments) != 2 {
		t.Fatalf("expected 2 arguments, got %d", len(app.Registry.Arguments))
	}
	if app.Registry.Arguments[0].Name != "alpha" || app.Registry.Arguments[1].Name != "zeta" {
		t.Fatalf("expected normalized argument order, got %#v", app.Registry.Arguments)
	}
	if !reflect.DeepEqual(app.Registry.Arguments[1].AllowedValues, []string{"a", "b"}) {
		t.Fatalf("expected normalized allowed values, got %#v", app.Registry.Arguments[1].AllowedValues)
	}
	if !reflect.DeepEqual(app.Registry.Arguments[1].Scopes, []domain.ArgumentScope{domain.ArgumentScopeNote, domain.ArgumentScopeRenderer}) {
		t.Fatalf("expected normalized scopes, got %#v", app.Registry.Arguments[1].Scopes)
	}
}

func TestAppDataValidatorConfiguredArgumentsDiagnostics(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source: "app",
		Reports: []domain.RawReport{{
			Title:    "R",
			Filepath: "reports/r.md",
			Sections: []domain.RawReportSection{{
				Title:     "S",
				Arguments: []string{"unknown=x", "fmt=true", "mode=z", "badarg", "=x", "k="},
			}},
		}},
		Notes: []domain.RawNote{{
			Name:      "n1",
			Title:     "N1",
			Arguments: []string{"fmt=x", "verbose=maybe"},
		}},
		Relationships: []domain.RawRelationship{{FromID: "n1", ToID: "n2", Label: "L"}},
		Registry: domain.RawArgumentRegistry{Arguments: []domain.RawArgumentDefinition{
			{Name: "fmt", ValueType: "boolean", Scopes: []string{"renderer"}},
			{Name: "mode", ValueType: "enum", Scopes: []string{"h3-section"}, AllowedValues: []string{"a", "b"}},
			{Name: "verbose", ValueType: "boolean", Scopes: []string{"note"}},
		}},
	}

	_, report, err := AppDataValidator{}.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	want := []struct {
		code string
		path string
	}{
		{code: "FBC001", path: "reports[0].sections[0].arguments[3]"},
		{code: "FBC001", path: "reports[0].sections[0].arguments[4]"},
		{code: "FBC001", path: "reports[0].sections[0].arguments[5]"},
		{code: "FBC002", path: "reports[0].sections[0].arguments[0]"},
		{code: "FBC003", path: "notes[0].arguments[0]"},
		{code: "FBC004", path: "notes[0].arguments[0]"},
		{code: "FBC004", path: "notes[0].arguments[1]"},
		{code: "FBC004", path: "reports[0].sections[0].arguments[2]"},
	}
	for _, item := range want {
		found := false
		for _, d := range report.Diagnostics {
			if d.Code == item.code && d.Path == item.path {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected configured-args diagnostic %s at %s, got %#v", item.code, item.path, report.Diagnostics)
		}
	}
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

	_, report, err := AppDataValidator{}.Validate(raw)
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

func TestAppDataValidatorCanonicalDiagnosticLocationAndMetadata(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source:  "app",
		Reports: []domain.RawReport{{Title: "", Filepath: "", Sections: []domain.RawReportSection{{Title: ""}}}},
	}

	_, report, err := AppDataValidator{}.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	want := []domain.Diagnostic{
		{Code: "FBV101", Severity: domain.SeverityError, Source: schemaValidationSource, Message: "missing required field: report title", Location: "reports[0].title", Path: "reports[0].title"},
		{Code: "FBV102", Severity: domain.SeverityError, Source: schemaValidationSource, Message: "missing required field: report filepath", Location: "reports[0].filepath", Path: "reports[0].filepath"},
		{Code: "FBV104", Severity: domain.SeverityError, Source: schemaValidationSource, Message: "missing required field: section title", Location: "reports[0].sections[0].title", Path: "reports[0].sections[0].title"},
		{Code: "FBV200", Severity: domain.SeverityError, Source: schemaValidationSource, Message: "missing required collection: notes", Location: "notes", Path: "notes"},
		{Code: "FBV300", Severity: domain.SeverityError, Source: schemaValidationSource, Message: "missing required collection: relationships", Location: "relationships", Path: "relationships"},
	}
	if len(report.Diagnostics) != len(want) {
		t.Fatalf("diagnostic count mismatch: got %d want %d", len(report.Diagnostics), len(want))
	}
	for i := range want {
		got := report.Diagnostics[i]
		if got.Code != want[i].Code || got.Severity != want[i].Severity || got.Source != want[i].Source || got.Message != want[i].Message || got.Location != want[i].Location || got.Path != want[i].Path {
			t.Fatalf("diagnostic %d mismatch: got %#v want %#v", i, got, want[i])
		}
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

	app1, report1, err1 := AppDataValidator{}.Validate(raw)
	if err1 != nil {
		t.Fatalf("first validate failed: %v", err1)
	}
	app2, report2, err2 := AppDataValidator{}.Validate(raw)
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
