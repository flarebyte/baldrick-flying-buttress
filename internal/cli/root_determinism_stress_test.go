package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func assertSameCommandResult(t *testing.T, runIndex int, code, wantCode int, out, wantOut, errOut, wantErr string) {
	t.Helper()
	if code != wantCode {
		t.Fatalf("run %d exit code mismatch: got %d want %d", runIndex, code, wantCode)
	}
	if out != wantOut {
		t.Fatalf("run %d stdout mismatch", runIndex)
	}
	if errOut != wantErr {
		t.Fatalf("run %d stderr mismatch", runIndex)
	}
}

func TestDeterminismStressNonFileCommands(t *testing.T) {
	t.Parallel()

	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	tests := []struct {
		name string
		args []string
	}{
		{name: "validate", args: []string{"validate", "--config", filepath.Join("testdata", "config.raw.json")}},
		{name: "list reports", args: []string{"list", "reports", "--config", filepath.Join("testdata", "config.raw.json")}},
		{name: "list names", args: []string{"list", "names", "--config", filepath.Join("testdata", "config.raw.json"), "--prefix", "n"}},
		{name: "lint names", args: []string{"lint", "names", "--config", filepath.Join("testdata", "config.lint.raw.json")}},
		{name: "lint orphans", args: []string{"lint", "orphans", "--config", filepath.Join("testdata", "config.orphans.raw.json"), "--subject-label", "ingredient"}},
		{name: "generate json", args: []string{"generate", "json", "--config", filepath.Join("testdata", "config.raw.json")}},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			firstCode, firstOut, firstErr := runCommandWithFactory(tc.args, loaderFactory, validator)
			for i := 0; i < 4; i++ {
				code, out, errOut := runCommandWithFactory(tc.args, loaderFactory, validator)
				assertSameCommandResult(t, i+2, code, firstCode, out, firstOut, errOut, firstErr)
			}
		})
	}
}

func TestDeterminismStressGenerateMarkdownFiles(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	firstCode, firstOut, firstErr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if firstCode != 0 {
		t.Fatalf("expected first exit code 0, got %d", firstCode)
	}
	firstAlpha, err := os.ReadFile(filepath.Join(tmp, "out", "alpha.md"))
	if err != nil {
		t.Fatalf("read first alpha failed: %v", err)
	}
	firstBeta, err := os.ReadFile(filepath.Join(tmp, "out", "beta.md"))
	if err != nil {
		t.Fatalf("read first beta failed: %v", err)
	}

	for i := 0; i < 4; i++ {
		code, out, errOut := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
		assertSameCommandResult(t, i+2, code, firstCode, out, firstOut, errOut, firstErr)
		alpha, err := os.ReadFile(filepath.Join(tmp, "out", "alpha.md"))
		if err != nil {
			t.Fatalf("read alpha failed on run %d: %v", i+2, err)
		}
		beta, err := os.ReadFile(filepath.Join(tmp, "out", "beta.md"))
		if err != nil {
			t.Fatalf("read beta failed on run %d: %v", i+2, err)
		}
		if string(alpha) != string(firstAlpha) {
			t.Fatalf("run %d alpha content mismatch", i+2)
		}
		if string(beta) != string(firstBeta) {
			t.Fatalf("run %d beta content mismatch", i+2)
		}
	}
}
