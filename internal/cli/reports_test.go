package cli

import (
	"bytes"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestEmitReportsTableDeterministic(t *testing.T) {
	t.Parallel()

	reports := []domain.Report{
		{ID: "memory-health", Title: "Memory Health"},
		{ID: "cpu-overview", Title: "CPU Overview"},
	}

	var first bytes.Buffer
	if err := emitReportsTable(&first, reports); err != nil {
		t.Fatalf("emit first failed: %v", err)
	}
	var second bytes.Buffer
	if err := emitReportsTable(&second, reports); err != nil {
		t.Fatalf("emit second failed: %v", err)
	}

	want := readGolden(t, "list_reports_table_output.golden")
	if first.String() != want {
		t.Fatalf("stdout mismatch\nwant: %q\n got: %q", want, first.String())
	}
	if second.String() != first.String() {
		t.Fatalf("non-deterministic table output\nfirst: %q\nsecond: %q", first.String(), second.String())
	}
}

func TestEmitReportsTableVaryingTitleLengths(t *testing.T) {
	t.Parallel()

	reports := []domain.Report{
		{ID: "out/short.md", Title: "A"},
		{ID: "out/very-long-report.md", Title: "A Very Long Report Title"},
		{ID: "out/medium.md", Title: "Medium"},
	}

	var out bytes.Buffer
	if err := emitReportsTable(&out, reports); err != nil {
		t.Fatalf("emit failed: %v", err)
	}

	want := readGolden(t, "list_reports_table_varying_output.golden")
	if out.String() != want {
		t.Fatalf("stdout mismatch\nwant: %q\n got: %q", want, out.String())
	}
}
