package cli

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/graph"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
	"github.com/flarebyte/baldrick-flying-buttress/internal/orphans"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	clioutput "github.com/flarebyte/baldrick-flying-buttress/internal/output"
	"github.com/flarebyte/baldrick-flying-buttress/internal/renderer"
)

type generateMarkdownAction struct {
	out io.Writer
}

func (a generateMarkdownAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = a.out
	_ = report
	diagnostics, err := writeMarkdownReports(validated)
	if err != nil {
		return err
	}
	if len(diagnostics) == 0 {
		return nil
	}
	outReport := domain.ValidationReport{Diagnostics: diagnostics}
	if err := clioutput.EmitDiagnostics(a.out, outReport); err != nil {
		return err
	}
	if outReport.HasErrors() {
		return outcome.ValidationBlockedError()
	}
	return nil
}

func (generateMarkdownAction) AllowOnValidationErrors() bool {
	return false
}

func writeMarkdownReports(app domain.ValidatedApp) ([]domain.Diagnostic, error) {
	noteByID := map[string]domain.Note{}
	for _, note := range ordering.Notes(app.Notes) {
		noteByID[note.ID] = note
	}
	diagnostics := make([]domain.Diagnostic, 0)
	registry := renderer.ResolveRegistry()

	for _, report := range ordering.MarkdownReports(app.MarkdownReports) {
		destination := filepath.Join(app.ConfigDir, report.Filepath)
		content, sectionDiagnostics, err := renderMarkdownReport(report, noteByID, app, registry)
		if err != nil {
			return nil, err
		}
		diagnostics = append(diagnostics, sectionDiagnostics...)
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return nil, fmt.Errorf("create report directory %s: %w", filepath.Dir(destination), err)
		}
		if err := os.WriteFile(destination, []byte(content), 0o644); err != nil {
			return nil, fmt.Errorf("write report %s: %w", destination, err)
		}
	}
	return ordering.Diagnostics(diagnostics), nil
}

