package validate

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func enrichDiagnostics(raw domain.RawApp, diagnostics []domain.Diagnostic) []domain.Diagnostic {
	if len(diagnostics) == 0 {
		return diagnostics
	}

	reportIDByTitle := make(map[string]string, len(raw.Reports))
	for _, report := range raw.Reports {
		reportID := domain.ReportIDFromFilepath(report.Filepath)
		if reportID == "" {
			reportID = strings.TrimSpace(report.Title)
		}
		if strings.TrimSpace(report.Title) != "" {
			reportIDByTitle[report.Title] = reportID
		}
	}

	noteTitleByName := make(map[string]string, len(raw.Notes))
	for _, note := range raw.Notes {
		if strings.TrimSpace(note.Name) != "" {
			noteTitleByName[note.Name] = note.Title
		}
	}

	configPathAbs := ""
	if strings.TrimSpace(raw.ConfigPath) != "" {
		if abs, err := filepath.Abs(raw.ConfigPath); err == nil {
			configPathAbs = abs
		}
	}

	out := make([]domain.Diagnostic, 0, len(diagnostics))
	for _, d := range diagnostics {
		if strings.TrimSpace(d.NormalizedPath) == "" {
			d.NormalizedPath = d.Path
		}
		if strings.TrimSpace(raw.ConfigPath) != "" {
			d.ConfigPath = raw.ConfigPath
		}
		if strings.TrimSpace(configPathAbs) != "" {
			d.ConfigPathAbs = configPathAbs
		}
		if strings.TrimSpace(d.ReportID) == "" && strings.TrimSpace(d.ReportTitle) != "" {
			d.ReportID = reportIDByTitle[d.ReportTitle]
		}
		if strings.TrimSpace(d.NoteTitle) == "" && strings.TrimSpace(d.NoteName) != "" {
			d.NoteTitle = noteTitleByName[d.NoteName]
		}
		d.RelatedNodes = relatedNodes(d)
		d.SuggestedFixes = suggestedFixes(d)
		out = append(out, d)
	}
	return out
}

func relatedNodes(d domain.Diagnostic) []string {
	nodes := make([]string, 0, 3)
	if strings.TrimSpace(d.NoteName) != "" {
		nodes = append(nodes, d.NoteName)
	}
	if strings.TrimSpace(d.RelationshipFrom) != "" {
		nodes = append(nodes, d.RelationshipFrom)
	}
	if strings.TrimSpace(d.RelationshipTo) != "" {
		nodes = append(nodes, d.RelationshipTo)
	}
	if len(nodes) == 0 {
		return nil
	}
	slices.Sort(nodes)
	return slices.Compact(nodes)
}

func suggestedFixes(d domain.Diagnostic) []string {
	switch d.Code {
	case "FBV000":
		return []string{"Set source to a non-empty string"}
	case "FBV100":
		return []string{"Add at least one report to reports"}
	case "FBV101":
		return []string{"Set report title to a non-empty string"}
	case "FBV102":
		return []string{"Set report filepath to a non-empty string"}
	case "FBV103":
		return []string{"Add at least one section to the report"}
	case "FBV104":
		return []string{"Set section title to a non-empty string"}
	case "FBV200":
		return []string{"Add at least one note to notes"}
	case "FBV201":
		return []string{"Set note name to a unique non-empty string"}
	case "FBV202":
		return []string{"Set note title to a non-empty string"}
	case "FBV300":
		return []string{"Add at least one relationship to relationships"}
	case "FBV301":
		return []string{"Set relationship from to an existing note name"}
	case "FBV302":
		return []string{"Set relationship to to an existing note name"}
	case "FBC001":
		return []string{"Use key=value format for configured arguments"}
	case "FBC002":
		return []string{"Remove the unknown argument or declare it in argumentRegistry"}
	case "FBC003":
		return []string{"Move the argument to an allowed scope or update argumentRegistry.scopes"}
	case "FBC004":
		return []string{"Use a value that matches the declared argument type"}
	case "FBR000":
		return []string{"Rename the duplicate argument so argument names are unique"}
	case "FBR001":
		return []string{"Set argument name to a non-empty string"}
	case "FBR002":
		return []string{"Set valueType for the argument"}
	case "FBR003":
		return []string{"Use one of the supported value types"}
	case "FBR004":
		return []string{"Declare at least one valid scope for the argument"}
	case "FBR005":
		return []string{"Use only valid scopes: h3-section, note, renderer"}
	case "FBR006":
		return []string{"Make enum defaultValue match one of allowedValues"}
	case "FBR007":
		return []string{"Remove duplicate enum allowedValues"}
	case "LABEL_REF_UNKNOWN":
		return []string{"Use an existing label or add the label to the referenced notes or relationships"}
	case "GRAPH_MISSING_NODE":
		return []string{"Add the missing note or fix the relationship endpoint"}
	case "GRAPH_ORPHAN_NODE":
		return []string{"Add a relationship or reference the note from a report section"}
	case "GRAPH_DUPLICATE_NOTE_NAME":
		return []string{"Rename duplicate notes so each note name is unique"}
	case "GRAPH_CROSS_REPORT_REFERENCE":
		return []string{"Move the related notes into a shared report or allow crossReportReference"}
	default:
		return nil
	}
}
