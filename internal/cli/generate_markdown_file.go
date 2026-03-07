package cli

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
	"github.com/flarebyte/baldrick-flying-buttress/internal/safety"
	"github.com/flarebyte/baldrick-flying-buttress/internal/textutil"
)

type noteArgs struct {
	formatCSV string
	includes  []csvFilter
	excludes  []csvFilter
}

type csvFilter struct {
	column string
	value  string
}

func renderNoteBody(ctx context.Context, note domain.Note, configDir string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
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
	if err := ctx.Err(); err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(note.Filepath))
	if ext == ".csv" {
		return renderFileCSV(ctx, data, args)
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
		key, value, ok := textutil.ParseKeyValue(arg)
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
			resolved.includes = append(resolved.includes, filter)
		case "csv-exclude":
			filter, err := parseCSVFilter(value)
			if err != nil {
				return noteArgs{}, err
			}
			resolved.excludes = append(resolved.excludes, filter)
		}
	}
	resolved.includes = sortCSVFilters(resolved.includes)
	resolved.excludes = sortCSVFilters(resolved.excludes)
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

func renderFileCSV(ctx context.Context, data []byte, args noteArgs) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if err := safety.CheckCSVFileSize(len(data)); err != nil {
		return "", err
	}
	if args.formatCSV == "" {
		args.formatCSV = "table"
	}
	if args.formatCSV != "table" && args.formatCSV != "raw" {
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
	sourceHeader := make([]string, len(rows[0]))
	copy(sourceHeader, rows[0])
	sortedHeader := ordering.Strings(append([]string(nil), sourceHeader...))
	indexByColumn := make(map[string]int, len(sourceHeader))
	for i, name := range sourceHeader {
		if _, exists := indexByColumn[name]; exists {
			continue
		}
		indexByColumn[name] = i
	}
	includeChecks := make([]csvFilterCheck, 0, len(args.includes))
	excludeChecks := make([]csvFilterCheck, 0, len(args.excludes))
	includeImpossible := false
	for _, filter := range args.includes {
		i, ok := indexByColumn[filter.column]
		if !ok {
			includeImpossible = true
			break
		}
		includeChecks = append(includeChecks, csvFilterCheck{index: i, value: filter.value})
	}
	for _, filter := range args.excludes {
		i, ok := indexByColumn[filter.column]
		if !ok {
			continue
		}
		excludeChecks = append(excludeChecks, csvFilterCheck{index: i, value: filter.value})
	}
	filteredRows := make([][]string, 0, len(rows))
	if !includeImpossible {
		for _, row := range rows[1:] {
			if err := ctx.Err(); err != nil {
				return "", err
			}
			if !matchesAllFilters(row, includeChecks) || matchesAnyFilter(row, excludeChecks) {
				continue
			}
			filteredRows = append(filteredRows, row)
			if err := safety.CheckCSVRenderedRows(len(filteredRows)); err != nil {
				return "", err
			}
		}
	}
	if args.formatCSV == "raw" {
		return renderRawCSV(sourceHeader, filteredRows)
	}
	var b strings.Builder
	b.WriteString("| ")
	b.WriteString(strings.Join(escapeMarkdownCells(sortedHeader), " | "))
	b.WriteString(" |\n")
	b.WriteString("| ")
	for i := range sortedHeader {
		if i > 0 {
			b.WriteString(" | ")
		}
		b.WriteString("---")
	}
	b.WriteString(" |\n")
	for _, row := range filteredRows {
		values := make([]string, 0, len(sortedHeader))
		for _, col := range sortedHeader {
			values = append(values, escapeMarkdownTableCell(valueAtColumn(row, indexByColumn[col])))
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
		out = append(out, escapeMarkdownTableCell(v))
	}
	return out
}

type csvFilterCheck struct {
	index int
	value string
}

func sortCSVFilters(filters []csvFilter) []csvFilter {
	if len(filters) == 0 {
		return nil
	}
	sorted := append([]csvFilter(nil), filters...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].column == sorted[j].column {
			return sorted[i].value < sorted[j].value
		}
		return sorted[i].column < sorted[j].column
	})
	return sorted
}

func matchesAllFilters(row []string, checks []csvFilterCheck) bool {
	for _, check := range checks {
		if check.index >= len(row) || row[check.index] != check.value {
			return false
		}
	}
	return true
}

func matchesAnyFilter(row []string, checks []csvFilterCheck) bool {
	for _, check := range checks {
		if check.index < len(row) && row[check.index] == check.value {
			return true
		}
	}
	return false
}

func valueAtColumn(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return row[index]
}

func escapeMarkdownTableCell(input string) string {
	normalized := strings.ReplaceAll(input, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	normalized = strings.ReplaceAll(normalized, "|", "\\|")
	return strings.ReplaceAll(normalized, "\n", "<br/>")
}

func renderRawCSV(header []string, rows [][]string) (string, error) {
	var b strings.Builder
	b.WriteString("```csv\n")
	writer := csv.NewWriter(&b)
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("render raw csv header: %w", err)
	}
	for _, row := range rows {
		normalized := make([]string, len(header))
		copy(normalized, row)
		if err := writer.Write(normalized); err != nil {
			return "", fmt.Errorf("render raw csv row: %w", err)
		}
	}
	writer.Flush()
	if writer.Error() != nil {
		return "", fmt.Errorf("render raw csv flush: %w", writer.Error())
	}
	b.WriteString("```")
	return b.String(), nil
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
	return ordering.Strings(textutil.SplitNonEmptyLines(input))
}
