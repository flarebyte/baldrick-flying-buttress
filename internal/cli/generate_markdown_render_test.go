package cli

import (
	"context"
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

	got, diagnostics, err := renderMarkdownReport(context.Background(), report, notes, domain.ValidatedApp{}, renderer.ResolveRegistry())
	if err != nil {
		t.Fatalf("render markdown report failed: %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got %#v", diagnostics)
	}
	want := "# Inventory\n\n## Overview\n\n### Ingredients\n\n#### Apple\n\nFresh apple.\n\n"
	if got != want {
		t.Fatalf("markdown mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderMarkdownPlainNoteWithLabelsEnabled(t *testing.T) {
	t.Parallel()

	report := domain.MarkdownReport{
		Title: "Inventory",
		Sections: []domain.MarkdownH2Section{{
			Title: "Overview",
			Sections: []domain.MarkdownH3Section{{
				Title:     "Ingredients",
				Arguments: []string{"show-labels=true"},
				NoteIDs:   []string{"n.apple"},
			}},
		}},
	}
	notes := map[string]domain.Note{
		"n.apple": {ID: "n.apple", Title: "Apple", Markdown: "Fresh apple.", LabelsCSV: "future,v1"},
	}

	got, diagnostics, err := renderMarkdownReport(context.Background(), report, notes, domain.ValidatedApp{}, renderer.ResolveRegistry())
	if err != nil {
		t.Fatalf("render markdown report failed: %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got %#v", diagnostics)
	}
	want := "# Inventory\n\n## Overview\n\n### Ingredients\n\n#### Apple\n\nLabels: future, v1\n\nFresh apple.\n\n"
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

	first, firstDiagnostics, firstErr := renderMarkdownReport(context.Background(), report, map[string]domain.Note{}, domain.ValidatedApp{}, renderer.ResolveRegistry())
	second, secondDiagnostics, secondErr := renderMarkdownReport(context.Background(), report, map[string]domain.Note{}, domain.ValidatedApp{}, renderer.ResolveRegistry())
	if firstErr != nil || secondErr != nil {
		t.Fatalf("expected no render errors, got first=%v second=%v", firstErr, secondErr)
	}
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

	got, diagnostics, err := renderMarkdownReport(context.Background(), domain.MarkdownReport{Title: "Title"}, map[string]domain.Note{}, domain.ValidatedApp{}, renderer.ResolveRegistry())
	if err != nil {
		t.Fatalf("render markdown report failed: %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got %#v", diagnostics)
	}
	if got[len(got)-1] != '\n' {
		t.Fatalf("expected trailing newline, got %q", got)
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
