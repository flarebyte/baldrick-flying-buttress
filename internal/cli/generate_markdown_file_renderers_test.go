package cli

import (
	"context"
	"strings"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/safety"
)

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

func TestRenderFileCSVTooLarge(t *testing.T) {
	t.Parallel()

	data := []byte(strings.Repeat("x", safety.MaxCSVFileBytes+1))
	_, err := renderFileCSV(context.Background(), data, noteArgs{formatCSV: "table"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "csv file too large") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderFileCSVRowsRenderedLimitExceeded(t *testing.T) {
	t.Parallel()

	var b strings.Builder
	b.WriteString("name,status\n")
	for i := 0; i < safety.MaxCSVRowsRenderedPerNote+1; i++ {
		b.WriteString("n")
		b.WriteString("x,ok\n")
	}
	_, err := renderFileCSV(context.Background(), []byte(b.String()), noteArgs{formatCSV: "table"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "csv rendered rows") {
		t.Fatalf("unexpected error: %v", err)
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
