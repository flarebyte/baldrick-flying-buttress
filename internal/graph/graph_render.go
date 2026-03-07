package graph

import (
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func RenderMarkdownText(selected Selected, shape Shape, policy CyclePolicy) string {
	adj := buildAdj(selected.Relationships)
	noteByID := map[string]domain.Note{}
	for _, note := range selected.Notes {
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
			renderTree(&b, root, 0, adj, noteByID)
		}
	case ShapeDAG:
		visits := map[string]int{}
		for _, root := range roots {
			renderDAG(&b, root, 0, adj, noteByID, visits, 2)
		}
	case ShapeCyclic:
		if policy == CyclePolicyAllow {
			rendered := map[string]bool{}
			for _, root := range roots {
				renderCyclic(&b, root, 0, adj, noteByID, rendered, map[string]bool{})
			}
		}
	}
	return b.String()
}

func renderTree(b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note) {
	note, ok := noteByID[nodeID]
	if !ok {
		return
	}
	writeNodeLine(b, note, depth)
	for _, child := range adj[nodeID] {
		renderTree(b, child, depth+1, adj, noteByID)
	}
}

func renderDAG(b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note, visits map[string]int, maxVisits int) {
	note, ok := noteByID[nodeID]
	if !ok {
		return
	}
	if visits[nodeID] >= maxVisits {
		indent := strings.Repeat("  ", depth)
		b.WriteString(indent)
		b.WriteString("- *(repeat limit reached: ")
		b.WriteString(note.Title)
		b.WriteString(")*\n")
		return
	}
	visits[nodeID]++
	writeNodeLine(b, note, depth)
	for _, child := range adj[nodeID] {
		renderDAG(b, child, depth+1, adj, noteByID, visits, maxVisits)
	}
}

func renderCyclic(b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note, rendered map[string]bool, stack map[string]bool) {
	note, ok := noteByID[nodeID]
	if !ok {
		return
	}
	if stack[nodeID] {
		indent := strings.Repeat("  ", depth)
		b.WriteString(indent)
		b.WriteString("- *(cycle back to ")
		b.WriteString(note.Title)
		b.WriteString(")*\n")
		return
	}
	if rendered[nodeID] {
		return
	}
	rendered[nodeID] = true
	writeNodeLine(b, note, depth)
	stack[nodeID] = true
	for _, child := range adj[nodeID] {
		renderCyclic(b, child, depth+1, adj, noteByID, rendered, stack)
	}
	delete(stack, nodeID)
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