func renderMarkdownReport(report domain.MarkdownReport, noteByID map[string]domain.Note, app domain.ValidatedApp, registry renderer.Registry) (string, []domain.Diagnostic, error) {
	var b strings.Builder
	diagnostics := make([]domain.Diagnostic, 0)
	writeMarkdownHeading(&b, 1, report.Title)
	writeMarkdownParagraph(&b, report.Description)

	for _, h2 := range ordering.MarkdownH2Sections(report.Sections) {
		writeMarkdownHeading(&b, 2, h2.Title)
		writeMarkdownParagraph(&b, h2.Description)
		for _, h3 := range ordering.MarkdownH3Sections(h2.Sections) {
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
				graphText, err := capability.Render(selected, shape, resolvedArgs)
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
				note, ok := noteByID[noteID]
				if !ok {
					continue
				}
				writeMarkdownHeading(&b, 4, note.Title)
				body, err := renderNoteBody(note, app.ConfigDir)
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

func resolveOrphanQuery(arguments []string) (orphans.Query, bool, error) {
	query := orphans.Query{Direction: orphans.DirectionEither}
	hasOrphanArg := false
	for _, entry := range ordering.Strings(arguments) {
		key, value, ok := parseKVArg(entry)
		if !ok {
			continue
		}
		switch key {
		case "orphan-subject-label":
			hasOrphanArg = true
			query.SubjectLabel = value
		case "orphan-edge-label":
			hasOrphanArg = true
			query.EdgeLabel = value
		case "orphan-counterpart-label":
			hasOrphanArg = true
			query.CounterpartLabel = value
		case "orphan-direction":
			hasOrphanArg = true
			query.Direction = orphans.Direction(value)
		}
	}
	if !hasOrphanArg {
		return orphans.Query{}, false, nil
	}
	if strings.TrimSpace(query.SubjectLabel) == "" {
		return orphans.Query{}, true, fmt.Errorf("orphan-subject-label is required")
	}
	if err := query.Validate(); err != nil {
		return orphans.Query{}, true, err
	}
	return query, true, nil
}

func renderOrphanRows(notes []domain.Note) string {
	var b strings.Builder
	b.WriteString("| name | title | labels |\n")
	b.WriteString("| --- | --- | --- |\n")
	for _, note := range ordering.Notes(notes) {
		labels := strings.Join(splitCSV(note.LabelsCSV), ", ")
		b.WriteString("| ")
		b.WriteString(escapeMarkdownCell(note.ID))
		b.WriteString(" | ")
		b.WriteString(escapeMarkdownCell(note.Title))
		b.WriteString(" | ")
		b.WriteString(escapeMarkdownCell(labels))
		b.WriteString(" |\n")
	}
	return b.String()
}

func parseKVArg(entry string) (string, string, bool) {
	parts := strings.SplitN(entry, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if key == "" || value == "" {
		return "", "", false
	}
	return key, value, true
}

func splitCSV(input string) []string {
	if strings.TrimSpace(input) == "" {
		return []string{}
	}
	items := strings.Split(input, ",")
	out := make([]string, 0, len(items))
	for _, item := range items {
		v := strings.TrimSpace(item)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}

func escapeMarkdownCell(input string) string {
	return strings.ReplaceAll(input, "|", "\\|")
}

type noteArgs struct {
	formatCSV string
	include   csvFilter
	exclude   csvFilter
}

type csvFilter struct {
	column string
	value  string
}

func (f csvFilter) empty() bool {
	return f.column == "" && f.value == ""
}

func renderNoteBody(note domain.Note, configDir string) (string, error) {
	if strings.TrimSpace(note.Filepath) == "" {
		return note.Markdown, nil
	}
	args, err := resolveNoteArgs(note)
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(configDir, note.Filepath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("read note file %s: %w", note.Filepath, err)
	}
	ext := strings.ToLower(filepath.Ext(note.Filepath))
	if ext == ".csv" {
		return renderFileCSV(data, args)
	}
	if isMediaExt(ext) {
		return renderFileMedia(note)
	}
	if isCodeExt(ext) {
		return renderFileCode(data, ext), nil
	}
	return "", fmt.Errorf("unsupported note file type: %s", ext)
}

func resolveNoteArgs(note domain.Note) (noteArgs, error) {
	resolved := noteArgs{formatCSV: "table"}
	for _, arg := range splitArgLines(note.ArgumentsCSV) {
		key, value, ok := parseKVArg(arg)
		if !ok {
			continue
		}
		switch key {
		case "format-csv":
			resolved.formatCSV = value
		case "csv-include":
			filter, err := parseCSVFilter(value)
			if err != nil {
				return noteArgs{}, err
			}
			resolved.include = filter
		case "csv-exclude":
			filter, err := parseCSVFilter(value)
			if err != nil {
				return noteArgs{}, err
			}
			resolved.exclude = filter
		}
	}
	return resolved, nil
}

func parseCSVFilter(value string) (csvFilter, error) {
	parts := strings.SplitN(strings.TrimSpace(value), ":", 2)
	if len(parts) != 2 {
		return csvFilter{}, fmt.Errorf("invalid csv filter: %s", value)
	}
	column := strings.TrimSpace(parts[0])
	target := strings.TrimSpace(parts[1])
	if column == "" || target == "" {
		return csvFilter{}, fmt.Errorf("invalid csv filter: %s", value)
	}
	return csvFilter{column: column, value: target}, nil
}

func renderFileCSV(data []byte, args noteArgs) (string, error) {
	if args.formatCSV == "" {
		args.formatCSV = "table"
	}
	if args.formatCSV != "table" {
		return "", fmt.Errorf("unsupported format-csv: %s", args.formatCSV)
	}
	reader := csv.NewReader(bytes.NewReader(data))
	rows, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("parse csv: %w", err)
	}
	if len(rows) == 0 {
		return "", nil
	}
	header := make([]string, len(rows[0]))
	copy(header, rows[0])
	header = ordering.Strings(header)
	indexByColumn := map[string]int{}
	for i, name := range rows[0] {
		indexByColumn[name] = i
	}
	includeIndex := -1
	excludeIndex := -1
	if !args.include.empty() {
		i, ok := indexByColumn[args.include.column]
		if ok {
			includeIndex = i
		}
	}
	if !args.exclude.empty() {
		i, ok := indexByColumn[args.exclude.column]
		if ok {
			excludeIndex = i
		}
	}
	var b strings.Builder
	b.WriteString("| ")
	b.WriteString(strings.Join(escapeMarkdownCells(header), " | "))
	b.WriteString(" |\n")
	b.WriteString("| ")
	for i := range header {
		if i > 0 {
			b.WriteString(" | ")
		}
		b.WriteString("---")
	}
	b.WriteString(" |\n")

	for _, row := range rows[1:] {
		if includeIndex >= 0 {
			if includeIndex >= len(row) || row[includeIndex] != args.include.value {
				continue
			}
		}
		if excludeIndex >= 0 {
			if excludeIndex < len(row) && row[excludeIndex] == args.exclude.value {
				continue
			}
		}
		byColumn := map[string]string{}
		for i, col := range rows[0] {
			if i < len(row) {
				byColumn[col] = row[i]
			} else {
				byColumn[col] = ""
			}
		}
		values := make([]string, 0, len(header))
		for _, col := range header {
			values = append(values, escapeMarkdownCell(byColumn[col]))
		}
		b.WriteString("| ")
		b.WriteString(strings.Join(values, " | "))
		b.WriteString(" |\n")
	}
	return strings.TrimSuffix(b.String(), "\n"), nil
}

func escapeMarkdownCells(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		out = append(out, escapeMarkdownCell(v))
	}
	return out
}

func renderFileMedia(note domain.Note) (string, error) {
	title := strings.TrimSpace(note.Title)
	if title == "" {
		title = note.ID
	}
	return fmt.Sprintf("![%s](%s)", title, filepath.ToSlash(note.Filepath)), nil
}

func renderFileCode(data []byte, ext string) string {
	lang := codeLanguage(ext)
	var b strings.Builder
	b.WriteString("```")
	b.WriteString(lang)
	b.WriteByte('\n')
	b.Write(data)
	if len(data) == 0 || data[len(data)-1] != '\n' {
		b.WriteByte('\n')
	}
	b.WriteString("```")
	return b.String()
}

func codeLanguage(ext string) string {
	switch ext {
	case ".md":
		return "markdown"
	case ".mmd", ".mermaid":
		return "mermaid"
	case ".puml", ".plantuml":
		return "plantuml"
	case ".yaml", ".yml":
		return "yaml"
	default:
		if strings.HasPrefix(ext, ".") {
			return ext[1:]
		}
		return ext
	}
}

func isMediaExt(ext string) bool {
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp":
		return true
	default:
		return false
	}
}

func isCodeExt(ext string) bool {
	switch ext {
	case ".go", ".ts", ".js", ".json", ".md", ".txt", ".mmd", ".mermaid", ".puml", ".plantuml", ".yaml", ".yml", ".sql", ".cue":
		return true
	default:
		return false
	}
}

func splitArgLines(input string) []string {
	if strings.TrimSpace(input) == "" {
		return []string{}
	}
	lines := strings.Split(input, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		v := strings.TrimSpace(line)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return ordering.Strings(out)
}
