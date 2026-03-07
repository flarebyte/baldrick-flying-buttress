package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (a namesAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
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
			if strings.HasPrefix(note.ID, prefix) {
				notes = append(notes, note)
			}
		}
	}

	if kind == namesKindAll || kind == namesKindRelationships {
		for _, relationship := range orderedRelationships {
			if strings.HasPrefix(relationship.FromID, prefix) || strings.HasPrefix(relationship.ToID, prefix) {
				relationships = append(relationships, relationship)
			}
		}
	}

	return notes, relationships, nil
}

func emitNamesTable(w io.Writer, notes []domain.Note, relationships []domain.Relationship) error {
	if _, err := io.WriteString(w, "KIND\tNAME\tFROM\tTO\n"); err != nil {
		return err
	}
	for _, note := range notes {
		if _, err := fmt.Fprintf(w, "note\t%s\t\t\n", note.ID); err != nil {
			return err
		}
	}
	for _, relationship := range relationships {
		if _, err := fmt.Fprintf(w, "relationship\t\t%s\t%s\n", relationship.FromID, relationship.ToID); err != nil {
			return err
		}
	}
	return nil
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
