package cli

import (
	"path/filepath"
	"testing"
)

func TestGenerateMarkdownGraphRendering(t *testing.T) {
	t.Parallel()

	tmp, code, stdout, stderr := runGenerateMarkdownFixture(t, "config.markdown.graph.raw.json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "graph.md"), "generate_markdown_graph_output.golden")
}

func TestGenerateMarkdownGraphCyclePolicyDisallowSkipsSection(t *testing.T) {
	t.Parallel()

	tmp, code, stdout, stderr := runGenerateMarkdownFixture(t, "config.markdown.graph.cycle.raw.json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "generate_markdown_graph_cycle_diagnostic_output.golden"), "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "cycle.md"), "generate_markdown_graph_cycle_output.golden")
}

func TestGenerateMarkdownRendererExplicitMermaid(t *testing.T) {
	t.Parallel()

	tmp, code, stdout, stderr := runGenerateMarkdownFixture(t, "config.markdown.renderer.explicit.raw.json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "renderer-explicit.md"), "generate_markdown_renderer_explicit_output.golden")
}

func TestGenerateMarkdownRendererFallbackWhenUnspecified(t *testing.T) {
	t.Parallel()

	tmp, code, stdout, stderr := runGenerateMarkdownFixture(t, "config.markdown.renderer.fallback.raw.json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "renderer-fallback.md"), "generate_markdown_renderer_fallback_output.golden")
}

func TestGenerateMarkdownRendererNoteLevelOverride(t *testing.T) {
	t.Parallel()

	tmp, code, stdout, stderr := runGenerateMarkdownFixture(t, "config.markdown.renderer.noteoverride.raw.json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "renderer-noteoverride.md"), "generate_markdown_renderer_noteoverride_output.golden")
}
