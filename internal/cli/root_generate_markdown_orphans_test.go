package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func TestGenerateMarkdownOrphansSubjectOnly(t *testing.T) {
	t.Parallel()

	tmp, code, stdout, stderr := runGenerateMarkdownFixture(t, "config.markdown.orphans.subject.raw.json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "orphans-subject.md"), "generate_markdown_orphans_subject_output.golden")
}

func TestGenerateMarkdownOrphansWithEdgeFilter(t *testing.T) {
	t.Parallel()

	tmp, code, stdout, stderr := runGenerateMarkdownFixture(t, "config.markdown.orphans.edge.raw.json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "orphans-edge.md"), "generate_markdown_orphans_edge_output.golden")
}

func TestGenerateMarkdownOrphansWithCounterpartFilter(t *testing.T) {
	t.Parallel()

	tmp, code, stdout, stderr := runGenerateMarkdownFixture(t, "config.markdown.orphans.counterpart.raw.json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "orphans-counterpart.md"), "generate_markdown_orphans_counterpart_output.golden")
}

func TestGenerateMarkdownOrphansWithDirectionOverride(t *testing.T) {
	t.Parallel()

	tmp, code, stdout, stderr := runGenerateMarkdownFixture(t, "config.markdown.orphans.direction.raw.json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "orphans-direction.md"), "generate_markdown_orphans_direction_output.golden")
}

func TestGenerateMarkdownOrphansDeterministicAcrossRuns(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.orphans.subject.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code1, stdout1, stderr1 := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code1 != 0 {
		t.Fatalf("expected first exit code 0, got %d", code1)
	}
	first, err := os.ReadFile(filepath.Join(tmp, "out", "orphans-subject.md"))
	if err != nil {
		t.Fatalf("read first orphan report failed: %v", err)
	}

	code2, stdout2, stderr2 := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code2 != 0 {
		t.Fatalf("expected second exit code 0, got %d", code2)
	}
	second, err := os.ReadFile(filepath.Join(tmp, "out", "orphans-subject.md"))
	if err != nil {
		t.Fatalf("read second orphan report failed: %v", err)
	}

	if stdout1 != stdout2 || stderr1 != stderr2 {
		t.Fatalf("non-deterministic orphan command output")
	}
	if string(first) != string(second) {
		t.Fatalf("non-deterministic orphan markdown\\nfirst: %q\\nsecond: %q", string(first), string(second))
	}
}
