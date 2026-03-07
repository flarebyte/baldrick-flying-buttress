package ordering

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestReportsSortedDeterministically(t *testing.T) {
	t.Parallel()

	in := []domain.Report{
		{ID: "b", Title: "B"},
		{ID: "a", Title: "C"},
		{ID: "a", Title: "A"},
	}
	got := Reports(in)
	want := []domain.Report{
		{ID: "a", Title: "A"},
		{ID: "a", Title: "C"},
		{ID: "b", Title: "B"},
	}
	assertReports(t, got, want)
}

func TestNotesSortedDeterministically(t *testing.T) {
	t.Parallel()

	in := []domain.Note{
		{ID: "2", Label: "svc"},
		{ID: "1", Label: "svc"},
		{ID: "9", Label: "api"},
	}
	got := Notes(in)
	want := []domain.Note{
		{ID: "9", Label: "api"},
		{ID: "1", Label: "svc"},
		{ID: "2", Label: "svc"},
	}
	assertNotes(t, got, want)
}

func TestRelationshipsSortedDeterministically(t *testing.T) {
	t.Parallel()

	in := []domain.Relationship{
		{FromID: "b", ToID: "a", Label: "z"},
		{FromID: "a", ToID: "c", Label: "a"},
		{FromID: "a", ToID: "b", Label: "b"},
		{FromID: "a", ToID: "b", Label: "a"},
	}
	got := Relationships(in)
	want := []domain.Relationship{
		{FromID: "a", ToID: "b", Label: "a"},
		{FromID: "a", ToID: "b", Label: "b"},
		{FromID: "a", ToID: "c", Label: "a"},
		{FromID: "b", ToID: "a", Label: "z"},
	}
	assertRelationships(t, got, want)
}

func TestDiagnosticsSortedDeterministically(t *testing.T) {
	t.Parallel()

	in := []domain.Diagnostic{
		{Code: "B", Severity: domain.SeverityError, Path: "p", Message: "m"},
		{Code: "A", Severity: domain.SeverityWarning, Path: "p2", Message: "z"},
		{Code: "A", Severity: domain.SeverityWarning, Path: "p1", Message: "a"},
	}
	got := Diagnostics(in)
	want := []domain.Diagnostic{
		{Code: "A", Severity: domain.SeverityWarning, Path: "p1", Message: "a"},
		{Code: "A", Severity: domain.SeverityWarning, Path: "p2", Message: "z"},
		{Code: "B", Severity: domain.SeverityError, Path: "p", Message: "m"},
	}
	assertDiagnostics(t, got, want)
}

func assertReports(t *testing.T, got, want []domain.Report) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("index %d mismatch: got %#v want %#v", i, got[i], want[i])
		}
	}
}

func assertNotes(t *testing.T, got, want []domain.Note) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("index %d mismatch: got %#v want %#v", i, got[i], want[i])
		}
	}
}

func assertRelationships(t *testing.T, got, want []domain.Relationship) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("index %d mismatch: got %#v want %#v", i, got[i], want[i])
		}
	}
}

func assertDiagnostics(t *testing.T, got, want []domain.Diagnostic) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("index %d mismatch: got %#v want %#v", i, got[i], want[i])
		}
	}
}
