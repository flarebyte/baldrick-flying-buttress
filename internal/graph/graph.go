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

func RenderMarkdownText(selected Selected, shape Shape, policy CyclePolicy) string {
	adj := buildAdj(selected.Relationships)
	noteByID := map[string]domain.Note{}
	for _, note := range selected.Notes {
		noteByID[note.ID] = note
	}
	roots := selected.Roots
	if len(roots) == 0 {
		roots = noteIDs(selected.Notes)
	}

	var b strings.Builder
	switch shape {
	case ShapeTree:
		for _, root := range roots {
			renderTree(&b, root, 0, adj, noteByID)
		}
	case ShapeDAG:
		visits := map[string]int{}
		for _, root := range roots {
			renderDAG(&b, root, 0, adj, noteByID, visits, 2)
		}
	case ShapeCyclic:
		if policy == CyclePolicyAllow {
			rendered := map[string]bool{}
			for _, root := range roots {
				renderCyclic(&b, root, 0, adj, noteByID, rendered, map[string]bool{})
			}
		}
	}
	return b.String()
}

func renderTree(b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note) {
	note, ok := noteByID[nodeID]
	if !ok {
		return
	}
	writeNodeLine(b, note, depth)
	for _, child := range adj[nodeID] {
		renderTree(b, child, depth+1, adj, noteByID)
	}
}

func renderDAG(b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note, visits map[string]int, maxVisits int) {
	note, ok := noteByID[nodeID]
	if !ok {
		return
	}
	if visits[nodeID] >= maxVisits {
		indent := strings.Repeat("  ", depth)
		b.WriteString(indent)
		b.WriteString("- *(repeat limit reached: ")
		b.WriteString(note.Title)
		b.WriteString(")*\n")
		return
	}
	visits[nodeID]++
	writeNodeLine(b, note, depth)
	for _, child := range adj[nodeID] {
		renderDAG(b, child, depth+1, adj, noteByID, visits, maxVisits)
	}
}

func renderCyclic(b *strings.Builder, nodeID string, depth int, adj map[string][]string, noteByID map[string]domain.Note, rendered map[string]bool, stack map[string]bool) {
	note, ok := noteByID[nodeID]
	if !ok {
		return
	}
	if stack[nodeID] {
		indent := strings.Repeat("  ", depth)
		b.WriteString(indent)
		b.WriteString("- *(cycle back to ")
		b.WriteString(note.Title)
		b.WriteString(")*\n")
		return
	}
	if rendered[nodeID] {
		return
	}
	rendered[nodeID] = true
	writeNodeLine(b, note, depth)
	stack[nodeID] = true
	for _, child := range adj[nodeID] {
		renderCyclic(b, child, depth+1, adj, noteByID, rendered, stack)
	}
	delete(stack, nodeID)
}

func writeNodeLine(b *strings.Builder, note domain.Note, depth int) {
	indent := strings.Repeat("  ", depth)
	b.WriteString(indent)
	b.WriteString("- ")
	if strings.TrimSpace(note.Title) != "" {
		b.WriteString(note.Title)
	} else {
		b.WriteString(note.ID)
	}
	if strings.TrimSpace(note.Markdown) != "" {
		b.WriteString(": ")
		b.WriteString(note.Markdown)
	}
	b.WriteByte('\n')
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

func noteHasLabel(note domain.Note, label string) bool {
	for _, v := range splitCSV(note.LabelsCSV) {
		if v == label {
			return true
		}
	}
	return note.Label == label
}

func relationshipHasLabel(rel domain.Relationship, label string) bool {
	if rel.Label == label {
		return true
	}
	for _, v := range splitCSV(rel.LabelsCSV) {
		if v == label {
			return true
		}
	}
	return false
}

func hasSubject(subject map[string]domain.Note, id string) bool {
	_, ok := subject[id]
	return ok
}

func parseArg(entry string) (string, string, bool) {
	parts := strings.SplitN(entry, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	k := strings.TrimSpace(parts[0])
	v := strings.TrimSpace(parts[1])
	if k == "" || v == "" {
		return "", "", false
	}
	return k, v, true
}

func splitCSV(input string) []string {
	if strings.TrimSpace(input) == "" {
		return []string{}
	}
	items := strings.Split(input, ",")
	out := make([]string, 0, len(items))
	for _, item := range items {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}
