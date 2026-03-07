package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func TestGenerateMarkdownFileBackedRendering(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureBundle(t, tmp, "config.markdown.file.raw.json", []string{
		"fixtures/data.csv",
		"fixtures/diagram.png",
		"fixtures/flow.mmd",
	})
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	output, err := os.ReadFile(filepath.Join(tmp, "out", "file.md"))
	if err != nil {
		t.Fatalf("read file-backed report failed: %v", err)
	}
	if string(output) != readGolden(t, "generate_markdown_file_output.golden") {
		t.Fatalf("file-backed markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_file_output.golden"), string(output))
	}
}

func TestGenerateMarkdownFileBackedMissingFileRuntimeFailure(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.file.missing.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != outcome.ExitCodeRuntimeFailure {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeRuntimeFailure, code)
	}
	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if stderr == "" {
		t.Fatal("expected stderr")
	}
}

func TestGenerateMarkdownFileBackedDeterministicAcrossRuns(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureBundle(t, tmp, "config.markdown.file.raw.json", []string{
		"fixtures/data.csv",
		"fixtures/diagram.png",
		"fixtures/flow.mmd",
	})
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code1, stdout1, stderr1 := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code1 != 0 {
		t.Fatalf("expected first exit code 0, got %d", code1)
	}
	first, err := os.ReadFile(filepath.Join(tmp, "out", "file.md"))
	if err != nil {
		t.Fatalf("read first file-backed report failed: %v", err)
	}

	code2, stdout2, stderr2 := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code2 != 0 {
		t.Fatalf("expected second exit code 0, got %d", code2)
	}
	second, err := os.ReadFile(filepath.Join(tmp, "out", "file.md"))
	if err != nil {
		t.Fatalf("read second file-backed report failed: %v", err)
	}

	if stdout1 != stdout2 || stderr1 != stderr2 {
		t.Fatalf("non-deterministic file-backed command output")
	}
	if string(first) != string(second) {
		t.Fatalf("non-deterministic file-backed markdown\\nfirst: %q\\nsecond: %q", string(first), string(second))
	}
}
