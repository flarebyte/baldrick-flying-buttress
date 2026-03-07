package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/graph"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
	"github.com/flarebyte/baldrick-flying-buttress/internal/orphans"
	"github.com/flarebyte/baldrick-flying-buttress/internal/renderer"
)

func renderMarkdownReport(ctx context.Context, report domain.MarkdownReport, noteByID map[string]domain.Note, app domain.ValidatedApp, registry renderer.Registry) (string, []domain.Diagnostic, error) {
	var b strings.Builder
	diagnostics := make([]domain.Diagnostic, 0)
	writeMarkdownHeading(&b, 1, report.Title)
	writeMarkdownParagraph(&b, report.Description)

	for _, h2 := range ordering.MarkdownH2Sections(report.Sections) {
		if err := ctx.Err(); err != nil {
			return "", nil, err
		}
		writeMarkdownHeading(&b, 2, h2.Title)
		writeMarkdownParagraph(&b, h2.Description)
		for _, h3 := range ordering.MarkdownH3Sections(h2.Sections) {
			if err := ctx.Err(); err != nil {
				return "", nil, err
			}
			writeMarkdownHeading(&b, 3, h3.Title)
			writeMarkdownParagraph(&b, h3.Description)
			orphanQuery, orphanMode, err := resolveOrphanQuery(h3.Arguments)
			if err != nil {
				diagnostics = append(diagnostics, domain.Diagnostic{
					Code:         "ORPHAN_QUERY_ARGS_INVALID",
					Severity:     domain.SeverityError,
					Source:       "args.orphan.query.resolve",
					Message:      err.Error(),
					Location:     h3.Path,
					Path:         h3.Path,
					ReportTitle:  h3.ReportTitle,
					SectionTitle: h3.H2Title,
				})
				continue
			}
			if orphanMode {
				orphanNotes := orphans.Find(app, orphanQuery)
				table := renderOrphanRows(orphanNotes)
				if table != "" {
					b.WriteString(table)
					if !strings.HasSuffix(table, "\n") {
						b.WriteByte('\n')
					}
					b.WriteByte('\n')
				}
				continue
			}
			if graph.HasGraphArgs(h3.Arguments) {
				query := graph.QueryFromArgs(h3.Arguments)
				selected := graph.Select(app, query)
				shape := graph.DetectShape(selected)
				cyclePolicy, err := graph.ResolveCyclePolicy(h3.Arguments)
				if err != nil {
					diagnostics = append(diagnostics, domain.Diagnostic{
						Code:         "GRAPH_CYCLE_POLICY_INVALID",
						Severity:     domain.SeverityError,
						Source:       "graph.policy.cycle",
						Message:      err.Error(),
						Location:     h3.Path,
						Path:         h3.Path,
						ReportTitle:  h3.ReportTitle,
						SectionTitle: h3.H2Title,
					})
					continue
				}
				if shape == graph.ShapeCyclic && cyclePolicy == graph.CyclePolicyDisallow {
					diagnostics = append(diagnostics, domain.Diagnostic{
						Code:         "GRAPH_CYCLE_DISALLOWED",
						Severity:     domain.SeverityWarning,
						Source:       "graph.policy.cycle",
						Message:      "cycle detected while cycle-policy=disallow; skipping graph section",
						Location:     h3.Path,
						Path:         h3.Path,
						ReportTitle:  h3.ReportTitle,
						SectionTitle: h3.H2Title,
					})
					continue
				}
				resolvedArgs, err := renderer.ResolveArgs(app, h3, noteByID)
				if err != nil {
					diagnostics = append(diagnostics, domain.Diagnostic{
						Code:         "GRAPH_RENDERER_ARGS_INVALID",
						Severity:     domain.SeverityError,
						Source:       "args.renderer.resolve",
						Message:      err.Error(),
						Location:     h3.Path,
						Path:         h3.Path,
						ReportTitle:  h3.ReportTitle,
						SectionTitle: h3.H2Title,
					})
					continue
				}
				capability, err := registry.Select(resolvedArgs.Renderer, shape)
				if err != nil {
					diagnostics = append(diagnostics, domain.Diagnostic{
						Code:         "GRAPH_RENDERER_UNSUPPORTED",
						Severity:     domain.SeverityError,
						Source:       "renderer.registry.resolve",
						Message:      err.Error(),
						Location:     h3.Path,
						Path:         h3.Path,
						ReportTitle:  h3.ReportTitle,
						SectionTitle: h3.H2Title,
					})
					continue
				}
				graphText, err := capability.Render(ctx, selected, shape, resolvedArgs)
				if err != nil {
					diagnostics = append(diagnostics, domain.Diagnostic{
						Code:         "GRAPH_RENDER_FAILED",
						Severity:     domain.SeverityError,
						Source:       "render.graph." + capability.Name,
						Message:      err.Error(),
						Location:     h3.Path,
						Path:         h3.Path,
						ReportTitle:  h3.ReportTitle,
						SectionTitle: h3.H2Title,
					})
					continue
				}
				if graphText != "" {
					b.WriteString(graphText)
					if !strings.HasSuffix(graphText, "\n") {
						b.WriteByte('\n')
					}
					b.WriteByte('\n')
				}
				continue
			}
			for _, noteID := range ordering.Strings(h3.NoteIDs) {
				if err := ctx.Err(); err != nil {
					return "", nil, err
				}
				note, ok := noteByID[noteID]
				if !ok {
					continue
				}
				writeMarkdownHeading(&b, 4, note.Title)
				body, err := renderNoteBody(ctx, note, app.ConfigDir)
				if err != nil {
					return "", nil, fmt.Errorf("render note %s: %w", note.ID, err)
				}
				writeMarkdownParagraph(&b, body)
			}
		}
	}

	if b.Len() == 0 || !strings.HasSuffix(b.String(), "\n") {
		b.WriteByte('\n')
	}
	return b.String(), diagnostics, nil
}

func writeMarkdownHeading(b *strings.Builder, level int, title string) {
	if strings.TrimSpace(title) == "" {
		return
	}
	b.WriteString(strings.Repeat("#", level))
	b.WriteByte(' ')
	b.WriteString(title)
	b.WriteString("\n\n")
}

func writeMarkdownParagraph(b *strings.Builder, text string) {
	if strings.TrimSpace(text) == "" {
		return
	}
	b.WriteString(text)
	b.WriteString("\n\n")
}

func escapeMarkdownCell(input string) string {
	return strings.ReplaceAll(input, "|", "\\|")
}
