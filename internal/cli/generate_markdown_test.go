package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
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
