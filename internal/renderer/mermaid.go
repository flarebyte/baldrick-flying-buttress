package renderer

import (
	"fmt"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/graph"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

func renderMermaid(selected graph.Selected, shape graph.Shape, args Args) (string, error) {
	_ = shape
	nodes := ordering.Notes(selected.Notes)
	rels := ordering.Relationships(selected.Relationships)

	aliasByID := map[string]string{}
	for i, note := range nodes {
		aliasByID[note.ID] = fmt.Sprintf("N%d", i+1)
	}

	var b strings.Builder
	b.WriteString("```mermaid\n")
	b.WriteString("graph ")
	b.WriteString(args.MermaidDirection)
	b.WriteByte('\n')

	for _, note := range nodes {
		alias := aliasByID[note.ID]
		b.WriteString("    ")
		b.WriteString(alias)
		b.WriteString("[")
		b.WriteString(escapeMermaidLabel(note.Title))
		b.WriteString("]\n")
	}

	for _, rel := range rels {
		from, okFrom := aliasByID[rel.FromID]
		to, okTo := aliasByID[rel.ToID]
		if !okFrom || !okTo {
			continue
		}
		b.WriteString("    ")
		b.WriteString(from)
		b.WriteString(" --> ")
		b.WriteString(to)
		b.WriteByte('\n')
	}

	b.WriteString("```\n")
	return b.String(), nil
}

func escapeMermaidLabel(label string) string {
	v := strings.ReplaceAll(label, "\"", "'")
	return strings.TrimSpace(v)
}
