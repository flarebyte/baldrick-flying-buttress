package cli

import (
	"path/filepath"
	"testing"
)

func TestGenerateMarkdownFileBackedRendering(t *testing.T) {
	t.Parallel()

	tmp, code, stdout, stderr := runGenerateMarkdownBundleFixture(t, "config.markdown.file.raw.json", []string{
		"fixtures/data.csv",
		"fixtures/diagram.png",
		"fixtures/flow.mmd",
	})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "file.md"), "generate_markdown_file_output.golden")
}

func TestGenerateMarkdownFileBackedMissingFileRuntimeFailure(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.file.missing.raw.json")
	code, stdout, stderr := runGenerateMarkdownWithConfig(configPath)
	assertRuntimeFailureOutput(t, code, stdout, stderr)
}

func TestGenerateMarkdownFileBackedDeterministicAcrossRuns(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureBundle(t, tmp, "config.markdown.file.raw.json", []string{
		"fixtures/data.csv",
		"fixtures/diagram.png",
		"fixtures/flow.mmd",
	})
	code1, stdout1, stderr1 := runGenerateMarkdownWithConfig(configPath)
	if code1 != 0 {
		t.Fatalf("expected first exit code 0, got %d", code1)
	}
	first := readGeneratedMarkdown(t, tmp, filepath.Join("out", "file.md"))

	code2, stdout2, stderr2 := runGenerateMarkdownWithConfig(configPath)
	if code2 != 0 {
		t.Fatalf("expected second exit code 0, got %d", code2)
	}
	second := readGeneratedMarkdown(t, tmp, filepath.Join("out", "file.md"))

	if stdout1 != stdout2 || stderr1 != stderr2 {
		t.Fatalf("non-deterministic file-backed command output")
	}
	if first != second {
		t.Fatalf("non-deterministic file-backed markdown\\nfirst: %q\\nsecond: %q", first, second)
	}
}
