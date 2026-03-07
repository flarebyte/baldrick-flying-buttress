package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

const (
	namesKindAll           = "all"
	namesKindNotes         = "notes"
	namesKindRelationships = "relationships"
	namesFormatTable       = "table"
	namesFormatJSON        = "json"
)

type namesAction struct {
	out    io.Writer
	prefix string
	kind   string
	format string
}

func (a namesAction) Execute(_ context.Context, validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = report
	notes, relationships, err := filterNames(validated, a.prefix, a.kind)
	if err != nil {
		return err
	}
	switch a.format {
	case namesFormatTable:
		return emitNamesTable(a.out, notes, relationships)
	case namesFormatJSON:
		return emitNamesJSON(a.out, notes, relationships)
	default:
		return fmt.Errorf("invalid format: %s", a.format)
	}
}

func (namesAction) AllowOnValidationErrors() bool {
	return false
}

func validateNamesKind(kind string) error {
	switch kind {
	case namesKindAll, namesKindNotes, namesKindRelationships:
		return nil
	default:
		return fmt.Errorf("invalid kind: %s", kind)
	}
}

func validateNamesFormat(format string) error {
	switch format {
	case namesFormatTable, namesFormatJSON:
		return nil
	default:
		return fmt.Errorf("invalid format: %s", format)
	}
}

func filterNames(app domain.ValidatedApp, prefix string, kind string) ([]domain.Note, []domain.Relationship, error) {
	if prefix == "" {
		return nil, nil, errors.New("prefix is required")
	}
	if err := validateNamesKind(kind); err != nil {
		return nil, nil, err
	}

	orderedNotes := ordering.Notes(app.Notes)
	orderedRelationships := ordering.Relationships(app.Relationships)

	notes := make([]domain.Note, 0)
	relationships := make([]domain.Relationship, 0)

	if kind == namesKindAll || kind == namesKindNotes {
		for _, note := range orderedNotes {
			if matchesNotePrefix(note.ID, prefix) {
				notes = append(notes, note)
			}
		}
	}

	if kind == namesKindAll || kind == namesKindRelationships {
		for _, relationship := range orderedRelationships {
			if matchesRelationshipPrefix(relationship, prefix) {
				relationships = append(relationships, relationship)
			}
		}
	}

	return notes, relationships, nil
}

func matchesNotePrefix(name string, prefix string) bool {
	return prefix == "" || strings.HasPrefix(name, prefix)
}

func matchesRelationshipPrefix(relationship domain.Relationship, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(relationship.FromID, prefix) || strings.HasPrefix(relationship.ToID, prefix)
}

func emitNamesTable(w io.Writer, notes []domain.Note, relationships []domain.Relationship) error {
	writeDivider := false
	if len(notes) > 0 {
		rows := make([][]string, 0, len(notes))
		for _, note := range notes {
			rows = append(rows, []string{
				note.ID,
				note.Title,
				joinSortedLabels(note.LabelsCSV),
			})
		}
		if err := emitAlignedTable(w, "notes", []string{"name", "title", "labels"}, rows); err != nil {
			return err
		}
		writeDivider = true
	}
	if len(relationships) > 0 {
		if writeDivider {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}
		rows := make([][]string, 0, len(relationships))
		for _, relationship := range relationships {
			rows = append(rows, []string{
				relationship.FromID,
				relationship.ToID,
				joinSortedLabels(relationship.LabelsCSV),
			})
		}
		if err := emitAlignedTable(w, "relationships", []string{"from", "to", "labels"}, rows); err != nil {
			return err
		}
	}
	return nil
}

func emitAlignedTable(w io.Writer, title string, headers []string, rows [][]string) error {
	if len(headers) == 0 {
		return nil
	}
	widths := make([]int, len(headers))
	for i := range headers {
		widths[i] = len(headers[i])
	}
	for _, row := range rows {
		for i := range headers {
			if i >= len(row) {
				continue
			}
			if len(row[i]) > widths[i] {
				widths[i] = len(row[i])
			}
		}
	}
	if _, err := io.WriteString(w, title+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(w, formatAlignedTableRow(headers, widths)); err != nil {
		return err
	}
	dividerCells := make([]string, len(headers))
	for i := range widths {
		dividerCells[i] = strings.Repeat("-", widths[i])
	}
	if _, err := io.WriteString(w, formatAlignedTableRow(dividerCells, widths)); err != nil {
		return err
	}
	for _, row := range rows {
		cells := make([]string, len(headers))
		copy(cells, row)
		if _, err := io.WriteString(w, formatAlignedTableRow(cells, widths)); err != nil {
			return err
		}
	}
	return nil
}

func joinSortedLabels(csv string) string {
	labels := splitCSVLabels(csv)
	if len(labels) == 0 {
		return "-"
	}
	sort.Strings(labels)
	return strings.Join(labels, ", ")
}

func splitCSVLabels(csv string) []string {
	if csv == "" {
		return nil
	}
	parts := strings.Split(csv, ",")
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

type namesJSONDTO struct {
	Notes         []namesJSONNoteDTO         `json:"notes"`
	Relationships []namesJSONRelationshipDTO `json:"relationships"`
}

type namesJSONNoteDTO struct {
	Name string `json:"name"`
}

type namesJSONRelationshipDTO struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func emitNamesJSON(w io.Writer, notes []domain.Note, relationships []domain.Relationship) error {
	dto := namesJSONDTO{
		Notes:         make([]namesJSONNoteDTO, 0, len(notes)),
		Relationships: make([]namesJSONRelationshipDTO, 0, len(relationships)),
	}
	for _, note := range notes {
		dto.Notes = append(dto.Notes, namesJSONNoteDTO{Name: note.ID})
	}
	for _, relationship := range relationships {
		dto.Relationships = append(dto.Relationships, namesJSONRelationshipDTO{From: relationship.FromID, To: relationship.ToID})
	}
	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	_, err = w.Write(append(data, '\n'))
	return err
}
