package output

import (
	"bytes"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestContractDiagnosticsEmitterBytesAndRunTwice(t *testing.T) {
	t.Parallel()

	report := domain.ValidationReport{
		Diagnostics: []domain.Diagnostic{{Code: "B", Severity: domain.SeverityError, Message: "m2", Path: "p2"}, {Code: "A", Severity: domain.SeverityWarning, Message: "m1", Path: "p1"}},
	}

	var first bytes.Buffer
	if err := EmitDiagnostics(&first, report); err != nil {
		t.Fatalf("emit diagnostics first: %v", err)
	}
	var second bytes.Buffer
	if err := EmitDiagnostics(&second, report); err != nil {
		t.Fatalf("emit diagnostics second: %v", err)
	}

	want := []byte("{\"diagnostics\":[{\"code\":\"A\",\"severity\":\"warning\",\"message\":\"m1\",\"path\":\"p1\"},{\"code\":\"B\",\"severity\":\"error\",\"message\":\"m2\",\"path\":\"p2\"}]}\n")
	if !bytes.Equal(first.Bytes(), want) {
		t.Fatalf("unexpected first bytes: %q", first.Bytes())
	}
	if !bytes.Equal(second.Bytes(), want) {
		t.Fatalf("unexpected second bytes: %q", second.Bytes())
	}
}

func TestContractReportListEmitterPermutationInvariant(t *testing.T) {
	t.Parallel()

	a := domain.ValidatedApp{Reports: []domain.Report{{ID: "b", Title: "B"}, {ID: "a", Title: "A"}}}
	b := domain.ValidatedApp{Reports: []domain.Report{{ID: "a", Title: "A"}, {ID: "b", Title: "B"}}}

	var outA bytes.Buffer
	if err := EmitReportList(&outA, a); err != nil {
		t.Fatalf("emit A: %v", err)
	}
	var outB bytes.Buffer
	if err := EmitReportList(&outB, b); err != nil {
		t.Fatalf("emit B: %v", err)
	}
	if !bytes.Equal(outA.Bytes(), outB.Bytes()) {
		t.Fatalf("bytes mismatch\nA: %q\nB: %q", outA.Bytes(), outB.Bytes())
	}
}

func TestContractGraphEmitterBytesAndRunTwice(t *testing.T) {
	t.Parallel()

	app := domain.ValidatedApp{
		Notes:         []domain.Note{{ID: "n2", Label: "b"}, {ID: "n1", Label: "a"}},
		Relationships: []domain.Relationship{{FromID: "n1", ToID: "n2", Label: "z"}, {FromID: "n1", ToID: "n2", Label: "a"}},
	}

	var first bytes.Buffer
	if err := EmitGraphJSON(&first, app); err != nil {
		t.Fatalf("emit graph first: %v", err)
	}
	var second bytes.Buffer
	if err := EmitGraphJSON(&second, app); err != nil {
		t.Fatalf("emit graph second: %v", err)
	}

	want := []byte("{\"notes\":[{\"id\":\"n1\",\"label\":\"a\"},{\"id\":\"n2\",\"label\":\"b\"}],\"relationships\":[{\"from\":\"n1\",\"to\":\"n2\",\"label\":\"a\"},{\"from\":\"n1\",\"to\":\"n2\",\"label\":\"z\"}]}\n")
	if !bytes.Equal(first.Bytes(), want) {
		t.Fatalf("unexpected first bytes: %q", first.Bytes())
	}
	if !bytes.Equal(second.Bytes(), want) {
		t.Fatalf("unexpected second bytes: %q", second.Bytes())
	}
}
