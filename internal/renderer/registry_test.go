package renderer

import (
	"context"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/graph"
)

func TestRegistryResolveContainsBuiltins(t *testing.T) {
	t.Parallel()
	r := ResolveRegistry()
	if len(r.Capabilities) != 2 || r.Capabilities[0].Name != "markdown-text" || r.Capabilities[1].Name != "mermaid" {
		t.Fatalf("unexpected registry: %#v", r)
	}
}

func TestRendererArgumentPrecedenceNoteOverH3OverDefault(t *testing.T) {
	t.Parallel()
	defaultRenderer := "markdown-text"
	defaultDirection := "TD"
	app := domain.ValidatedApp{Registry: domain.ArgumentRegistry{Arguments: []domain.ArgumentDefinition{
		{Name: "graph-renderer", DefaultValue: &defaultRenderer},
		{Name: "mermaid-direction", DefaultValue: &defaultDirection},
	}}}
	h3 := domain.MarkdownH3Section{Arguments: []string{"graph-renderer=mermaid", "mermaid-direction=LR"}, NoteIDs: []string{"n1"}}
	noteByID := map[string]domain.Note{"n1": {ID: "n1", ArgumentsCSV: "graph-renderer=markdown-text"}}

	args, err := ResolveArgs(app, h3, noteByID)
	if err != nil {
		t.Fatalf("resolve args failed: %v", err)
	}
	if args.Renderer != "markdown-text" || args.MermaidDirection != "LR" {
		t.Fatalf("unexpected resolved args: %#v", args)
	}
}

func TestDeterministicFallbackWhenRendererUnspecified(t *testing.T) {
	t.Parallel()
	args, err := ResolveArgs(domain.ValidatedApp{}, domain.MarkdownH3Section{}, map[string]domain.Note{})
	if err != nil {
		t.Fatalf("resolve args failed: %v", err)
	}
	if args.Renderer != "markdown-text" {
		t.Fatalf("expected markdown-text fallback, got %s", args.Renderer)
	}
	if args.RendererExplicit {
		t.Fatal("expected renderer to be non-explicit")
	}
}

func TestMermaidRenderingTreeDAGCyclic(t *testing.T) {
	t.Parallel()
	selected := graph.Selected{
		Notes:         []domain.Note{{ID: "a", Title: "A"}, {ID: "b", Title: "B"}, {ID: "c", Title: "C"}},
		Relationships: []domain.Relationship{{FromID: "a", ToID: "b"}, {FromID: "b", ToID: "c"}, {FromID: "c", ToID: "a"}},
	}
	got, err := renderMermaid(context.Background(), selected, graph.ShapeCyclic, Args{MermaidDirection: "TD"})
	if err != nil {
		t.Fatalf("render mermaid failed: %v", err)
	}
	want := "```mermaid\ngraph TD\n    N1[A]\n    N2[B]\n    N3[C]\n    N1 --> N2\n    N2 --> N3\n    N3 --> N1\n```\n"
	if got != want {
		t.Fatalf("mermaid mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestShapeCompatibilityHandling(t *testing.T) {
	t.Parallel()
	r := Registry{Capabilities: []Capability{{Name: "x", SupportedShapes: []graph.Shape{graph.ShapeTree}}}}
	_, err := r.Select("x", graph.ShapeDAG)
	if err == nil {
		t.Fatal("expected shape incompatibility error")
	}
}

func TestSelectResolvedFallbackByShape(t *testing.T) {
	t.Parallel()

	r := ResolveRegistry()
	args := Args{Renderer: "markdown-text", MermaidDirection: "TD", RendererExplicit: false}

	treeCap, err := r.SelectResolved(args, graph.ShapeTree)
	if err != nil {
		t.Fatalf("select tree fallback failed: %v", err)
	}
	if treeCap.Name != "markdown-text" {
		t.Fatalf("expected markdown-text for tree fallback, got %s", treeCap.Name)
	}

	dagCap, err := r.SelectResolved(args, graph.ShapeDAG)
	if err != nil {
		t.Fatalf("select dag fallback failed: %v", err)
	}
	if dagCap.Name != "markdown-text" {
		t.Fatalf("expected markdown-text for dag fallback, got %s", dagCap.Name)
	}

	cyclicCap, err := r.SelectResolved(args, graph.ShapeCyclic)
	if err != nil {
		t.Fatalf("select cyclic fallback failed: %v", err)
	}
	if cyclicCap.Name != "mermaid" {
		t.Fatalf("expected mermaid for cyclic fallback, got %s", cyclicCap.Name)
	}
}

func TestSelectResolvedExplicitRendererOverridesFallback(t *testing.T) {
	t.Parallel()

	r := ResolveRegistry()
	args := Args{Renderer: "markdown-text", MermaidDirection: "TD", RendererExplicit: true}
	capability, err := r.SelectResolved(args, graph.ShapeCyclic)
	if err != nil {
		t.Fatalf("select explicit renderer failed: %v", err)
	}
	if capability.Name != "markdown-text" {
		t.Fatalf("expected explicit markdown-text override, got %s", capability.Name)
	}
}
