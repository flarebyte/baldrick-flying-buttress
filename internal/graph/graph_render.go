package graph

import (
	"context"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
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

	var b strings.Builder
	switch shape {
	case ShapeTree:
		for _, root := range roots {
			if err := renderTree(ctx, &b, root, 0, adj, noteByID); err != nil {
				return "", err
			}
		}
	case ShapeDAG:
		visits := map[string]int{}
		for _, root := range roots {
			if err := renderDAG(ctx, &b, root, 0, adj, noteByID, visits, 2); err != nil {
				return "", err
			}
		}
	case ShapeCyclic:
		if policy == CyclePolicyAllow {
			rendered := map[string]bool{}
			for _, root := range roots {
				if err := renderCyclic(ctx, &b, root, 0, adj, noteByID, rendered, map[string]bool{}); err != nil {
					return "", err
				}
			}
		}
	}
	return b.String(), nil
}

func renderTree(ctx context.Context, b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	note, ok := noteByID[nodeID]
	if !ok {
		return nil
	}
	writeNodeLine(b, note, depth)
	for _, child := range adj[nodeID] {
		if err := renderTree(ctx, b, child, depth+1, adj, noteByID); err != nil {
			return err
		}
	}
	return nil
}

func renderDAG(ctx context.Context, b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note, visits map[string]int, maxVisits int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	note, ok := noteByID[nodeID]
	if !ok {
		return nil
	}
	if visits[nodeID] >= maxVisits {
		indent := strings.Repeat("  ", depth)
		b.WriteString(indent)
		b.WriteString("- *(repeat limit reached: ")
		b.WriteString(note.Title)
		b.WriteString(")*\n")
		return nil
	}
	visits[nodeID]++
	writeNodeLine(b, note, depth)
	for _, child := range adj[nodeID] {
		if err := renderDAG(ctx, b, child, depth+1, adj, noteByID, visits, maxVisits); err != nil {
			return err
		}
	}
	return nil
}

func renderCyclic(ctx context.Context, b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note, rendered map[string]bool, stack map[string]bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	note, ok := noteByID[nodeID]
	if !ok {
		return nil
	}
	if stack[nodeID] {
		indent := strings.Repeat("  ", depth)
		b.WriteString(indent)
		b.WriteString("- *(cycle back to ")
		b.WriteString(note.Title)
		b.WriteString(")*\n")
		return nil
	}
	if rendered[nodeID] {
		return nil
	}
	rendered[nodeID] = true
	writeNodeLine(b, note, depth)
	stack[nodeID] = true
	for _, child := range adj[nodeID] {
		if err := renderCyclic(ctx, b, child, depth+1, adj, noteByID, rendered, stack); err != nil {
			return err
		}
	}
	delete(stack, nodeID)
	return nil
}

func writeNodeLine(b *strings.Builder, note domain.Note, depth int) {
	indent := strings.Repeat("  ", depth)
	b.WriteString(indent)
	b.WriteString("- ")
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
