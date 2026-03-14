package graph

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/safety"
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
	got, err := RenderMarkdownText(context.Background(), selected, ShapeTree, CyclePolicyDisallow, Query{})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	want := "- <a id=\"graph-node-root\"></a> Root: root body\n  - <a id=\"graph-node-left\"></a> Left: left body\n  - <a id=\"graph-node-right\"></a> Right: right body\n"
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
	got, err := RenderMarkdownText(context.Background(), selected, ShapeDAG, CyclePolicyDisallow, Query{})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	want := "- <a id=\"graph-node-a\"></a> A: a\n  - <a id=\"graph-node-b\"></a> B: b\n    - <a id=\"graph-node-d\"></a> D: d\n  - <a id=\"graph-node-c\"></a> C: c\n    - *(see [D](#graph-node-d))*\n"
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
	got, err := RenderMarkdownText(context.Background(), selected, ShapeCyclic, CyclePolicyAllow, Query{})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	want := "- <a id=\"graph-node-a\"></a> A: a\n  - <a id=\"graph-node-b\"></a> B: b\n    - <a id=\"graph-node-c\"></a> C: c\n      - *(cycle back to [A](#graph-node-a))*\n\nAdjacency summary:\n- [A](#graph-node-a) -> [B](#graph-node-b)\n- [B](#graph-node-b) -> [C](#graph-node-c)\n- [C](#graph-node-c) -> [A](#graph-node-a)\n"
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
	got, err := RenderMarkdownText(context.Background(), selected, ShapeCyclic, CyclePolicyDisallow, Query{})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty render when cycle-policy disallow, got %q", got)
	}
}

func TestRenderMarkdownTextNodeLimitExceeded(t *testing.T) {
	t.Parallel()

	notes := make([]domain.Note, 0, safety.MaxGraphRenderNodesPerSection+1)
	for i := 0; i < safety.MaxGraphRenderNodesPerSection+1; i++ {
		notes = append(notes, domain.Note{
			ID:       "n" + strconv.Itoa(i),
			Title:    "N",
			Markdown: "m",
		})
	}
	selected := Selected{Notes: notes}

	_, err := RenderMarkdownText(context.Background(), selected, ShapeTree, CyclePolicyDisallow, Query{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "graph render limit exceeded") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStableAnchorSlugDeterministic(t *testing.T) {
	t.Parallel()

	got := stableAnchorSlug("  CLI.Root::A/B  ")
	want := "cli-root-a-b"
	if got != want {
		t.Fatalf("slug mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestBuildNoteAnchorsWithCollisionDeterministic(t *testing.T) {
	t.Parallel()

	anchors := buildNoteAnchors([]domain.Note{
		{ID: "n/a"},
		{ID: "n-a"},
	})
	if anchors["n-a"] != "graph-node-n-a" {
		t.Fatalf("unexpected anchor for n-a: %q", anchors["n-a"])
	}
	if anchors["n/a"] != "graph-node-n-a-2" {
		t.Fatalf("unexpected anchor for n/a: %q", anchors["n/a"])
	}
}

func TestSelectWithIncludeExcludeLabels(t *testing.T) {
	t.Parallel()

	app := domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "main", Title: "Main", Label: "graph", LabelsCSV: "graph,main"},
			{ID: "future", Title: "Future", Label: "graph", LabelsCSV: "graph,future"},
			{ID: "helper", Title: "Helper", Label: "graph", LabelsCSV: "graph,helper"},
		},
		Relationships: []domain.Relationship{
			{FromID: "main", ToID: "future", Label: "edge", LabelsCSV: "edge"},
			{FromID: "main", ToID: "helper", Label: "edge", LabelsCSV: "edge"},
		},
	}

	selected := Select(app, Query{SubjectLabel: "graph", IncludeLabels: []string{"future", "main"}, ExcludeLabels: []string{"helper"}})
	if len(selected.Notes) != 2 || len(selected.Relationships) != 1 {
		t.Fatalf("unexpected selection with include/exclude labels: %#v", selected)
	}
}

func TestSelectWithDepthLimit(t *testing.T) {
	t.Parallel()

	selected := Select(fixtureChainApp(), Query{SubjectLabel: "graph", StartNode: "a", MaxDepth: 1, MaxDepthSet: true})
	if len(selected.Notes) != 2 {
		t.Fatalf("expected 2 notes with depth limit, got %d", len(selected.Notes))
	}
	if len(selected.Relationships) != 1 {
		t.Fatalf("expected 1 relationship with depth limit, got %d", len(selected.Relationships))
	}
}

func TestSelectWithHiddenHelperNodes(t *testing.T) {
	t.Parallel()

	app := domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "a", Title: "A", Markdown: "a", Label: "graph", LabelsCSV: "graph"},
			{ID: "h", Title: "Helper", Markdown: "h", Label: "graph", LabelsCSV: "graph,helper"},
			{ID: "b", Title: "B", Markdown: "b", Label: "graph", LabelsCSV: "graph"},
		},
		Relationships: []domain.Relationship{
			{FromID: "a", ToID: "h", Label: "edge", LabelsCSV: "edge"},
			{FromID: "h", ToID: "b", Label: "edge", LabelsCSV: "edge"},
		},
	}

	selected := Select(app, Query{SubjectLabel: "graph", ShowHelperNodes: false, ShowHelpersSet: true, HelperLabel: "helper"})
	if len(selected.Notes) != 2 {
		t.Fatalf("expected helper node hidden, got notes=%d", len(selected.Notes))
	}
	if len(selected.Relationships) != 1 || selected.Relationships[0].FromID != "a" || selected.Relationships[0].ToID != "b" {
		t.Fatalf("expected helper path to collapse, got %#v", selected.Relationships)
	}
}

