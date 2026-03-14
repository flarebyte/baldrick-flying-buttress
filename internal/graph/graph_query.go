package graph

import (
	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

type ChildOrder string

const (
	ChildOrderID    ChildOrder = "id"
	ChildOrderTitle ChildOrder = "title"
)

func pruneByLabels(notes map[string]domain.Note, rels []domain.Relationship, query Query) (map[string]domain.Note, []domain.Relationship) {
	if len(query.IncludeLabels) == 0 && len(query.ExcludeLabels) == 0 {
		return notes, rels
	}

	filteredNotes := map[string]domain.Note{}
	for id, note := range notes {
		if noteAllowedByQueryLabels(note, query) {
			filteredNotes[id] = note
		}
	}

	return filteredNotes, filterRelationshipsWithinNotes(filteredNotes, rels)
}

func noteAllowedByQueryLabels(note domain.Note, query Query) bool {
	if len(query.IncludeLabels) > 0 {
		matched := false
		for _, label := range query.IncludeLabels {
			if noteHasLabel(note, label) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	for _, label := range query.ExcludeLabels {
		if noteHasLabel(note, label) {
			return false
		}
	}
	return true
}

func collapseHelperNodes(notes map[string]domain.Note, rels []domain.Relationship, helperLabel string) (map[string]domain.Note, []domain.Relationship) {
	helperIDs := map[string]struct{}{}
	for id, note := range notes {
		if noteHasLabel(note, helperLabel) {
			helperIDs[id] = struct{}{}
		}
	}
	if len(helperIDs) == 0 {
		return notes, rels
	}

	parents := map[string][]string{}
	children := map[string][]string{}
	for _, rel := range rels {
		children[rel.FromID] = append(children[rel.FromID], rel.ToID)
		parents[rel.ToID] = append(parents[rel.ToID], rel.FromID)
	}

	rewired := map[string]domain.Relationship{}
	for _, rel := range rels {
		if _, isHelperFrom := helperIDs[rel.FromID]; isHelperFrom {
			continue
		}
		if _, isHelperTo := helperIDs[rel.ToID]; isHelperTo {
			continue
		}
		rewired[rel.FromID+"->"+rel.ToID] = rel
	}

	for helperID := range helperIDs {
		parentIDs := collapseNonHelperAncestors(helperID, parents, helperIDs)
		childIDs := collapseNonHelperDescendants(helperID, children, helperIDs)
		for _, parentID := range parentIDs {
			for _, childID := range childIDs {
				if parentID == childID {
					continue
				}
				key := parentID + "->" + childID
				rewired[key] = domain.Relationship{FromID: parentID, ToID: childID}
			}
		}
	}

	filteredNotes := map[string]domain.Note{}
	for id, note := range notes {
		if _, isHelper := helperIDs[id]; isHelper {
			continue
		}
		filteredNotes[id] = note
	}

	return filteredNotes, filterRelationshipsWithinNotes(filteredNotes, mapRelationships(rewired))
}

func filterRelationshipsWithinNotes(notes map[string]domain.Note, rels []domain.Relationship) []domain.Relationship {
	filtered := make([]domain.Relationship, 0, len(rels))
	for _, rel := range rels {
		if _, ok := notes[rel.FromID]; !ok {
			continue
		}
		if _, ok := notes[rel.ToID]; !ok {
			continue
		}
		filtered = append(filtered, rel)
	}
	return filtered
}

func mapRelationships(input map[string]domain.Relationship) []domain.Relationship {
	out := make([]domain.Relationship, 0, len(input))
	for _, rel := range input {
		out = append(out, rel)
	}
	return out
}

func collapseNonHelperAncestors(nodeID string, parents map[string][]string, helperIDs map[string]struct{}) []string {
	return collapseBoundaryNodes(nodeID, parents, helperIDs)
}

func collapseNonHelperDescendants(nodeID string, children map[string][]string, helperIDs map[string]struct{}) []string {
	return collapseBoundaryNodes(nodeID, children, helperIDs)
}

func collapseBoundaryNodes(nodeID string, links map[string][]string, helperIDs map[string]struct{}) []string {
	seen := map[string]struct{}{}
	stack := append([]string(nil), links[nodeID]...)
	out := make([]string, 0)
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if _, ok := seen[current]; ok {
			continue
		}
		seen[current] = struct{}{}
		if _, isHelper := helperIDs[current]; !isHelper {
			out = append(out, current)
			continue
		}
		stack = append(stack, links[current]...)
	}
	return ordering.Strings(out)
}

func pruneToDepth(notes map[string]domain.Note, rels []domain.Relationship, startNode string, maxDepth int) (map[string]domain.Note, []domain.Relationship) {
	if maxDepth < 0 {
		return notes, rels
	}

	outNotes := map[string]domain.Note{}
	if len(notes) == 0 {
		return outNotes, nil
	}

	noteSlice := make([]domain.Note, 0, len(notes))
	for _, note := range notes {
		noteSlice = append(noteSlice, note)
	}
	adj := buildAdj(rels)
	startIDs := roots(noteSlice, rels)
	if startNode != "" {
		if _, ok := notes[startNode]; ok {
			startIDs = []string{startNode}
		}
	}
	if len(startIDs) == 0 {
		startIDs = noteIDs(noteSlice)
	}

	depthByID := map[string]int{}
	queue := append([]string(nil), startIDs...)
	for _, startID := range startIDs {
		depthByID[startID] = 0
	}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		depth := depthByID[current]
		if depth >= maxDepth {
			continue
		}
		for _, childID := range adj[current] {
			nextDepth := depth + 1
			if prevDepth, ok := depthByID[childID]; ok && prevDepth <= nextDepth {
				continue
			}
			depthByID[childID] = nextDepth
			queue = append(queue, childID)
		}
	}

	for id, note := range notes {
		depth, ok := depthByID[id]
		if !ok || depth > maxDepth {
			continue
		}
		outNotes[id] = note
	}

	outRels := make([]domain.Relationship, 0, len(rels))
	for _, rel := range rels {
		fromDepth, okFrom := depthByID[rel.FromID]
		toDepth, okTo := depthByID[rel.ToID]
		if !okFrom || !okTo || fromDepth >= maxDepth || toDepth > maxDepth {
			continue
		}
		if _, ok := outNotes[rel.FromID]; !ok {
			continue
		}
		if _, ok := outNotes[rel.ToID]; !ok {
			continue
		}
		outRels = append(outRels, rel)
	}
	return outNotes, outRels
}
