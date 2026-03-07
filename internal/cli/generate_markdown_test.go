package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/orphans"
	"github.com/flarebyte/baldrick-flying-buttress/internal/renderer"
)

func TestRenderMarkdownPlainNote(t *testing.T) {
	t.Parallel()

	report := domain.MarkdownReport{
		Title: "Inventory",
		Sections: []domain.MarkdownH2Section{{
			Title: "Overview",
			Sections: []domain.MarkdownH3Section{{
				Title:   "Ingredients",
				NoteIDs: []string{"n.apple"},
			}},
		}},
	}
	notes := map[string]domain.Note{"n.apple": {ID: "n.apple", Title: "Apple", Markdown: "Fresh apple."}}

	got, diagnostics := renderMarkdownReport(report, notes, domain.ValidatedApp{}, renderer.ResolveRegistry())
	if len(diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got %#v", diagnostics)
	}
	want := "# Inventory\n\n## Overview\n\n### Ingredients\n\n#### Apple\n\nFresh apple.\n\n"
	if got != want {
		t.Fatalf("markdown mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderMarkdownDeterministicSections(t *testing.T) {
	t.Parallel()

	report := domain.MarkdownReport{
		Title: "R",
		Sections: []domain.MarkdownH2Section{
			{Title: "B", Sections: []domain.MarkdownH3Section{{Title: "Y"}}},
			{Title: "A", Sections: []domain.MarkdownH3Section{{Title: "Z"}, {Title: "X"}}},
		},
	}

	first, firstDiagnostics := renderMarkdownReport(report, map[string]domain.Note{}, domain.ValidatedApp{}, renderer.ResolveRegistry())
	second, secondDiagnostics := renderMarkdownReport(report, map[string]domain.Note{}, domain.ValidatedApp{}, renderer.ResolveRegistry())
	if len(firstDiagnostics) != 0 || len(secondDiagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got first=%#v second=%#v", firstDiagnostics, secondDiagnostics)
	}
	want := "# R\n\n## A\n\n### X\n\n### Z\n\n## B\n\n### Y\n\n"
	if first != want {
		t.Fatalf("markdown mismatch\nwant: %q\n got: %q", want, first)
	}
	if second != first {
		t.Fatalf("non-deterministic markdown\nfirst: %q\nsecond: %q", first, second)
	}
}

func TestRenderMarkdownTrailingNewline(t *testing.T) {
	t.Parallel()

	got, diagnostics := renderMarkdownReport(domain.MarkdownReport{Title: "Title"}, map[string]domain.Note{}, domain.ValidatedApp{}, renderer.ResolveRegistry())
	if len(diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got %#v", diagnostics)
	}
	if got[len(got)-1] != '\n' {
		t.Fatalf("expected trailing newline, got %q", got)
	}
}

func TestWriteMarkdownReportsDeterministic(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	app := domain.ValidatedApp{
		ConfigDir: tmp,
		Notes: []domain.Note{
			{ID: "n.b", Title: "B", Markdown: "Body B"},
			{ID: "n.a", Title: "A", Markdown: "Body A"},
		},
		MarkdownReports: []domain.MarkdownReport{{
			Title:    "Report",
			Filepath: "out/report.md",
			Sections: []domain.MarkdownH2Section{{
				Title: "H2",
				Sections: []domain.MarkdownH3Section{{
					Title:   "H3",
					NoteIDs: []string{"n.b", "n.a"},
				}},
			}},
		}},
	}

	if diagnostics, err := writeMarkdownReports(app); err != nil {
		t.Fatalf("first write failed: %v", err)
	} else if len(diagnostics) != 0 {
		t.Fatalf("expected no diagnostics on first write, got %#v", diagnostics)
	}
	first, err := os.ReadFile(filepath.Join(tmp, "out/report.md"))
	if err != nil {
		t.Fatalf("read first output failed: %v", err)
	}

	if diagnostics, err := writeMarkdownReports(app); err != nil {
		t.Fatalf("second write failed: %v", err)
	} else if len(diagnostics) != 0 {
		t.Fatalf("expected no diagnostics on second write, got %#v", diagnostics)
	}
	second, err := os.ReadFile(filepath.Join(tmp, "out/report.md"))
	if err != nil {
		t.Fatalf("read second output failed: %v", err)
	}

	if string(first) != string(second) {
		t.Fatalf("non-deterministic file output\nfirst: %q\nsecond: %q", string(first), string(second))
	}
}

func TestResolveOrphanQueryFromH3Arguments(t *testing.T) {
	t.Parallel()

	query, orphanMode, err := resolveOrphanQuery([]string{"orphan-subject-label=ingredient", "orphan-direction=out"})
	if err != nil {
		t.Fatalf("resolve orphan query failed: %v", err)
	}
	if !orphanMode {
		t.Fatal("expected orphan mode")
	}
	if query.SubjectLabel != "ingredient" || query.Direction != "out" {
		t.Fatalf("unexpected orphan query: %#v", query)
	}
}

func TestRenderOrphanRowsDeterministic(t *testing.T) {
	t.Parallel()

	notes := []domain.Note{
		{ID: "n.b", Title: "Note B", LabelsCSV: "beta, ingredient"},
		{ID: "n.a", Title: "Note A", LabelsCSV: "alpha, ingredient"},
	}
	first := renderOrphanRows(notes)
	second := renderOrphanRows(notes)
	want := "| name | title | labels |\n| --- | --- | --- |\n| n.a | Note A | alpha, ingredient |\n| n.b | Note B | beta, ingredient |\n"
	if first != want {
		t.Fatalf("orphan rows mismatch\nwant: %q\n got: %q", want, first)
	}
	if second != first {
		t.Fatalf("non-deterministic orphan rows\nfirst: %q\nsecond: %q", first, second)
	}
}

func TestRenderOrphanRowsEmpty(t *testing.T) {
	t.Parallel()

	got := renderOrphanRows(nil)
	want := "| name | title | labels |\n| --- | --- | --- |\n"
	if got != want {
		t.Fatalf("unexpected empty orphan rows\nwant: %q\n got: %q", want, got)
	}
}

func TestOrphanRenderingWithFilters(t *testing.T) {
	t.Parallel()

	app := domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "n.a", Label: "ingredient", Title: "A", LabelsCSV: "ingredient"},
			{ID: "n.b", Label: "ingredient", Title: "B", LabelsCSV: "ingredient"},
			{ID: "n.c", Label: "tool", Title: "C", LabelsCSV: "tool"},
			{ID: "n.d", Label: "ingredient", Title: "D", LabelsCSV: "ingredient"},
		},
		Relationships: []domain.Relationship{
			{FromID: "n.a", ToID: "n.c", Label: "uses", LabelsCSV: "uses"},
			{FromID: "n.c", ToID: "n.b", Label: "feeds", LabelsCSV: "feeds"},
		},
	}

	query, _, err := resolveOrphanQuery([]string{
		"orphan-subject-label=ingredient",
		"orphan-edge-label=uses",
		"orphan-counterpart-label=tool",
		"orphan-direction=out",
	})
	if err != nil {
		t.Fatalf("resolve orphan query failed: %v", err)
	}
	got := orphans.Find(app, query)
	if len(got) != 2 || got[0].ID != "n.b" || got[1].ID != "n.d" {
		t.Fatalf("unexpected filtered orphans: %#v", got)
	}
}
