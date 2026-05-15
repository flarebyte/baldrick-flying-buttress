// purpose: Implements graph operations in graph_args_labels.go so report generation and analysis can operate on deterministic graph views.
// responsibilities: select graph subsets; evaluate shape/properties; render or traverse graph structures; support query helpers used by CLI/report flows
// architecture_notes: Graph code keeps selection, shape detection, and rendering separated to prevent coupling traversal rules with presentation concerns.
package graph

import (
	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/textutil"
)

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
	return textutil.ParseKeyValue(entry)
}

func splitCSV(input string) []string {
	return textutil.SplitCSV(input)
}
