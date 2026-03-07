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
		includes:  []csvFilter{{column: "status", value: "active"}},
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
		excludes:  []csvFilter{{column: "status", value: "inactive"}},
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
		includes:  []csvFilter{{column: "status", value: "active"}},
		excludes:  []csvFilter{{column: "name", value: "cli.jobs"}},
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

func TestRenderFileCSVMultipleIncludeFilters(t *testing.T) {
	t.Parallel()

	data := []byte("name,kind,status\ncli.root,note,active\ncli.root,task,active\n")
	got, err := renderFileCSV(context.Background(), data, noteArgs{
		formatCSV: "table",
		includes: []csvFilter{
			{column: "status", value: "active"},
			{column: "kind", value: "note"},
		},
	})
	if err != nil {
		t.Fatalf("render csv multiple includes failed: %v", err)
	}
	want := "| kind | name | status |\n| --- | --- | --- |\n| note | cli.root | active |"
	if got != want {
		t.Fatalf("csv multiple include mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderFileCSVMultipleExcludeFilters(t *testing.T) {
	t.Parallel()

	data := []byte("name,kind,status\ncli.root,note,active\ncli.jobs,note,active\ncli.worker,note,inactive\n")
	got, err := renderFileCSV(context.Background(), data, noteArgs{
		formatCSV: "table",
		excludes: []csvFilter{
			{column: "name", value: "cli.jobs"},
			{column: "status", value: "inactive"},
		},
	})
	if err != nil {
		t.Fatalf("render csv multiple excludes failed: %v", err)
	}
	want := "| kind | name | status |\n| --- | --- | --- |\n| note | cli.root | active |"
	if got != want {
		t.Fatalf("csv multiple exclude mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderFileCSVExactMatchOnly(t *testing.T) {
	t.Parallel()

	data := []byte("name,status\ncli.root,active\n")
	got, err := renderFileCSV(context.Background(), data, noteArgs{
		formatCSV: "table",
		includes:  []csvFilter{{column: "status", value: "act"}},
	})
	if err != nil {
		t.Fatalf("render csv exact-match failed: %v", err)
	}
	want := "| name | status |\n| --- | --- |"
	if got != want {
		t.Fatalf("csv exact-match mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderFileCSVUnknownIncludeColumnMatchesNoRows(t *testing.T) {
	t.Parallel()

	data := []byte("name,status\ncli.root,active\n")
	got, err := renderFileCSV(context.Background(), data, noteArgs{
		formatCSV: "table",
		includes:  []csvFilter{{column: "missing", value: "x"}},
	})
	if err != nil {
		t.Fatalf("render csv unknown include failed: %v", err)
	}
	want := "| name | status |\n| --- | --- |"
	if got != want {
		t.Fatalf("csv unknown include mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderFileCSVUnknownExcludeColumnIgnored(t *testing.T) {
	t.Parallel()

	data := []byte("name,status\ncli.root,active\n")
	got, err := renderFileCSV(context.Background(), data, noteArgs{
		formatCSV: "table",
		excludes:  []csvFilter{{column: "missing", value: "x"}},
	})
	if err != nil {
		t.Fatalf("render csv unknown exclude failed: %v", err)
	}
	want := "| name | status |\n| --- | --- |\n| cli.root | active |"
	if got != want {
		t.Fatalf("csv unknown exclude mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderFileCSVEmptyCellsAndEscaping(t *testing.T) {
	t.Parallel()

	data := []byte("name,desc\ncli.root,\"alpha|beta\"\ncli.worker,\n")
	got, err := renderFileCSV(context.Background(), data, noteArgs{formatCSV: "table"})
	if err != nil {
		t.Fatalf("render csv escaping failed: %v", err)
	}
	want := "| desc | name |\n| --- | --- |\n| alpha\\|beta | cli.root |\n|  | cli.worker |"
	if got != want {
		t.Fatalf("csv escaping mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestRenderFileCSVRawMode(t *testing.T) {
	t.Parallel()

	data := []byte("name,status\ncli.root,active\ncli.worker,inactive\n")
	got, err := renderFileCSV(context.Background(), data, noteArgs{
		formatCSV: "raw",
		includes:  []csvFilter{{column: "status", value: "active"}},
	})
	if err != nil {
		t.Fatalf("render csv raw failed: %v", err)
	}
	want := "```csv\nname,status\ncli.root,active\n```"
	if got != want {
		t.Fatalf("csv raw mismatch\nwant: %q\n got: %q", want, got)
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
