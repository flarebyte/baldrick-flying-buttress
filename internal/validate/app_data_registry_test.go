package validate

import (
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

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

	_, report := validateRaw(t, raw)

	want := []diagExpectation{
		{code: "FBR000", path: "argumentRegistry.arguments[name=\"dup\"]"},
		{code: "FBR001", path: "argumentRegistry.arguments[0].name"},
		{code: "FBR002", path: "argumentRegistry.arguments[0].valueType"},
		{code: "FBR003", path: "argumentRegistry.arguments[name=\"fmt\"].valueType"},
		{code: "FBR004", path: "argumentRegistry.arguments[0].scopes"},
		{code: "FBR005", path: "argumentRegistry.arguments[name=\"fmt\"].scopes"},
		{code: "FBR006", path: "argumentRegistry.arguments[name=\"mode\"].defaultValue"},
		{code: "FBR007", path: "argumentRegistry.arguments[name=\"mode\"].allowedValues"},
	}
	assertHasDiagnostics(t, report.Diagnostics, want, "registry")
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

	app, _ := validateRaw(t, raw)
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
