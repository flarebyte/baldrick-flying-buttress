package validate

import (
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestEnrichDiagnosticsAddsMachineFriendlyFields(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		ConfigPath: filepath.Join("testdata", "config.raw.json"),
		Reports: []domain.RawReport{{
			Title:    "CPU Overview",
			Filepath: "reports/cpu-overview.md",
		}},
		Notes: []domain.RawNote{{
			Name:  "n1",
			Title: "Service API",
		}},
	}

	got := enrichDiagnostics(raw, []domain.Diagnostic{{
		Code:             "GRAPH_MISSING_NODE",
		Path:             "relationships[0].to",
		ReportTitle:      "CPU Overview",
		NoteName:         "n1",
		RelationshipFrom: "n1",
		RelationshipTo:   "missing",
	}})
	if len(got) != 1 {
		t.Fatalf("expected one diagnostic, got %#v", got)
	}
	d := got[0]
	if d.NormalizedPath != "relationships[0].to" {
		t.Fatalf("unexpected normalized path: %#v", d)
	}
	if d.ConfigPath != filepath.Join("testdata", "config.raw.json") {
		t.Fatalf("unexpected config path: %#v", d)
	}
	if d.ReportID != "cpu-overview" {
		t.Fatalf("unexpected report id: %#v", d)
	}
	if d.NoteTitle != "Service API" {
		t.Fatalf("unexpected note title: %#v", d)
	}
	if len(d.RelatedNodes) != 2 || d.RelatedNodes[0] != "missing" || d.RelatedNodes[1] != "n1" {
		t.Fatalf("unexpected related nodes: %#v", d.RelatedNodes)
	}
	if len(d.SuggestedFixes) == 0 {
		t.Fatalf("expected suggested fixes: %#v", d)
	}
}
