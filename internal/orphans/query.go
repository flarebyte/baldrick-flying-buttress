package orphans

import (
	"fmt"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

type Direction string

const (
	DirectionIn     Direction = "in"
	DirectionOut    Direction = "out"
	DirectionEither Direction = "either"
)

type Query struct {
	SubjectLabel     string
	EdgeLabel        string
	CounterpartLabel string
	Direction        Direction
}

func (q Query) Validate() error {
	if q.SubjectLabel == "" {
		return fmt.Errorf("subject-label is required")
	}
	switch q.Direction {
	case DirectionIn, DirectionOut, DirectionEither:
		return nil
	default:
		return fmt.Errorf("invalid direction: %s", q.Direction)
	}
}

func Find(app domain.ValidatedApp, query Query) []domain.Note {
	orderedNotes := ordering.Notes(app.Notes)
	orderedRelationships := ordering.Relationships(app.Relationships)
	noteByID := make(map[string]domain.Note, len(orderedNotes))
	for _, note := range orderedNotes {
		noteByID[note.ID] = note
	}

	orphans := make([]domain.Note, 0)
	for _, note := range orderedNotes {
		if note.Label != query.SubjectLabel {
			continue
		}
		if hasMatchingRelationship(note, orderedRelationships, noteByID, query) {
			continue
		}
		orphans = append(orphans, note)
	}
	return orphans
}

func hasMatchingRelationship(subject domain.Note, relationships []domain.Relationship, noteByID map[string]domain.Note, query Query) bool {
	for _, relationship := range relationships {
		if query.EdgeLabel != "" && relationship.Label != query.EdgeLabel {
			continue
		}

		counterpartID, ok := relationshipMatchesDirection(relationship, subject.ID, query.Direction)
		if !ok {
			continue
		}

		if query.CounterpartLabel != "" {
			counterpart, exists := noteByID[counterpartID]
			if !exists || counterpart.Label != query.CounterpartLabel {
				continue
			}
		}

		return true
	}
	return false
}

func relationshipMatchesDirection(relationship domain.Relationship, subjectID string, direction Direction) (string, bool) {
	switch direction {
	case DirectionIn:
		if relationship.ToID == subjectID {
			return relationship.FromID, true
		}
		return "", false
	case DirectionOut:
		if relationship.FromID == subjectID {
			return relationship.ToID, true
		}
		return "", false
	case DirectionEither:
		if relationship.FromID == subjectID {
			return relationship.ToID, true
		}
		if relationship.ToID == subjectID {
			return relationship.FromID, true
		}
		return "", false
	default:
		return "", false
	}
}
