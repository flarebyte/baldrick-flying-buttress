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

func TestFSAppLoaderLoadCuePackageFromFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	appPath := filepath.Join(tmp, "app.cue")
	if err := os.WriteFile(appPath, []byte(`package flyb

import "strings"

source: strings.ToUpper("fixture-app")
reports: [{
	title: "Overview"
	filepath: "reports/overview.md"
	sections: [{title: "Summary"}]
}]
notes: [{
	name: "n1"
	title: "service.api"
}]
`), 0o644); err != nil {
		t.Fatalf("write app.cue failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "relationships.cue"), []byte(`package flyb

relationships: [{
	from: "n1"
	to: "n1"
	label: "depends_on"
}]
`), 0o644); err != nil {
		t.Fatalf("write relationships.cue failed: %v", err)
	}

	got, err := FSAppLoader{ConfigPath: appPath}.Load(context.Background())
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if got.ConfigPath != appPath {
		t.Fatalf("expected config path %q, got %q", appPath, got.ConfigPath)
	}
	if got.Source != "FIXTURE-APP" {
		t.Fatalf("expected source FIXTURE-APP, got %q", got.Source)
	}
	if len(got.Relationships) != 1 {
		t.Fatalf("expected 1 relationship, got %d", len(got.Relationships))
	}
}

func TestFSAppLoaderLoadDirectoryResolvesAppCue(t *testing.T) {
	t.Parallel()

	got, err := FSAppLoader{ConfigPath: filepath.Join("..", "..", "doc", "design-meta")}.Load(context.Background())
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	wantPath := filepath.Join("..", "..", "doc", "design-meta", "app.cue")
	if got.ConfigPath != wantPath {
		t.Fatalf("expected config path %q, got %q", wantPath, got.ConfigPath)
	}
	if got.Source != "design-migration" {
		t.Fatalf("expected source design-migration, got %q", got.Source)
	}
	if len(got.Relationships) == 0 {
		t.Fatal("expected relationships loaded from sibling package files")
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
