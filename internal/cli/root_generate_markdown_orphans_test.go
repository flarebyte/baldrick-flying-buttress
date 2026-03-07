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

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.orphans.subject.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	output, err := os.ReadFile(filepath.Join(tmp, "out", "orphans-subject.md"))
	if err != nil {
		t.Fatalf("read orphan subject report failed: %v", err)
	}
	if string(output) != readGolden(t, "generate_markdown_orphans_subject_output.golden") {
		t.Fatalf("orphan subject markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_orphans_subject_output.golden"), string(output))
	}
}

func TestGenerateMarkdownOrphansWithEdgeFilter(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.orphans.edge.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	output, err := os.ReadFile(filepath.Join(tmp, "out", "orphans-edge.md"))
	if err != nil {
		t.Fatalf("read orphan edge report failed: %v", err)
	}
	if string(output) != readGolden(t, "generate_markdown_orphans_edge_output.golden") {
		t.Fatalf("orphan edge markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_orphans_edge_output.golden"), string(output))
	}
}

func TestGenerateMarkdownOrphansWithCounterpartFilter(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.orphans.counterpart.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	output, err := os.ReadFile(filepath.Join(tmp, "out", "orphans-counterpart.md"))
	if err != nil {
		t.Fatalf("read orphan counterpart report failed: %v", err)
	}
	if string(output) != readGolden(t, "generate_markdown_orphans_counterpart_output.golden") {
		t.Fatalf("orphan counterpart markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_orphans_counterpart_output.golden"), string(output))
	}
}

func TestGenerateMarkdownOrphansWithDirectionOverride(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.orphans.direction.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	output, err := os.ReadFile(filepath.Join(tmp, "out", "orphans-direction.md"))
	if err != nil {
		t.Fatalf("read orphan direction report failed: %v", err)
	}
	if string(output) != readGolden(t, "generate_markdown_orphans_direction_output.golden") {
		t.Fatalf("orphan direction markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_orphans_direction_output.golden"), string(output))
	}
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
