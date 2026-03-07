package graph

import (
	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

func DetectShape(selected Selected) Shape {
	adj := buildAdj(selected.Relationships)
	if hasCycle(selected.Notes, adj) {
		return ShapeCyclic
	}
	if isTree(selected.Notes, selected.Relationships, adj) {
		return ShapeTree
	}
	return ShapeDAG
}

func buildAdj(rels []domain.Relationship) map[string][]string {
	adj := map[string][]string{}
	for _, rel := range ordering.Relationships(rels) {
		adj[rel.FromID] = append(adj[rel.FromID], rel.ToID)
	}
	for key := range adj {
		adj[key] = ordering.Strings(adj[key])
	}
	return adj
}

func roots(notes []domain.Note, rels []domain.Relationship) []string {
	indegree := map[string]int{}
	for _, note := range notes {
		indegree[note.ID] = 0
	}
	for _, rel := range rels {
		if _, ok := indegree[rel.ToID]; ok {
			indegree[rel.ToID]++
		}
	}
	rootIDs := make([]string, 0)
	for id, deg := range indegree {
		if deg == 0 {
			rootIDs = append(rootIDs, id)
		}
	}
	return ordering.Strings(rootIDs)
}

func noteIDs(notes []domain.Note) []string {
	ids := make([]string, 0, len(notes))
	for _, note := range notes {
		ids = append(ids, note.ID)
	}
	return ordering.Strings(ids)
}

func hasCycle(notes []domain.Note, adj map[string][]string) bool {
	state := map[string]int{}
	for _, note := range notes {
		if state[note.ID] != 0 {
			continue
		}
		if visitCycle(note.ID, adj, state) {
			return true
		}
	}
	return false
}

func visitCycle(node string, adj map[string][]string, state map[string]int) bool {
	state[node] = 1
	for _, child := range adj[node] {
		if state[child] == 1 {
			return true
		}
		if state[child] == 0 && visitCycle(child, adj, state) {
			return true
		}
	}
	state[node] = 2
	return false
}

func isTree(notes []domain.Note, rels []domain.Relationship, adj map[string][]string) bool {
	if len(notes) == 0 {
		return true
	}
	if len(rels) != len(notes)-1 {
		return false
	}
	indegree := map[string]int{}
	for _, note := range notes {
		indegree[note.ID] = 0
	}
	for _, rel := range rels {
		if _, ok := indegree[rel.ToID]; ok {
			indegree[rel.ToID]++
			if indegree[rel.ToID] > 1 {
				return false
			}
		}
	}
	rootCount := 0
	rootID := ""
	for id, deg := range indegree {
		if deg == 0 {
			rootCount++
			rootID = id
		}
	}
	if rootCount != 1 {
		return false
	}
	seen := map[string]bool{}
	walk(rootID, adj, seen)
	return len(seen) == len(notes)
}

func walk(node string, adj map[string][]string, seen map[string]bool) {
	if seen[node] {
		return
	}
	seen[node] = true
	for _, child := range adj[node] {
		walk(child, adj, seen)
	}
}

func pruneFromStart(notes map[string]domain.Note, rels []domain.Relationship, start string) (map[string]domain.Note, []domain.Relationship) {
	if _, ok := notes[start]; !ok {
		return notes, rels
	}
	adj := map[string][]string{}
	for _, rel := range rels {
		adj[rel.FromID] = append(adj[rel.FromID], rel.ToID)
	}
	stack := []string{start}
	reachable := map[string]bool{}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if reachable[n] {
			continue
		}
		reachable[n] = true
		stack = append(stack, adj[n]...)
	}
	outNotes := map[string]domain.Note{}
	for id, note := range notes {
		if reachable[id] {
			outNotes[id] = note
		}
	}
	outRels := make([]domain.Relationship, 0)
	for _, rel := range rels {
		if reachable[rel.FromID] && reachable[rel.ToID] {
			outRels = append(outRels, rel)
		}
	}
	return outNotes, outRels
}
