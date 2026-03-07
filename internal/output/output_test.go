package output

import (
	"bytes"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestEmitDiagnosticsExactBytes(t *testing.T) {
	t.Parallel()

	report := domain.ValidationReport{
		Diagnostics: []domain.Diagnostic{
			{Code: "FB002", Severity: domain.SeverityError, Message: "error", Path: "b"},
			{Code: "FB001", Severity: domain.SeverityWarning, Message: "warn", Path: "a"},
		},
	}

	var b bytes.Buffer
	if err := EmitDiagnostics(&b, report); err != nil {
		t.Fatalf("emit diagnostics: %v", err)
	}

	want := "{\"diagnostics\":[{\"code\":\"FB001\",\"severity\":\"warning\",\"message\":\"warn\",\"path\":\"a\"},{\"code\":\"FB002\",\"severity\":\"error\",\"message\":\"error\",\"path\":\"b\"}]}\n"
	if b.String() != want {
		t.Fatalf("bytes mismatch\nwant: %q\n got: %q", want, b.String())
	}
}

func TestEmitReportListExactBytes(t *testing.T) {
	t.Parallel()

	app := domain.ValidatedApp{
		Reports: []domain.Report{
			{ID: "z", Title: "Zeta"},
			{ID: "a", Title: "Alpha"},
		},
	}

	var b bytes.Buffer
	if err := EmitReportList(&b, app); err != nil {
		t.Fatalf("emit report list: %v", err)
	}

	want := "{\"reports\":[{\"id\":\"a\",\"title\":\"Alpha\"},{\"id\":\"z\",\"title\":\"Zeta\"}]}\n"
	if b.String() != want {
		t.Fatalf("bytes mismatch\nwant: %q\n got: %q", want, b.String())
	}
}

func TestEmitGraphJSONExactBytes(t *testing.T) {
	t.Parallel()

	app := domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "n2", Label: "service.db"},
			{ID: "n1", Label: "service.api"},
		},
		Relationships: []domain.Relationship{
			{FromID: "n1", ToID: "n2", Label: "owns"},
			{FromID: "n1", ToID: "n2", Label: "depends_on"},
		},
	}

	var b bytes.Buffer
	if err := EmitGraphJSON(&b, app); err != nil {
		t.Fatalf("emit graph: %v", err)
	}

	want := "{\"notes\":[{\"id\":\"n1\",\"label\":\"service.api\"},{\"id\":\"n2\",\"label\":\"service.db\"}],\"relationships\":[{\"from\":\"n1\",\"to\":\"n2\",\"label\":\"depends_on\"},{\"from\":\"n1\",\"to\":\"n2\",\"label\":\"owns\"}]}\n"
	if b.String() != want {
		t.Fatalf("bytes mismatch\nwant: %q\n got: %q", want, b.String())
	}
}

func TestEmitterOutputUnaffectedByUnmappedDomainFields(t *testing.T) {
	t.Parallel()

	app := domain.ValidatedApp{
		Name:    "ignored",
		Modules: []string{"x", "y"},
		Reports: []domain.Report{
			{ID: "a", Title: "A"},
		},
	}

	var b bytes.Buffer
	if err := EmitReportList(&b, app); err != nil {
		t.Fatalf("emit report list: %v", err)
	}

	want := "{\"reports\":[{\"id\":\"a\",\"title\":\"A\"}]}\n"
	if b.String() != want {
		t.Fatalf("bytes mismatch\nwant: %q\n got: %q", want, b.String())
	}
}
