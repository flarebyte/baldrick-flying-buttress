package ordering

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestContractReportsPermutationInvariant(t *testing.T) {
	t.Parallel()

	a := []domain.Report{{ID: "b", Title: "B"}, {ID: "a", Title: "A"}}
	b := []domain.Report{{ID: "a", Title: "A"}, {ID: "b", Title: "B"}}
	assertEqualSlices(t, Reports(a), Reports(b))
}

func TestContractNotesPermutationInvariant(t *testing.T) {
	t.Parallel()

	a := []domain.Note{{ID: "2", Label: "svc"}, {ID: "1", Label: "svc"}, {ID: "0", Label: "api"}}
	b := []domain.Note{{ID: "0", Label: "api"}, {ID: "1", Label: "svc"}, {ID: "2", Label: "svc"}}
	assertEqualSlices(t, Notes(a), Notes(b))
}

func TestContractRelationshipsPermutationInvariant(t *testing.T) {
	t.Parallel()

	a := []domain.Relationship{{FromID: "a", ToID: "b", Label: "z"}, {FromID: "a", ToID: "b", Label: "a"}}
	b := []domain.Relationship{{FromID: "a", ToID: "b", Label: "a"}, {FromID: "a", ToID: "b", Label: "z"}}
	assertEqualSlices(t, Relationships(a), Relationships(b))
}

func TestContractDiagnosticsPermutationInvariant(t *testing.T) {
	t.Parallel()

	a := []domain.Diagnostic{{Code: "B", Severity: domain.SeverityError, Message: "m2", Path: "p"}, {Code: "A", Severity: domain.SeverityWarning, Message: "m1", Path: "p"}}
	b := []domain.Diagnostic{{Code: "A", Severity: domain.SeverityWarning, Message: "m1", Path: "p"}, {Code: "B", Severity: domain.SeverityError, Message: "m2", Path: "p"}}
	assertEqualSlices(t, Diagnostics(a), Diagnostics(b))
}