func TestRenderMarkdownTextWithBranchPriorityAndTitleOrder(t *testing.T) {
	t.Parallel()

	selected := Select(domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "root", Title: "Root", Markdown: "r", Label: "graph", LabelsCSV: "graph"},
			{ID: "n.future", Title: "Z Future", Markdown: "f", Label: "graph", LabelsCSV: "graph,future"},
			{ID: "n.main", Title: "A Main", Markdown: "m", Label: "graph", LabelsCSV: "graph,main"},
		},
		Relationships: []domain.Relationship{
			{FromID: "root", ToID: "n.future", Label: "edge", LabelsCSV: "edge"},
			{FromID: "root", ToID: "n.main", Label: "edge", LabelsCSV: "edge"},
		},
	}, Query{SubjectLabel: "graph"})

	got, err := RenderMarkdownText(context.Background(), selected, ShapeTree, CyclePolicyDisallow, Query{
		ChildOrder:     ChildOrderTitle,
		BranchPriority: []string{"main", "future"},
	})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	want := "- <a id=\"graph-node-root\"></a> Root: r\n  - <a id=\"graph-node-n-main\"></a> A Main: m\n  - <a id=\"graph-node-n-future\"></a> Z Future: f\n"
	if got != want {
		t.Fatalf("branch-priority render mismatch\nwant: %q\n got: %q", want, got)
	}
}

func TestResolveQueryInvalidArgs(t *testing.T) {
	t.Parallel()

	_, err := ResolveQuery([]string{"graph-max-depth=-1"})
	if err == nil {
		t.Fatal("expected invalid depth error")
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
	return fixtureAppWithEdges(
		[]domain.Note{
			{ID: "root", Title: "Root", Markdown: "root body", Label: "graph", LabelsCSV: "graph"},
			{ID: "left", Title: "Left", Markdown: "left body", Label: "graph", LabelsCSV: "graph"},
			{ID: "right", Title: "Right", Markdown: "right body", Label: "graph", LabelsCSV: "graph"},
		},
		[][2]string{{"root", "left"}, {"root", "right"}},
	)
}

func fixtureDAGApp() domain.ValidatedApp {
	return fixtureAppWithEdges(
		[]domain.Note{
			{ID: "a", Title: "A", Markdown: "a", Label: "graph", LabelsCSV: "graph"},
			{ID: "b", Title: "B", Markdown: "b", Label: "graph", LabelsCSV: "graph"},
			{ID: "c", Title: "C", Markdown: "c", Label: "graph", LabelsCSV: "graph"},
			{ID: "d", Title: "D", Markdown: "d", Label: "graph", LabelsCSV: "graph"},
		},
		[][2]string{{"a", "b"}, {"a", "c"}, {"b", "d"}, {"c", "d"}},
	)
}

func fixtureCycleApp() domain.ValidatedApp {
	return fixtureAppWithEdges(
		[]domain.Note{
			{ID: "a", Title: "A", Markdown: "a", Label: "graph", LabelsCSV: "graph"},
			{ID: "b", Title: "B", Markdown: "b", Label: "graph", LabelsCSV: "graph"},
			{ID: "c", Title: "C", Markdown: "c", Label: "graph", LabelsCSV: "graph"},
		},
		[][2]string{{"a", "b"}, {"b", "c"}, {"c", "a"}},
	)
}

func fixtureChainApp() domain.ValidatedApp {
	return fixtureAppWithEdges(
		[]domain.Note{
			{ID: "a", Title: "A", Markdown: "a", Label: "graph", LabelsCSV: "graph"},
			{ID: "b", Title: "B", Markdown: "b", Label: "graph", LabelsCSV: "graph"},
			{ID: "c", Title: "C", Markdown: "c", Label: "graph", LabelsCSV: "graph"},
		},
		[][2]string{{"a", "b"}, {"b", "c"}},
	)
}

func fixtureAppWithEdges(notes []domain.Note, edges [][2]string) domain.ValidatedApp {
	rels := make([]domain.Relationship, 0, len(edges))
	for _, edge := range edges {
		rels = append(rels, domain.Relationship{FromID: edge[0], ToID: edge[1], Label: "edge", LabelsCSV: "edge"})
	}
	return domain.ValidatedApp{Notes: notes, Relationships: rels}
}
