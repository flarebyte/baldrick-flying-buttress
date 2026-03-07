package ordering

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestContractReportsPermutationInvariant(t *testing.T) {
	t.Parallel()

	a := []domain.Report{{ID: "b", Title: "B"}, {ID: "a", Title: "A"}}
	b := []domain.Report{{ID: "a", Title: "A"}, {ID: "b", Title: "B"}}
	assertReportsEqual(t, Reports(a), Reports(b))
}

func TestContractNotesPermutationInvariant(t *testing.T) {
	t.Parallel()

	a := []domain.Note{{ID: "2", Label: "svc"}, {ID: "1", Label: "svc"}, {ID: "0", Label: "api"}}
	b := []domain.Note{{ID: "0", Label: "api"}, {ID: "1", Label: "svc"}, {ID: "2", Label: "svc"}}
	assertNotesEqual(t, Notes(a), Notes(b))
}

func TestContractRelationshipsPermutationInvariant(t *testing.T) {
	t.Parallel()

	a := []domain.Relationship{{FromID: "a", ToID: "b", Label: "z"}, {FromID: "a", ToID: "b", Label: "a"}}
	b := []domain.Relationship{{FromID: "a", ToID: "b", Label: "a"}, {FromID: "a", ToID: "b", Label: "z"}}
	assertRelationshipsEqual(t, Relationships(a), Relationships(b))
}

func TestContractDiagnosticsPermutationInvariant(t *testing.T) {
	t.Parallel()

	a := []domain.Diagnostic{{Code: "B", Severity: domain.SeverityError, Message: "m2", Path: "p"}, {Code: "A", Severity: domain.SeverityWarning, Message: "m1", Path: "p"}}
	b := []domain.Diagnostic{{Code: "A", Severity: domain.SeverityWarning, Message: "m1", Path: "p"}, {Code: "B", Severity: domain.SeverityError, Message: "m2", Path: "p"}}
	assertDiagnosticsEqual(t, Diagnostics(a), Diagnostics(b))
}

func assertReportsEqual(t *testing.T, got, want []domain.Report) {
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

func assertNotesEqual(t *testing.T, got, want []domain.Note) {
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

func assertRelationshipsEqual(t *testing.T, got, want []domain.Relationship) {
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

func assertDiagnosticsEqual(t *testing.T, got, want []domain.Diagnostic) {
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
