package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
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

	got := renderMarkdownReport(report, notes)
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

	first := renderMarkdownReport(report, map[string]domain.Note{})
	second := renderMarkdownReport(report, map[string]domain.Note{})
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

	got := renderMarkdownReport(domain.MarkdownReport{Title: "Title"}, map[string]domain.Note{})
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

	if err := writeMarkdownReports(app); err != nil {
		t.Fatalf("first write failed: %v", err)
	}
	first, err := os.ReadFile(filepath.Join(tmp, "out/report.md"))
	if err != nil {
		t.Fatalf("read first output failed: %v", err)
	}

	if err := writeMarkdownReports(app); err != nil {
		t.Fatalf("second write failed: %v", err)
	}
	second, err := os.ReadFile(filepath.Join(tmp, "out/report.md"))
	if err != nil {
		t.Fatalf("read second output failed: %v", err)
	}

	if string(first) != string(second) {
		t.Fatalf("non-deterministic file output\nfirst: %q\nsecond: %q", string(first), string(second))
	}
}
