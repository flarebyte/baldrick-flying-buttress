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
