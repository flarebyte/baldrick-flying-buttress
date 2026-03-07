package graph

import (
	"fmt"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

type Shape string

const (
	ShapeTree   Shape = "tree"
	ShapeDAG    Shape = "dag"
	ShapeCyclic Shape = "cyclic"
)

type CyclePolicy string

const (
	CyclePolicyDisallow CyclePolicy = "disallow"
	CyclePolicyAllow    CyclePolicy = "allow"
)

type Query struct {
	SubjectLabel     string
	EdgeLabel        string
	CounterpartLabel string
	StartNode        string
}

type Selected struct {
	Notes         []domain.Note
	Relationships []domain.Relationship
	Roots         []string
}

func Select(app domain.ValidatedApp, query Query) Selected {
	notes := ordering.Notes(app.Notes)
	rels := ordering.Relationships(app.Relationships)
	noteByID := map[string]domain.Note{}
	for _, note := range notes {
		noteByID[note.ID] = note
	}

	subject := map[string]domain.Note{}
	for _, note := range notes {
		if query.SubjectLabel == "" || noteHasLabel(note, query.SubjectLabel) {
			subject[note.ID] = note
		}
	}

	includedNotes := map[string]domain.Note{}
	includedRels := make([]domain.Relationship, 0)
	for _, rel := range rels {
		if query.EdgeLabel != "" && !relationshipHasLabel(rel, query.EdgeLabel) {
			continue
		}

		fromSubject := hasSubject(subject, rel.FromID)
		toSubject := hasSubject(subject, rel.ToID)
		if !fromSubject && !toSubject {
			continue
		}

		if query.CounterpartLabel != "" {
			if fromSubject {
				if toNote, ok := noteByID[rel.ToID]; !ok || !noteHasLabel(toNote, query.CounterpartLabel) {
					fromSubject = false
				}
			}
			if toSubject {
				if fromNote, ok := noteByID[rel.FromID]; !ok || !noteHasLabel(fromNote, query.CounterpartLabel) {
					toSubject = false
				}
			}
			if !fromSubject && !toSubject {
				continue
			}
		}

		includedRels = append(includedRels, rel)
		if note, ok := noteByID[rel.FromID]; ok {
			includedNotes[rel.FromID] = note
		}
		if note, ok := noteByID[rel.ToID]; ok {
			includedNotes[rel.ToID] = note
		}
	}

	for id, note := range subject {
		includedNotes[id] = note
	}

	if strings.TrimSpace(query.StartNode) != "" {
		includedNotes, includedRels = pruneFromStart(includedNotes, includedRels, query.StartNode)
	}

	outNotes := make([]domain.Note, 0, len(includedNotes))
	for _, note := range includedNotes {
		outNotes = append(outNotes, note)
	}
	outNotes = ordering.Notes(outNotes)
	includedRels = ordering.Relationships(includedRels)

	return Selected{
		Notes:         outNotes,
		Relationships: includedRels,
		Roots:         roots(outNotes, includedRels),
	}
}

func ResolveCyclePolicy(arguments []string) (CyclePolicy, error) {
	policy := CyclePolicyDisallow
	for _, arg := range arguments {
		k, v, ok := parseArg(arg)
		if !ok || k != "cycle-policy" {
			continue
		}
		switch v {
		case string(CyclePolicyAllow):
			policy = CyclePolicyAllow
		case string(CyclePolicyDisallow):
			policy = CyclePolicyDisallow
		default:
			return "", fmt.Errorf("invalid cycle-policy: %s", v)
		}
	}
	return policy, nil
}

func QueryFromArgs(arguments []string) Query {
	query := Query{}
	for _, arg := range arguments {
		k, v, ok := parseArg(arg)
		if !ok {
			continue
		}
		switch k {
		case "graph-subject-label":
			query.SubjectLabel = v
		case "graph-edge-label":
			query.EdgeLabel = v
		case "graph-counterpart-label":
			query.CounterpartLabel = v
		case "graph-start-node":
			query.StartNode = v
		}
	}
	return query
}

func HasGraphArgs(arguments []string) bool {
	for _, arg := range arguments {
		k, _, ok := parseArg(arg)
		if !ok {
			continue
		}
		switch k {
		case "graph-subject-label", "graph-edge-label", "graph-counterpart-label", "graph-start-node":
			return true
		}
	}
	return false
}
