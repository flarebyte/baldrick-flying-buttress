package validate

import (
	"fmt"
	"slices"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func resolveGraphIntegrityPolicy(raw domain.RawGraphIntegrityPolicy) domain.GraphIntegrityPolicy {
	return domain.GraphIntegrityPolicy{
		MissingNode:          normalizePolicySeverity(raw.MissingNode, domain.PolicySeverityIgnore),
		OrphanNode:           normalizePolicySeverity(raw.OrphanNode, domain.PolicySeverityIgnore),
		DuplicateNoteName:    normalizePolicySeverity(raw.DuplicateNoteName, domain.PolicySeverityIgnore),
		CrossReportReference: normalizeCrossReportPolicy(raw.CrossReportReference),
	}
}

func normalizePolicySeverity(value string, fallback domain.PolicySeverity) domain.PolicySeverity {
	switch domain.PolicySeverity(strings.TrimSpace(value)) {
	case domain.PolicySeverityError, domain.PolicySeverityWarning, domain.PolicySeverityIgnore:
		return domain.PolicySeverity(strings.TrimSpace(value))
	default:
		return fallback
	}
}

func normalizeCrossReportPolicy(value string) domain.CrossReportPolicy {
	switch domain.CrossReportPolicy(strings.TrimSpace(value)) {
	case domain.CrossReportPolicyAllow, domain.CrossReportPolicyDisallow:
		return domain.CrossReportPolicy(strings.TrimSpace(value))
	default:
		return domain.CrossReportPolicyAllow
	}
}

func validateGraphIntegrity(raw domain.RawApp, policy domain.GraphIntegrityPolicy) []domain.Diagnostic {
	diagnostics := make([]domain.Diagnostic, 0)
	diagnostics = append(diagnostics, checkMissingNodes(raw, policy)...)
	diagnostics = append(diagnostics, checkOrphans(raw, policy.OrphanNode)...)
	diagnostics = append(diagnostics, checkDuplicateNoteNames(raw, policy.DuplicateNoteName)...)
	diagnostics = append(diagnostics, checkCrossReportReferences(raw, policy)...)
	return diagnostics
}

func checkMissingNodes(raw domain.RawApp, policy domain.GraphIntegrityPolicy) []domain.Diagnostic {
	if policy.MissingNode == domain.PolicySeverityIgnore {
		return []domain.Diagnostic{}
	}
	known := map[string]struct{}{}
	for _, note := range raw.Notes {
		name := strings.TrimSpace(note.Name)
		if name != "" {
			known[name] = struct{}{}
		}
	}
	diagnostics := make([]domain.Diagnostic, 0)
	for i, rel := range raw.Relationships {
		if rel.FromID != "" {
			if _, ok := known[rel.FromID]; !ok {
				diagnostics = append(diagnostics, newGraphDiagnostic(policy.MissingNode, graphMissingNodesSource, "GRAPH_MISSING_NODE", "relationship references missing from node", relationshipLocation(i, "from"), "", "", "", rel.FromID, rel.ToID))
			}
		}
		if rel.ToID != "" {
			if _, ok := known[rel.ToID]; !ok {
				diagnostics = append(diagnostics, newGraphDiagnostic(policy.MissingNode, graphMissingNodesSource, "GRAPH_MISSING_NODE", "relationship references missing to node", relationshipLocation(i, "to"), "", "", "", rel.FromID, rel.ToID))
			}
		}
	}
	return diagnostics
}

func checkOrphans(raw domain.RawApp, policy domain.PolicySeverity) []domain.Diagnostic {
	if policy == domain.PolicySeverityIgnore {
		return []domain.Diagnostic{}
	}
	connected := map[string]struct{}{}
	for _, rel := range raw.Relationships {
		if rel.FromID != "" {
			connected[rel.FromID] = struct{}{}
		}
		if rel.ToID != "" {
			connected[rel.ToID] = struct{}{}
		}
	}
	for _, report := range raw.Reports {
		for _, section := range report.Sections {
			for _, noteName := range section.Notes {
				n := strings.TrimSpace(noteName)
				if n != "" {
					connected[n] = struct{}{}
				}
			}
		}
	}
	diagnostics := make([]domain.Diagnostic, 0)
	for i, note := range raw.Notes {
		name := strings.TrimSpace(note.Name)
		if name == "" {
			continue
		}
		if _, ok := connected[name]; ok {
			continue
		}
		diagnostics = append(diagnostics, newGraphDiagnostic(policy, graphOrphansSource, "GRAPH_ORPHAN_NODE", "note is orphaned from relationships and report sections", noteLocation(i, "name"), "", "", name, "", ""))
	}
	return diagnostics
}

func checkDuplicateNoteNames(raw domain.RawApp, policy domain.PolicySeverity) []domain.Diagnostic {
	if policy == domain.PolicySeverityIgnore {
		return []domain.Diagnostic{}
	}
	counts := map[string]int{}
	for _, note := range raw.Notes {
		name := strings.TrimSpace(note.Name)
		if name != "" {
			counts[name]++
		}
	}
	names := make([]string, 0)
	for name, count := range counts {
		if count > 1 {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	diagnostics := make([]domain.Diagnostic, 0, len(names))
	for _, name := range names {
		location := fmt.Sprintf("notes[name=%q]", name)
		diagnostics = append(diagnostics, newGraphDiagnostic(policy, graphDuplicateNoteNamesSource, "GRAPH_DUPLICATE_NOTE_NAME", "duplicate note name", location, "", "", name, "", ""))
	}
	return diagnostics
}

func checkCrossReportReferences(raw domain.RawApp, policy domain.GraphIntegrityPolicy) []domain.Diagnostic {
	if policy.CrossReportReference == domain.CrossReportPolicyAllow {
		return []domain.Diagnostic{}
	}
	membership := collectNoteReportMembership(raw)
	diagnostics := make([]domain.Diagnostic, 0)
	for i, rel := range raw.Relationships {
		if rel.FromID == "" || rel.ToID == "" {
			continue
		}
		fromReports := membership[rel.FromID]
		toReports := membership[rel.ToID]
		if len(fromReports) == 0 || len(toReports) == 0 {
			continue
		}
		if hasReportOverlap(fromReports, toReports) {
			continue
		}
		diagnostics = append(diagnostics, newGraphDiagnostic(domain.PolicySeverityError, graphCrossReportSource, "GRAPH_CROSS_REPORT_REFERENCE", "cross-report relationship reference is disallowed", fmt.Sprintf("relationships[%d]", i), "", "", "", rel.FromID, rel.ToID))
	}
	return diagnostics
}

func collectNoteReportMembership(raw domain.RawApp) map[string][]string {
	membership := map[string]map[string]struct{}{}
	for _, report := range raw.Reports {
		reportID := domain.ReportIDFromFilepath(report.Filepath)
		if reportID == "" {
			reportID = report.Title
		}
		for _, section := range report.Sections {
			for _, noteName := range section.Notes {
				name := strings.TrimSpace(noteName)
				if name == "" {
					continue
				}
				if _, ok := membership[name]; !ok {
					membership[name] = map[string]struct{}{}
				}
				membership[name][reportID] = struct{}{}
			}
		}
	}

	out := map[string][]string{}
	for noteName, reportSet := range membership {
		reports := make([]string, 0, len(reportSet))
		for reportID := range reportSet {
			reports = append(reports, reportID)
		}
		slices.Sort(reports)
		out[noteName] = reports
	}
	return out
}

func hasReportOverlap(a, b []string) bool {
	for _, left := range a {
		for _, right := range b {
			if left == right {
				return true
			}
		}
	}
	return false
}
