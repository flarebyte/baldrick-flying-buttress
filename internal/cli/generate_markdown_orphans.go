package cli

import (
	"fmt"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
	"github.com/flarebyte/baldrick-flying-buttress/internal/orphans"
	"github.com/flarebyte/baldrick-flying-buttress/internal/textutil"
)

func resolveOrphanQuery(arguments []string) (orphans.Query, bool, error) {
	query := orphans.Query{Direction: orphans.DirectionEither}
	hasOrphanArg := false
	for _, entry := range ordering.Strings(arguments) {
		key, value, ok := textutil.ParseKeyValue(entry)
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
		labels := strings.Join(textutil.SplitCSV(note.LabelsCSV), ", ")
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
