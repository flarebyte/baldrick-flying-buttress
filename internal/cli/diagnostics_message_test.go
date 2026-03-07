package cli

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestWithDiagnosticContextMessageDeterministicOrder(t *testing.T) {
	t.Parallel()

	got := withDiagnosticContextMessage(domain.Diagnostic{
		Message:          "m",
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
	want := "m [reportTitle=R, sectionTitle=S, noteName=n1, relationship=a->b, argumentName=arg, subjectLabel=sub, edgeLabel=edge, counterpartLabel=cp, labelValue=l1]"
	if got.Message != want {
		t.Fatalf("message mismatch\nwant: %q\n got: %q", want, got.Message)
	}
}
