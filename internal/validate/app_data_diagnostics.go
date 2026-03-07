package validate

import (
	"fmt"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func newDiagnostic(source, code, message, location string) domain.Diagnostic {
	return domain.Diagnostic{
		Code:     code,
		Severity: domain.SeverityError,
		Source:   source,
		Message:  message,
		Location: location,
		Path:     location,
	}
}

func newArgumentDiagnostic(code, message, location, reportTitle, sectionTitle, noteName, argumentName string) domain.Diagnostic {
	d := newDiagnostic(configArgsValidationSource, code, message, location)
	d.ReportTitle = reportTitle
	d.SectionTitle = sectionTitle
	d.NoteName = noteName
	d.ArgumentName = argumentName
	d.Message = formatDiagnosticMessage(d.Message, d)
	return d
}

func newLabelReferenceDiagnostic(code, message, location, reportTitle, sectionTitle, noteName, argumentName, labelValue string) domain.Diagnostic {
	d := newDiagnostic(labelReferenceValidationSource, code, message, location)
	d.Severity = domain.SeverityWarning
	d.ReportTitle = reportTitle
	d.SectionTitle = sectionTitle
	d.NoteName = noteName
	d.ArgumentName = argumentName
	d.LabelValue = labelValue
	d.Message = formatDiagnosticMessage(d.Message, d)
	return d
}

func newGraphDiagnostic(policy domain.PolicySeverity, source, code, message, location, reportTitle, sectionTitle, noteName, relationshipFrom, relationshipTo string) domain.Diagnostic {
	if policy == domain.PolicySeverityIgnore {
		return domain.Diagnostic{}
	}
	d := newDiagnostic(source, code, message, location)
	if policy == domain.PolicySeverityWarning {
		d.Severity = domain.SeverityWarning
	} else {
		d.Severity = domain.SeverityError
	}
	d.ReportTitle = reportTitle
	d.SectionTitle = sectionTitle
	d.NoteName = noteName
	d.RelationshipFrom = relationshipFrom
	d.RelationshipTo = relationshipTo
	d.Message = formatDiagnosticMessage(d.Message, d)
	return d
}

func newRegistryDiagnostic(code, message, location, argumentName string) domain.Diagnostic {
	d := newDiagnostic(registryValidationSource, code, message, location)
	d.ArgumentName = argumentName
	d.Message = formatDiagnosticMessage(d.Message, d)
	return d
}

func noteLocation(i int, field string) string {
	return fmt.Sprintf("notes[%d].%s", i, field)
}

func relationshipLocation(i int, field string) string {
	return fmt.Sprintf("relationships[%d].%s", i, field)
}

func registryArgLocation(name string, index int) string {
	if name != "" {
		return fmt.Sprintf("argumentRegistry.arguments[name=%q]", name)
	}
	return fmt.Sprintf("argumentRegistry.arguments[%d]", index)
}

func registryArgNameLocation(name string) string {
	return fmt.Sprintf("argumentRegistry.arguments[name=%q]", name)
}

func formatDiagnosticMessage(base string, d domain.Diagnostic) string {
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
