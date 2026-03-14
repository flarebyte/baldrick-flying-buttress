package graph

import (
	"fmt"
	"strconv"
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
	IncludeLabels    []string
	ExcludeLabels    []string
	MaxDepth         int
	MaxDepthSet      bool
	ChildOrder       ChildOrder
	BranchPriority   []string
	ShowHelperNodes  bool
	ShowHelpersSet   bool
	HelperLabel      string
}

type Selected struct {
	Notes         []domain.Note
	Relationships []domain.Relationship
	Roots         []string
}

func Select(app domain.ValidatedApp, query Query) Selected {
	query = normalizeQuery(query)
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
	if len(query.IncludeLabels) > 0 || len(query.ExcludeLabels) > 0 {
		includedNotes, includedRels = pruneByLabels(includedNotes, includedRels, query)
	}
	if !query.ShowHelperNodes {
		includedNotes, includedRels = collapseHelperNodes(includedNotes, includedRels, query.HelperLabel)
	}
	if query.MaxDepth >= 0 {
		includedNotes, includedRels = pruneToDepth(includedNotes, includedRels, query.StartNode, query.MaxDepth)
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
	query := Query{MaxDepth: -1, ChildOrder: ChildOrderID, ShowHelperNodes: true}
	_ = applyQueryArgs(&query, arguments, false)
	return normalizeQuery(query)
}

func ResolveQuery(arguments []string) (Query, error) {
	query := Query{MaxDepth: -1, ChildOrder: ChildOrderID, ShowHelperNodes: true}
	if err := applyQueryArgs(&query, arguments, true); err != nil {
		return Query{}, err
	}
	return normalizeQuery(query), nil
}

func applyQueryArgs(query *Query, arguments []string, strict bool) error {
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
		case "graph-include-label":
			query.IncludeLabels = append(query.IncludeLabels, v)
		case "graph-exclude-label":
			query.ExcludeLabels = append(query.ExcludeLabels, v)
		case "graph-max-depth":
			maxDepth, err := strconv.Atoi(v)
			if err != nil || maxDepth < 0 {
				if !strict {
					continue
				}
				return fmt.Errorf("invalid graph-max-depth: %s", v)
			}
			query.MaxDepth = maxDepth
			query.MaxDepthSet = true
		case "graph-child-order":
			if ChildOrder(v) != ChildOrderID && ChildOrder(v) != ChildOrderTitle {
				if !strict {
					query.ChildOrder = ChildOrder(v)
					continue
				}
				return fmt.Errorf("invalid graph-child-order: %s", v)
			}
			query.ChildOrder = ChildOrder(v)
		case "graph-branch-priority-label":
			query.BranchPriority = append(query.BranchPriority, v)
		case "graph-show-helper-nodes":
			showHelpers, err := strconv.ParseBool(v)
			if err != nil {
				if !strict {
					continue
				}
				return fmt.Errorf("invalid graph-show-helper-nodes: %s", v)
			}
			query.ShowHelperNodes = showHelpers
			query.ShowHelpersSet = true
		case "graph-helper-label":
			query.HelperLabel = v
		}
	}
	return nil
}

func HasGraphArgs(arguments []string) bool {
	for _, arg := range arguments {
		k, _, ok := parseArg(arg)
		if !ok {
			continue
		}
		switch k {
		case "graph-subject-label", "graph-edge-label", "graph-counterpart-label", "graph-start-node",
			"graph-include-label", "graph-exclude-label", "graph-max-depth", "graph-child-order",
			"graph-branch-priority-label", "graph-show-helper-nodes", "graph-helper-label":
			return true
		}
	}
	return false
}

func normalizeQuery(query Query) Query {
	query.IncludeLabels = ordering.Strings(compactStrings(query.IncludeLabels))
	query.ExcludeLabels = ordering.Strings(compactStrings(query.ExcludeLabels))
	query.BranchPriority = compactStrings(query.BranchPriority)
	if !query.MaxDepthSet || query.MaxDepth < -1 {
		query.MaxDepth = -1
	}
	switch query.ChildOrder {
	case ChildOrderID, ChildOrderTitle:
	default:
		query.ChildOrder = ChildOrderID
	}
	if !query.ShowHelpersSet {
		query.ShowHelperNodes = true
	}
	if strings.TrimSpace(query.HelperLabel) == "" {
		query.HelperLabel = "helper"
	}
	return query
}

func compactStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
