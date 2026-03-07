package cli

import (
	"context"
	"errors"
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

	if diagnostics, err := writeMarkdownReports(context.Background(), app); err != nil {
		t.Fatalf("first write failed: %v", err)
	} else if len(diagnostics) != 0 {
		t.Fatalf("expected no diagnostics on first write, got %#v", diagnostics)
	}
	first, err := os.ReadFile(filepath.Join(tmp, "out/report.md"))
	if err != nil {
		t.Fatalf("read first output failed: %v", err)
	}

	if diagnostics, err := writeMarkdownReports(context.Background(), app); err != nil {
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

func TestRenderFileCSVDeterministic(t *testing.T) {
	t.Parallel()

	data := []byte("name,kind,status\ncli.root,note,active\ncli.worker,note,inactive\n")
	first, err := renderFileCSV(context.Background(), data, noteArgs{formatCSV: "table"})
	if err != nil {
		t.Fatalf("render csv failed: %v", err)
	}
	second, err := renderFileCSV(context.Background(), data, noteArgs{formatCSV: "table"})
	if err != nil {
		t.Fatalf("render csv second failed: %v", err)
	}
	want := "| kind | name | status |\n| --- | --- | --- |\n| note | cli.root | active |\n| note | cli.worker | inactive |"
	if first != want {
		t.Fatalf("csv table mismatch\nwant: %q\n got: %q", want, first)
	}
	if second != first {
		t.Fatalf("csv rendering is non-deterministic\nfirst: %q\nsecond: %q", first, second)
	}
}

func TestRenderFileCSVIncludeFilter(t *testing.T) {
	t.Parallel()

	data := []byte("name,kind,status\ncli.root,note,active\ncli.worker,note,inactive\n")
	got, err := renderFileCSV(context.Background(), data, noteArgs{
		formatCSV: "table",
		include:   csvFilter{column: "status", value: "active"},
	})
	if err != nil {
		t.Fatalf("render csv include failed: %v", err)
	}
	want := "| kind | name | status |\n| --- | --- | --- |\n| note | cli.root | active |"
	if got != want {
		t.Fatalf("csv include mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderFileCSVExcludeFilter(t *testing.T) {
	t.Parallel()

	data := []byte("name,kind,status\ncli.root,note,active\ncli.worker,note,inactive\n")
	got, err := renderFileCSV(context.Background(), data, noteArgs{
		formatCSV: "table",
		exclude:   csvFilter{column: "status", value: "inactive"},
	})
	if err != nil {
		t.Fatalf("render csv exclude failed: %v", err)
	}
	want := "| kind | name | status |\n| --- | --- | --- |\n| note | cli.root | active |"
	if got != want {
		t.Fatalf("csv exclude mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderFileCSVIncludeExcludeFilters(t *testing.T) {
	t.Parallel()

	data := []byte("name,kind,status\ncli.root,note,active\ncli.worker,note,inactive\ncli.jobs,note,active\n")
	got, err := renderFileCSV(context.Background(), data, noteArgs{
		formatCSV: "table",
		include:   csvFilter{column: "status", value: "active"},
		exclude:   csvFilter{column: "name", value: "cli.jobs"},
	})
	if err != nil {
		t.Fatalf("render csv include/exclude failed: %v", err)
	}
	want := "| kind | name | status |\n| --- | --- | --- |\n| note | cli.root | active |"
	if got != want {
		t.Fatalf("csv include/exclude mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderFileMedia(t *testing.T) {
	t.Parallel()

	got, err := renderFileMedia(domain.Note{ID: "n.image", Title: "Architecture", Filepath: "assets/arch.png"})
	if err != nil {
		t.Fatalf("render media failed: %v", err)
	}
	want := "![Architecture](assets/arch.png)"
	if got != want {
		t.Fatalf("media mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderFileCode(t *testing.T) {
	t.Parallel()

	got := renderFileCode([]byte("graph TD\nA-->B\n"), ".mmd")
	want := "```mermaid\ngraph TD\nA-->B\n```"
	if got != want {
		t.Fatalf("code mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestWriteMarkdownReportsCancelledContextNoOutputFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	app := domain.ValidatedApp{
		ConfigDir: tmp,
		Notes: []domain.Note{
			{ID: "n.a", Title: "A", Markdown: "Body A"},
		},
		MarkdownReports: []domain.MarkdownReport{{
			Title:    "Report",
			Filepath: "out/report.md",
			Sections: []domain.MarkdownH2Section{{
				Title: "H2",
				Sections: []domain.MarkdownH3Section{{
					Title:   "H3",
					NoteIDs: []string{"n.a"},
				}},
			}},
		}},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := writeMarkdownReports(ctx, app)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(tmp, "out", "report.md")); !os.IsNotExist(statErr) {
		t.Fatalf("expected no generated file, stat err: %v", statErr)
	}
}
