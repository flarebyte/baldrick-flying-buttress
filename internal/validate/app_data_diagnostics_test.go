package validate

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestFormatDiagnosticMessageDeterministicOrder(t *testing.T) {
	t.Parallel()

	got := formatDiagnosticMessage("base", domain.Diagnostic{
		SectionTitle:     "S",
		ReportTitle:      "R",
		ArgumentName:     "arg",
		NoteName:         "n1",
		LabelValue:       "l1",
		RelationshipFrom: "a",
		RelationshipTo:   "b",
		SubjectLabel:     "sub",
		EdgeLabel:        "edge",
		CounterpartLabel: "cp",
	})
	want := "base [reportTitle=R, sectionTitle=S, noteName=n1, relationship=a->b, argumentName=arg, subjectLabel=sub, edgeLabel=edge, counterpartLabel=cp, labelValue=l1]"
	if got != want {
		t.Fatalf("message mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestNewArgumentDiagnosticIncludesContextInMessage(t *testing.T) {
	t.Parallel()

	got := newArgumentDiagnostic("FBC002", "unknown configured argument key", "reports[0].sections[0].arguments[0]", "CPU", "Overview", "n1", "x")
	want := "unknown configured argument key [reportTitle=CPU, sectionTitle=Overview, noteName=n1, argumentName=x]"
	if got.Message != want {
		t.Fatalf("message mismatch\nwant: %q\n got: %q", want, got.Message)
	}
}
