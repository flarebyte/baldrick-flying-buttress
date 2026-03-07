package load

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestFSAppLoaderLoadSuccess(t *testing.T) {
	t.Parallel()

	path := filepath.Join("testdata", "app.raw.json")
	got, err := FSAppLoader{ConfigPath: path}.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if got.ConfigPath != path {
		t.Fatalf("expected config path %q, got %q", path, got.ConfigPath)
	}
	if got.Source != "fixture-app" {
		t.Fatalf("expected source fixture-app, got %q", got.Source)
	}
}

func TestFSAppLoaderMissingFile(t *testing.T) {
	t.Parallel()

	_, err := FSAppLoader{ConfigPath: filepath.Join("testdata", "missing.raw.json")}.Load()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "read config") {
		t.Fatalf("expected read config error, got %v", err)
	}
}

func TestFSAppLoaderMalformedInput(t *testing.T) {
	t.Parallel()

	_, err := FSAppLoader{ConfigPath: filepath.Join("testdata", "app.malformed.txt")}.Load()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "parse config") {
		t.Fatalf("expected parse config error, got %v", err)
	}
}
