package graph

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
	"github.com/flarebyte/baldrick-flying-buttress/internal/safety"
)

func RenderMarkdownText(ctx context.Context, selected Selected, shape Shape, policy CyclePolicy) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	adj := buildAdj(selected.Relationships)
	noteByID := map[string]domain.Note{}
	for _, note := range selected.Notes {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		noteByID[note.ID] = note
	}
	roots := selected.Roots
	if len(roots) == 0 {
		roots = noteIDs(selected.Notes)
	}
	anchors := buildNoteAnchors(selected.Notes)

	var b strings.Builder
	budget := renderBudget{}
	switch shape {
	case ShapeTree:
		for _, root := range roots {
			if err := renderTree(ctx, &b, root, 0, adj, noteByID, anchors, &budget); err != nil {
				return "", err
			}
		}
	case ShapeDAG:
		rendered := map[string]bool{}
		for _, root := range roots {
			if err := renderDAG(ctx, &b, root, 0, adj, noteByID, anchors, rendered, &budget); err != nil {
				return "", err
			}
		}
	case ShapeCyclic:
		if policy == CyclePolicyAllow {
			rendered := map[string]bool{}
			for _, root := range roots {
				if rendered[root] {
					continue
				}
				if err := renderCyclic(ctx, &b, root, 0, adj, noteByID, anchors, rendered, map[string]bool{}, &budget); err != nil {
					return "", err
				}
			}
			writeCyclicAdjacencySummary(ctx, &b, selected.Notes, noteByID, adj, anchors)
		}
	}
	return b.String(), nil
}

func renderTree(ctx context.Context, b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note, anchors map[string]string, budget *renderBudget) error {
	note, ok, err := resolveRenderableNote(ctx, nodeID, noteByID, budget)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	writeNodeLine(b, note, depth, anchors[note.ID])
	for _, child := range adj[nodeID] {
		if err := renderTree(ctx, b, child, depth+1, adj, noteByID, anchors, budget); err != nil {
			return err
		}
	}
	return nil
}

func renderDAG(ctx context.Context, b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note, anchors map[string]string, rendered map[string]bool, budget *renderBudget) error {
	note, ok, err := resolveRenderableNote(ctx, nodeID, noteByID, budget)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	if rendered[nodeID] {
		writeReferenceLine(b, note, depth, anchors[note.ID], "see")
		return nil
	}
	rendered[nodeID] = true
	writeNodeLine(b, note, depth, anchors[note.ID])
	for _, child := range adj[nodeID] {
		if err := renderDAG(ctx, b, child, depth+1, adj, noteByID, anchors, rendered, budget); err != nil {
			return err
		}
	}
	return nil
}

func renderCyclic(ctx context.Context, b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note, anchors map[string]string, rendered map[string]bool, stack map[string]bool, budget *renderBudget) error {
	note, ok, err := resolveRenderableNote(ctx, nodeID, noteByID, budget)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	if stack[nodeID] {
		writeReferenceLine(b, note, depth, anchors[note.ID], "cycle back to")
		return nil
	}
	if rendered[nodeID] {
		writeReferenceLine(b, note, depth, anchors[note.ID], "see")
		return nil
	}
	rendered[nodeID] = true
	writeNodeLine(b, note, depth, anchors[note.ID])
	stack[nodeID] = true
	for _, child := range adj[nodeID] {
		if err := renderCyclic(ctx, b, child, depth+1, adj, noteByID, anchors, rendered, stack, budget); err != nil {
			return err
		}
	}
	delete(stack, nodeID)
	return nil
}

type renderBudget struct {
	used int
}

func (b *renderBudget) use() error {
	b.used++
	if err := safety.CheckGraphRenderNodeCount(b.used); err != nil {
		return fmt.Errorf("graph render limit exceeded: %w", err)
	}
	return nil
}

func resolveRenderableNote(ctx context.Context, nodeID string, noteByID map[string]domain.Note, budget *renderBudget) (domain.Note, bool, error) {
	if err := ctx.Err(); err != nil {
		return domain.Note{}, false, err
	}
	if err := budget.use(); err != nil {
		return domain.Note{}, false, err
	}
	note, ok := noteByID[nodeID]
	return note, ok, nil
}

func writeNodeLine(b *strings.Builder, note domain.Note, depth int, anchor string) {
	indent := strings.Repeat("  ", depth)
	b.WriteString(indent)
	b.WriteString("- ")
	if strings.TrimSpace(anchor) != "" {
		b.WriteString("<a id=\"")
		b.WriteString(anchor)
		b.WriteString("\"></a> ")
	}
	if strings.TrimSpace(note.Title) != "" {
		b.WriteString(note.Title)
	} else {
		b.WriteString(note.ID)
	}
	if strings.TrimSpace(note.Markdown) != "" {
		b.WriteString(": ")
		b.WriteString(note.Markdown)
	}
	b.WriteByte('\n')
}

func writeReferenceLine(b *strings.Builder, note domain.Note, depth int, anchor, prefix string) {
	indent := strings.Repeat("  ", depth)
	b.WriteString(indent)
	b.WriteString("- *(")
	b.WriteString(prefix)
	b.WriteString(" [")
	b.WriteString(noteDisplayName(note))
	b.WriteString("](#")
	b.WriteString(anchor)
	b.WriteString("))*\n")
}

func noteDisplayName(note domain.Note) string {
	if strings.TrimSpace(note.Title) != "" {
		return note.Title
	}
	return note.ID
}

func buildNoteAnchors(notes []domain.Note) map[string]string {
	ids := make([]string, 0, len(notes))
	for _, note := range notes {
		ids = append(ids, note.ID)
	}
	ids = ordering.Strings(ids)
	anchors := make(map[string]string, len(ids))
	seenBase := map[string]int{}
	for _, id := range ids {
		base := stableAnchorSlug(id)
		seenBase[base]++
		if seenBase[base] == 1 {
			anchors[id] = "graph-node-" + base
			continue
		}
		anchors[id] = fmt.Sprintf("graph-node-%s-%d", base, seenBase[base])
	}
	return anchors
}

func stableAnchorSlug(value string) string {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	var b strings.Builder
	lastDash := false
	for _, r := range trimmed {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		return "node"
	}
	return slug
}

func writeCyclicAdjacencySummary(ctx context.Context, b *strings.Builder, notes []domain.Note, noteByID map[string]domain.Note, adj map[string][]string, anchors map[string]string) {
	if len(notes) == 0 {
		return
	}
	ids := noteIDs(notes)
	lines := make([]string, 0, len(ids))
	for _, fromID := range ids {
		if err := ctx.Err(); err != nil {
			return
		}
		toIDs := adj[fromID]
		if len(toIDs) == 0 {
			continue
		}
		fromNote, ok := noteByID[fromID]
		if !ok {
			continue
		}
		parts := make([]string, 0, len(toIDs))
		for _, toID := range toIDs {
			toNote, ok := noteByID[toID]
			if !ok {
				continue
			}
			parts = append(parts, fmt.Sprintf("[%s](#%s)", noteDisplayName(toNote), anchors[toID]))
		}
		if len(parts) == 0 {
			continue
		}
		lines = append(lines, fmt.Sprintf("- [%s](#%s) -> %s", noteDisplayName(fromNote), anchors[fromID], strings.Join(parts, ", ")))
	}
	if len(lines) == 0 {
		return
	}
	if b.Len() > 0 {
		b.WriteByte('\n')
	}
	b.WriteString("Adjacency summary:\n")
	for _, line := range lines {
		b.WriteString(line)
		b.WriteByte('\n')
	}
}
