// purpose: Declares core domain models in diagnostic_message.go used across validation, pipeline, and CLI layers.
// responsibilities: define canonical data structures; hold stable typed contracts; provide shared semantic primitives
// architecture_notes: Domain types are dependency-light and reused broadly to keep cross-layer contracts explicit and avoid cyclic package coupling.
package domain

import "strings"

func FormatDiagnosticMessage(base string, d Diagnostic) string {
	parts := make([]string, 0, 9)
	if strings.TrimSpace(d.ReportTitle) != "" {
		parts = append(parts, "reportTitle="+d.ReportTitle)
	}
	if strings.TrimSpace(d.SectionTitle) != "" {
		parts = append(parts, "sectionTitle="+d.SectionTitle)
	}
	if strings.TrimSpace(d.NoteName) != "" {
		parts = append(parts, "noteName="+d.NoteName)
	}
	if strings.TrimSpace(d.RelationshipFrom) != "" || strings.TrimSpace(d.RelationshipTo) != "" {
		parts = append(parts, "relationship="+d.RelationshipFrom+"->"+d.RelationshipTo)
	}
	if strings.TrimSpace(d.ArgumentName) != "" {
		parts = append(parts, "argumentName="+d.ArgumentName)
	}
	if strings.TrimSpace(d.SubjectLabel) != "" {
		parts = append(parts, "subjectLabel="+d.SubjectLabel)
	}
	if strings.TrimSpace(d.EdgeLabel) != "" {
		parts = append(parts, "edgeLabel="+d.EdgeLabel)
	}
	if strings.TrimSpace(d.CounterpartLabel) != "" {
		parts = append(parts, "counterpartLabel="+d.CounterpartLabel)
	}
	if strings.TrimSpace(d.LabelValue) != "" {
		parts = append(parts, "labelValue="+d.LabelValue)
	}
	if len(parts) == 0 {
		return base
	}
	return base + " [" + strings.Join(parts, ", ") + "]"
}
