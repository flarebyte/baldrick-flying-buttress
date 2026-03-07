package fsio

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteFileAtomicCreateNewFile(t *testing.T) {
	dir := t.TempDir()
	destination := filepath.Join(dir, "out", "report.md")
	want := []byte("# title\n")

	if err := WriteFileAtomic(context.Background(), destination, want, 0o644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	got, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("content mismatch\nwant: %q\ngot: %q", string(want), string(got))
	}
}

func TestWriteFileAtomicReplaceExistingFile(t *testing.T) {
	dir := t.TempDir()
	destination := filepath.Join(dir, "report.md")
	if err := os.WriteFile(destination, []byte("old\n"), 0o644); err != nil {
		t.Fatalf("seed file failed: %v", err)
	}

	want := []byte("new\n")
	if err := WriteFileAtomic(context.Background(), destination, want, 0o644); err != nil {
		t.Fatalf("replace failed: %v", err)
	}

	got, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("content mismatch\nwant: %q\ngot: %q", string(want), string(got))
	}
}

func TestWriteFileAtomicRenameFailureLeavesOriginalUnchanged(t *testing.T) {
	dir := t.TempDir()
	destination := filepath.Join(dir, "report.md")
	if err := os.WriteFile(destination, []byte("original\n"), 0o644); err != nil {
		t.Fatalf("seed file failed: %v", err)
	}

	originalRename := renameFn
	renameFn = func(string, string) error {
		return errors.New("rename failed")
	}
	defer func() { renameFn = originalRename }()

	err := WriteFileAtomic(context.Background(), destination, []byte("new\n"), 0o644)
	if err == nil {
		t.Fatal("expected error")
	}

	got, readErr := os.ReadFile(destination)
	if readErr != nil {
		t.Fatalf("read failed: %v", readErr)
	}
	if string(got) != "original\n" {
		t.Fatalf("original file changed\ngot: %q", string(got))
	}
}

func TestWriteFileAtomicFailureCleansTempFileBestEffort(t *testing.T) {
	dir := t.TempDir()
	destination := filepath.Join(dir, "report.md")

	originalRename := renameFn
	renameFn = func(string, string) error {
		return errors.New("rename failed")
	}
	defer func() { renameFn = originalRename }()

	err := WriteFileAtomic(context.Background(), destination, []byte("new\n"), 0o644)
	if err == nil {
		t.Fatal("expected error")
	}

	entries, listErr := os.ReadDir(dir)
	if listErr != nil {
		t.Fatalf("list dir failed: %v", listErr)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".flyb-write-") {
			t.Fatalf("temp file not cleaned up: %s", entry.Name())
		}
	}
}

func TestWriteFileAtomicCancellationNoDestinationAndNoTempFile(t *testing.T) {
	dir := t.TempDir()
	destination := filepath.Join(dir, "report.md")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := WriteFileAtomic(ctx, destination, []byte("new\n"), 0o644)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}

	if _, statErr := os.Stat(destination); !os.IsNotExist(statErr) {
		t.Fatalf("expected destination to not exist, got err: %v", statErr)
	}

	entries, listErr := os.ReadDir(dir)
	if listErr != nil {
		t.Fatalf("list dir failed: %v", listErr)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".flyb-write-") {
			t.Fatalf("temp file not cleaned up: %s", entry.Name())
		}
	}
}
