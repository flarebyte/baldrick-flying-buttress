package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func TestGenerateMarkdownGraphRendering(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.graph.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	graphReport, err := os.ReadFile(filepath.Join(tmp, "out", "graph.md"))
	if err != nil {
		t.Fatalf("read graph report failed: %v", err)
	}
	if string(graphReport) != readGolden(t, "generate_markdown_graph_output.golden") {
		t.Fatalf("graph markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_graph_output.golden"), string(graphReport))
	}
}

func TestGenerateMarkdownGraphCyclePolicyDisallowSkipsSection(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.graph.cycle.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "generate_markdown_graph_cycle_diagnostic_output.golden"), "")

	graphReport, err := os.ReadFile(filepath.Join(tmp, "out", "cycle.md"))
	if err != nil {
		t.Fatalf("read cycle report failed: %v", err)
	}
	if string(graphReport) != readGolden(t, "generate_markdown_graph_cycle_output.golden") {
		t.Fatalf("cycle markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_graph_cycle_output.golden"), string(graphReport))
	}
}

func TestGenerateMarkdownRendererExplicitMermaid(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.renderer.explicit.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	output, err := os.ReadFile(filepath.Join(tmp, "out", "renderer-explicit.md"))
	if err != nil {
		t.Fatalf("read renderer explicit report failed: %v", err)
	}
	if string(output) != readGolden(t, "generate_markdown_renderer_explicit_output.golden") {
		t.Fatalf("renderer explicit markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_renderer_explicit_output.golden"), string(output))
	}
}

func TestGenerateMarkdownRendererFallbackWhenUnspecified(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.renderer.fallback.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	output, err := os.ReadFile(filepath.Join(tmp, "out", "renderer-fallback.md"))
	if err != nil {
		t.Fatalf("read renderer fallback report failed: %v", err)
	}
	if string(output) != readGolden(t, "generate_markdown_renderer_fallback_output.golden") {
		t.Fatalf("renderer fallback markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_renderer_fallback_output.golden"), string(output))
	}
}

func TestGenerateMarkdownRendererNoteLevelOverride(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.renderer.noteoverride.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	output, err := os.ReadFile(filepath.Join(tmp, "out", "renderer-noteoverride.md"))
	if err != nil {
		t.Fatalf("read renderer noteoverride report failed: %v", err)
	}
	if string(output) != readGolden(t, "generate_markdown_renderer_noteoverride_output.golden") {
		t.Fatalf("renderer noteoverride markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_renderer_noteoverride_output.golden"), string(output))
	}
}
