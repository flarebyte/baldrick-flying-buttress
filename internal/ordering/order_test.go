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
	want := []domain.Report{
		{ID: "a", Title: "A"},
		{ID: "a", Title: "C"},
		{ID: "b", Title: "B"},
	}
	assertEqualSlices(t, Reports(in), want)
}

func TestNotesSortedDeterministically(t *testing.T) {
	t.Parallel()

	in := []domain.Note{
		{ID: "2", Label: "svc"},
		{ID: "1", Label: "svc"},
		{ID: "9", Label: "api"},
	}
	want := []domain.Note{
		{ID: "9", Label: "api"},
		{ID: "1", Label: "svc"},
		{ID: "2", Label: "svc"},
	}
	assertEqualSlices(t, Notes(in), want)
}

func TestRelationshipsSortedDeterministically(t *testing.T) {
	t.Parallel()

	in := []domain.Relationship{
		{FromID: "b", ToID: "a", Label: "z"},
		{FromID: "a", ToID: "c", Label: "a"},
		{FromID: "a", ToID: "b", Label: "b"},
		{FromID: "a", ToID: "b", Label: "a"},
	}
	want := []domain.Relationship{
		{FromID: "a", ToID: "b", Label: "a"},
		{FromID: "a", ToID: "b", Label: "b"},
		{FromID: "a", ToID: "c", Label: "a"},
		{FromID: "b", ToID: "a", Label: "z"},
	}
	assertEqualSlices(t, Relationships(in), want)
}

func TestDiagnosticsSortedDeterministically(t *testing.T) {
	t.Parallel()

	in := []domain.Diagnostic{
		{Code: "B", Severity: domain.SeverityError, Path: "p", Message: "m"},
		{Code: "A", Severity: domain.SeverityWarning, Path: "p2", Message: "z"},
		{Code: "A", Severity: domain.SeverityWarning, Path: "p1", Message: "a"},
	}
	want := []domain.Diagnostic{
		{Code: "A", Severity: domain.SeverityWarning, Path: "p1", Message: "a"},
		{Code: "A", Severity: domain.SeverityWarning, Path: "p2", Message: "z"},
		{Code: "B", Severity: domain.SeverityError, Path: "p", Message: "m"},
	}
	assertEqualSlices(t, Diagnostics(in), want)
}
