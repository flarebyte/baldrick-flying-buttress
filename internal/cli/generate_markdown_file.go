package cli

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
	"github.com/flarebyte/baldrick-flying-buttress/internal/textutil"
)

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
	return ordering.Strings(textutil.SplitNonEmptyLines(input))
}
