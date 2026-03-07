package load

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/safety"
)

func TestFSAppLoaderLoadSuccess(t *testing.T) {
	t.Parallel()

	path := filepath.Join("testdata", "app.cue")
	got, err := FSAppLoader{ConfigPath: path}.Load(context.Background())
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

	_, err := FSAppLoader{ConfigPath: filepath.Join("testdata", "missing.cue")}.Load(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "stat config") {
		t.Fatalf("expected stat config error, got %v", err)
	}
}

func TestFSAppLoaderMalformedInput(t *testing.T) {
	t.Parallel()

	_, err := FSAppLoader{ConfigPath: filepath.Join("testdata", "app.malformed.cue")}.Load(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "parse config") {
		t.Fatalf("expected parse config error, got %v", err)
	}
}

func TestFSAppLoaderConfigTooLarge(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "too-large.cue")
	content := strings.Repeat("a", int(safety.MaxConfigFileBytes)+1)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	_, err := FSAppLoader{ConfigPath: path}.Load(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "config file too large") {
		t.Fatalf("expected size limit error, got %v", err)
	}
}
