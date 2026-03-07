package graph

import (
	"context"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestSelectByLabels(t *testing.T) {
	t.Parallel()
	app := fixtureGraphApp()
	selected := Select(app, Query{SubjectLabel: "ingredient", EdgeLabel: "uses", CounterpartLabel: "tool"})
	if len(selected.Notes) != 2 || len(selected.Relationships) != 1 {
		t.Fatalf("unexpected selection: notes=%d rels=%d", len(selected.Notes), len(selected.Relationships))
	}
}

func TestDetectTreeShapeAndRender(t *testing.T) {
	t.Parallel()
	selected := Select(fixtureTreeApp(), Query{SubjectLabel: "graph"})
	if shape := DetectShape(selected); shape != ShapeTree {
		t.Fatalf("expected tree shape, got %s", shape)
	}
	got, err := RenderMarkdownText(context.Background(), selected, ShapeTree, CyclePolicyDisallow)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	want := "- Root: root body\n  - Left: left body\n  - Right: right body\n"
	if got != want {
		t.Fatalf("tree render mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestDetectDAGAndRenderRepetitionLimit(t *testing.T) {
	t.Parallel()
	selected := Select(fixtureDAGApp(), Query{SubjectLabel: "graph"})
	if shape := DetectShape(selected); shape != ShapeDAG {
		t.Fatalf("expected dag shape, got %s", shape)
	}
	got, err := RenderMarkdownText(context.Background(), selected, ShapeDAG, CyclePolicyDisallow)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	want := "- A: a\n  - B: b\n    - D: d\n  - C: c\n    - D: d\n"
	if got != want {
		t.Fatalf("dag render mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestDetectCyclicAndRenderCycleBack(t *testing.T) {
	t.Parallel()
	selected := Select(fixtureCycleApp(), Query{SubjectLabel: "graph", StartNode: "a"})
	if shape := DetectShape(selected); shape != ShapeCyclic {
		t.Fatalf("expected cyclic shape, got %s", shape)
	}
	got, err := RenderMarkdownText(context.Background(), selected, ShapeCyclic, CyclePolicyAllow)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	want := "- A: a\n  - B: b\n    - C: c\n      - *(cycle back to A)*\n"
	if got != want {
		t.Fatalf("cycle render mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestResolveCyclePolicy(t *testing.T) {
	t.Parallel()
	policy, err := ResolveCyclePolicy([]string{"cycle-policy=disallow"})
	if err != nil || policy != CyclePolicyDisallow {
		t.Fatalf("expected disallow policy, got %s err=%v", policy, err)
	}
}

func TestCyclePolicyDisallowSkipsCyclicRendering(t *testing.T) {
	t.Parallel()
	selected := Select(fixtureCycleApp(), Query{SubjectLabel: "graph", StartNode: "a"})
	got, err := RenderMarkdownText(context.Background(), selected, ShapeCyclic, CyclePolicyDisallow)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty render when cycle-policy disallow, got %q", got)
	}
}

func fixtureGraphApp() domain.ValidatedApp {
	return domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "n.ingredient", Title: "Ingredient", Markdown: "i", Label: "ingredient", LabelsCSV: "ingredient"},
			{ID: "n.tool", Title: "Tool", Markdown: "t", Label: "tool", LabelsCSV: "tool"},
			{ID: "n.extra", Title: "Extra", Markdown: "e", Label: "other", LabelsCSV: "other"},
		},
		Relationships: []domain.Relationship{{FromID: "n.ingredient", ToID: "n.tool", Label: "uses", LabelsCSV: "uses"}},
	}
}

func fixtureTreeApp() domain.ValidatedApp {
	return domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "root", Title: "Root", Markdown: "root body", Label: "graph", LabelsCSV: "graph"},
			{ID: "left", Title: "Left", Markdown: "left body", Label: "graph", LabelsCSV: "graph"},
			{ID: "right", Title: "Right", Markdown: "right body", Label: "graph", LabelsCSV: "graph"},
		},
		Relationships: []domain.Relationship{{FromID: "root", ToID: "left", Label: "edge", LabelsCSV: "edge"}, {FromID: "root", ToID: "right", Label: "edge", LabelsCSV: "edge"}},
	}
}

func fixtureDAGApp() domain.ValidatedApp {
	return domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "a", Title: "A", Markdown: "a", Label: "graph", LabelsCSV: "graph"},
			{ID: "b", Title: "B", Markdown: "b", Label: "graph", LabelsCSV: "graph"},
			{ID: "c", Title: "C", Markdown: "c", Label: "graph", LabelsCSV: "graph"},
			{ID: "d", Title: "D", Markdown: "d", Label: "graph", LabelsCSV: "graph"},
		},
		Relationships: []domain.Relationship{{FromID: "a", ToID: "b", Label: "edge", LabelsCSV: "edge"}, {FromID: "a", ToID: "c", Label: "edge", LabelsCSV: "edge"}, {FromID: "b", ToID: "d", Label: "edge", LabelsCSV: "edge"}, {FromID: "c", ToID: "d", Label: "edge", LabelsCSV: "edge"}},
	}
}

func fixtureCycleApp() domain.ValidatedApp {
	return domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "a", Title: "A", Markdown: "a", Label: "graph", LabelsCSV: "graph"},
			{ID: "b", Title: "B", Markdown: "b", Label: "graph", LabelsCSV: "graph"},
			{ID: "c", Title: "C", Markdown: "c", Label: "graph", LabelsCSV: "graph"},
		},
		Relationships: []domain.Relationship{{FromID: "a", ToID: "b", Label: "edge", LabelsCSV: "edge"}, {FromID: "b", ToID: "c", Label: "edge", LabelsCSV: "edge"}, {FromID: "c", ToID: "a", Label: "edge", LabelsCSV: "edge"}},
	}
}
