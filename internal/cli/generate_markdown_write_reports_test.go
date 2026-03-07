package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

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

func TestWriteMarkdownReportsWorkersOneVsManyEquivalent(t *testing.T) {
	t.Parallel()

	appFor := func(dir string) domain.ValidatedApp {
		return domain.ValidatedApp{
			ConfigDir: dir,
			Notes: []domain.Note{
				{ID: "n.a", Title: "A", Markdown: "Body A"},
				{ID: "n.b", Title: "B", Markdown: "Body B"},
				{ID: "n.c", Title: "C", Markdown: "Body C"},
			},
			MarkdownReports: []domain.MarkdownReport{
				{Title: "R1", Filepath: "out/r1.md", Sections: []domain.MarkdownH2Section{{Title: "H2", Sections: []domain.MarkdownH3Section{{Title: "H3", NoteIDs: []string{"n.a"}}}}}},
				{Title: "R2", Filepath: "out/r2.md", Sections: []domain.MarkdownH2Section{{Title: "H2", Sections: []domain.MarkdownH3Section{{Title: "H3", NoteIDs: []string{"n.b"}}}}}},
				{Title: "R3", Filepath: "out/r3.md", Sections: []domain.MarkdownH2Section{{Title: "H2", Sections: []domain.MarkdownH3Section{{Title: "H3", NoteIDs: []string{"n.c"}}}}}},
			},
		}
	}

	dirOne := t.TempDir()
	dirMany := t.TempDir()

	diagnosticsOne, err := writeMarkdownReportsWithWorkers(context.Background(), appFor(dirOne), 1)
	if err != nil {
		t.Fatalf("workers=1 failed: %v", err)
	}
	if len(diagnosticsOne) != 0 {
		t.Fatalf("expected no diagnostics for workers=1, got %#v", diagnosticsOne)
	}

	diagnosticsMany, err := writeMarkdownReportsWithWorkers(context.Background(), appFor(dirMany), 4)
	if err != nil {
		t.Fatalf("workers=4 failed: %v", err)
	}
	if len(diagnosticsMany) != 0 {
		t.Fatalf("expected no diagnostics for workers=4, got %#v", diagnosticsMany)
	}

	paths := []string{"out/r1.md", "out/r2.md", "out/r3.md"}
	for _, rel := range paths {
		one, err := os.ReadFile(filepath.Join(dirOne, rel))
		if err != nil {
			t.Fatalf("read workers=1 output failed: %v", err)
		}
		many, err := os.ReadFile(filepath.Join(dirMany, rel))
		if err != nil {
			t.Fatalf("read workers=4 output failed: %v", err)
		}
		if string(one) != string(many) {
			t.Fatalf("output mismatch for %s\nworkers=1: %q\nworkers=4: %q", rel, string(one), string(many))
		}
	}
}

func TestWriteMarkdownReportsWorkerFailureCancelsAndNoFilesCommitted(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	app := domain.ValidatedApp{
		ConfigDir: tmp,
		Notes: []domain.Note{
			{ID: "n.ok", Title: "OK", Markdown: "Body"},
			{ID: "n.missing", Title: "Missing", Filepath: "fixtures/missing.md"},
		},
		MarkdownReports: []domain.MarkdownReport{
			{Title: "OK report", Filepath: "out/ok.md", Sections: []domain.MarkdownH2Section{{Title: "H2", Sections: []domain.MarkdownH3Section{{Title: "H3", NoteIDs: []string{"n.ok"}}}}}},
			{Title: "Bad report", Filepath: "out/bad.md", Sections: []domain.MarkdownH2Section{{Title: "H2", Sections: []domain.MarkdownH3Section{{Title: "H3", NoteIDs: []string{"n.missing"}}}}}},
		},
	}

	_, err := writeMarkdownReportsWithWorkers(context.Background(), app, 4)
	if err == nil {
		t.Fatal("expected error")
	}

	if _, statErr := os.Stat(filepath.Join(tmp, "out", "ok.md")); !os.IsNotExist(statErr) {
		t.Fatalf("expected no committed ok output, got stat err: %v", statErr)
	}
	if _, statErr := os.Stat(filepath.Join(tmp, "out", "bad.md")); !os.IsNotExist(statErr) {
		t.Fatalf("expected no committed bad output, got stat err: %v", statErr)
	}
}

func TestWriteMarkdownReportsConcurrentRacePath(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	reports := make([]domain.MarkdownReport, 0, 50)
	notes := make([]domain.Note, 0, 50)
	for i := 0; i < 50; i++ {
		id := "n." + strconv.Itoa(i)
		notes = append(notes, domain.Note{ID: id, Title: id, Markdown: "Body"})
		reports = append(reports, domain.MarkdownReport{
			Title:    "R" + strconv.Itoa(i),
			Filepath: "out/r" + strconv.Itoa(i) + ".md",
			Sections: []domain.MarkdownH2Section{{Title: "H2", Sections: []domain.MarkdownH3Section{{Title: "H3", NoteIDs: []string{id}}}}},
		})
	}
	app := domain.ValidatedApp{ConfigDir: tmp, Notes: notes, MarkdownReports: reports}

	diagnostics, err := writeMarkdownReportsWithWorkers(context.Background(), app, 8)
	if err != nil {
		t.Fatalf("concurrent write failed: %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got %#v", diagnostics)
	}
	for i := 0; i < 50; i++ {
		if _, err := os.ReadFile(filepath.Join(tmp, "out", "r"+strconv.Itoa(i)+".md")); err != nil {
			t.Fatalf("expected generated report r%d.md, err=%v", i, err)
		}
	}
}
