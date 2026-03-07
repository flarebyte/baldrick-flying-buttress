package domain

import "testing"

func TestValidationReportCanonicalSetsEmptySlice(t *testing.T) {
	t.Parallel()

	report := ValidationReport{}
	canonical := report.Canonical()
	if canonical.Diagnostics == nil {
		t.Fatal("expected diagnostics slice to be non-nil")
	}
	if len(canonical.Diagnostics) != 0 {
		t.Fatalf("expected empty diagnostics, got %d", len(canonical.Diagnostics))
	}
}

func TestValidatedAppModelFieldsAreExplicit(t *testing.T) {
	t.Parallel()

	app := ValidatedApp{
		Name:    "stub-app",
		Modules: []string{"core", "edge"},
		Reports: []Report{{ID: "r1", Title: "Report"}},
		Notes: []Note{
			{ID: "n1", Label: "service.api"},
			{ID: "n2", Label: "service.db"},
		},
		Relationships: []Relationship{{FromID: "n1", ToID: "n2", Label: "depends_on"}},
	}

	if len(app.Reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(app.Reports))
	}
	if len(app.Notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(app.Notes))
	}
	if len(app.Relationships) != 1 {
		t.Fatalf("expected 1 relationship, got %d", len(app.Relationships))
	}
}
