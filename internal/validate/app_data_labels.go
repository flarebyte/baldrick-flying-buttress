package validate

import (
	"fmt"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func collectDatasetLabels(raw domain.RawApp) []string {
	labels := make([]string, 0)
	for _, note := range raw.Notes {
		labels = append(labels, note.Labels...)
	}
	for _, relationship := range raw.Relationships {
		labels = append(labels, relationship.Labels...)
		if relationship.Label != "" {
			labels = append(labels, relationship.Label)
		}
	}
	return normalizeAllowedValues(labels)
}

func validateLabelReferences(raw domain.RawApp, labelSet []string) []domain.Diagnostic {
	diagnostics := make([]domain.Diagnostic, 0)

	for i, report := range raw.Reports {
		for j, section := range report.Sections {
			for k, entry := range section.Arguments {
				name, value, ok := parseConfiguredArgument(entry)
				if !ok || !isLabelReferenceArgument(name) {
					continue
				}
				if containsString(labelSet, value) {
					continue
				}
				location := fmt.Sprintf("reports[%d].sections[%d].arguments[%d]", i, j, k)
				diagnostics = append(diagnostics, newLabelReferenceDiagnostic("LABEL_REF_UNKNOWN", "unknown label reference", location, report.Title, section.Title, "", name, value))
			}
		}
	}

	for i, note := range raw.Notes {
		for j, entry := range note.Arguments {
			name, value, ok := parseConfiguredArgument(entry)
			if !ok || !isLabelReferenceArgument(name) {
				continue
			}
			if containsString(labelSet, value) {
				continue
			}
			location := fmt.Sprintf("notes[%d].arguments[%d]", i, j)
			diagnostics = append(diagnostics, newLabelReferenceDiagnostic("LABEL_REF_UNKNOWN", "unknown label reference", location, "", "", note.Name, name, value))
		}
	}

	return diagnostics
}
