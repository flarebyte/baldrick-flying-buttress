package renderer

import (
	"fmt"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/graph"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

type Args struct {
	Renderer         string
	MermaidDirection string
}

type Capability struct {
	Name            string
	SupportedShapes []graph.Shape
	SupportedArgs   []string
	Render          func(graph.Selected, graph.Shape, Args) (string, error)
}

type Registry struct {
	Capabilities []Capability
}

func ResolveRegistry() Registry {
	caps := []Capability{
		{
			Name:            "markdown-text",
			SupportedShapes: []graph.Shape{graph.ShapeTree, graph.ShapeDAG, graph.ShapeCyclic},
			SupportedArgs:   []string{"graph-renderer"},
			Render: func(selected graph.Selected, shape graph.Shape, args Args) (string, error) {
				_ = args
				return graph.RenderMarkdownText(selected, shape, graph.CyclePolicyAllow), nil
			},
		},
		{
			Name:            "mermaid",
			SupportedShapes: []graph.Shape{graph.ShapeTree, graph.ShapeDAG, graph.ShapeCyclic},
			SupportedArgs:   []string{"graph-renderer", "mermaid-direction"},
			Render:          renderMermaid,
		},
	}
	return Registry{Capabilities: caps}
}

func (r Registry) Select(name string, shape graph.Shape) (Capability, error) {
	target := strings.TrimSpace(name)
	for _, capability := range r.Capabilities {
		if capability.Name != target {
			continue
		}
		for _, supported := range capability.SupportedShapes {
			if supported == shape {
				return capability, nil
			}
		}
		return Capability{}, fmt.Errorf("renderer %s does not support shape %s", target, shape)
	}
	return Capability{}, fmt.Errorf("unsupported renderer: %s", target)
}

func ResolveArgs(app domain.ValidatedApp, h3 domain.MarkdownH3Section, noteByID map[string]domain.Note) (Args, error) {
	resolved := Args{Renderer: "markdown-text", MermaidDirection: "TD"}

	for _, def := range app.Registry.Arguments {
		if def.DefaultValue == nil {
			continue
		}
		switch def.Name {
		case "graph-renderer":
			resolved.Renderer = strings.TrimSpace(*def.DefaultValue)
		case "mermaid-direction":
			resolved.MermaidDirection = strings.TrimSpace(*def.DefaultValue)
		}
	}

	applyArgs := func(arguments []string) {
		for _, entry := range ordering.Strings(arguments) {
			key, value, ok := parseArg(entry)
			if !ok {
				continue
			}
			switch key {
			case "graph-renderer":
				resolved.Renderer = value
			case "mermaid-direction":
				resolved.MermaidDirection = value
			}
		}
	}

	applyArgs(h3.Arguments)

	for _, noteID := range ordering.Strings(h3.NoteIDs) {
		note, ok := noteByID[noteID]
		if !ok {
			continue
		}
		applyArgs(splitArgs(note.ArgumentsCSV))
	}

	if resolved.Renderer == "" {
		resolved.Renderer = "markdown-text"
	}
	if resolved.MermaidDirection == "" {
		resolved.MermaidDirection = "TD"
	}
	if resolved.MermaidDirection != "TD" && resolved.MermaidDirection != "LR" {
		return Args{}, fmt.Errorf("invalid mermaid-direction: %s", resolved.MermaidDirection)
	}
	return resolved, nil
}

func parseArg(entry string) (string, string, bool) {
	parts := strings.SplitN(entry, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	k := strings.TrimSpace(parts[0])
	v := strings.TrimSpace(parts[1])
	if k == "" || v == "" {
		return "", "", false
	}
	return k, v, true
}

func splitArgs(csv string) []string {
	if strings.TrimSpace(csv) == "" {
		return []string{}
	}
	parts := strings.Split(csv, "\n")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		v := strings.TrimSpace(part)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}
