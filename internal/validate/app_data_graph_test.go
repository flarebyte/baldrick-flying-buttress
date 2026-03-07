package validate

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestGraphIntegrityPolicyResolutionDefaultsAndNormalization(t *testing.T) {
	t.Parallel()

	policy := resolveGraphIntegrityPolicy(domain.RawGraphIntegrityPolicy{
		MissingNode:          "warning",
		OrphanNode:           "invalid",
		DuplicateNoteName:    "ignore",
		CrossReportReference: "disallow",
	})

	if policy.MissingNode != domain.PolicySeverityWarning {
		t.Fatalf("expected missing node warning, got %s", policy.MissingNode)
	}
	if policy.OrphanNode != domain.PolicySeverityIgnore {
		t.Fatalf("expected orphan fallback ignore, got %s", policy.OrphanNode)
	}
	if policy.DuplicateNoteName != domain.PolicySeverityIgnore {
		t.Fatalf("expected duplicate ignore, got %s", policy.DuplicateNoteName)
	}
	if policy.CrossReportReference != domain.CrossReportPolicyDisallow {
		t.Fatalf("expected cross-report disallow, got %s", policy.CrossReportReference)
	}
}

func TestValidatedAppIncludesResolvedGraphIntegrityPolicy(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source:        "app",
		Reports:       []domain.RawReport{{Title: "R", Filepath: "reports/r.md", Sections: []domain.RawReportSection{{Title: "S"}}}},
		Notes:         []domain.RawNote{{Name: "n1", Title: "N1"}},
		Relationships: []domain.RawRelationship{{FromID: "n1", ToID: "n1", Label: "L"}},
		GraphIntegrityPolicy: domain.RawGraphIntegrityPolicy{
			MissingNode: "warning",
		},
	}

	app, _, err := AppDataValidator{}.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if app.GraphIntegrityPolicy.MissingNode != domain.PolicySeverityWarning {
		t.Fatalf("expected resolved policy in validated app, got %#v", app.GraphIntegrityPolicy)
	}
}

func TestGraphIntegrityValidationChecks(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source: "app",
		Reports: []domain.RawReport{
			{
				Title:    "R1",
				Filepath: "reports/r1.md",
				Sections: []domain.RawReportSection{{Title: "S1", Notes: []string{"n1"}}},
			},
			{
				Title:    "R2",
				Filepath: "reports/r2.md",
				Sections: []domain.RawReportSection{{Title: "S2", Notes: []string{"n2"}}},
			},
		},
		Notes: []domain.RawNote{
			{Name: "n1", Title: "N1"},
			{Name: "n2", Title: "N2"},
			{Name: "n3", Title: "N3"},
			{Name: "dup", Title: "D1"},
			{Name: "dup", Title: "D2"},
		},
		Relationships: []domain.RawRelationship{
			{FromID: "n1", ToID: "missing", Label: "depends_on"},
			{FromID: "n1", ToID: "n2", Label: "depends_on"},
		},
		GraphIntegrityPolicy: domain.RawGraphIntegrityPolicy{
			MissingNode:          "error",
			OrphanNode:           "warning",
			DuplicateNoteName:    "warning",
			CrossReportReference: "disallow",
		},
	}

	_, report, err := AppDataValidator{}.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	expected := []struct {
		code     string
		source   string
		location string
		severity domain.Severity
	}{
		{"GRAPH_MISSING_NODE", graphMissingNodesSource, "relationships[0].to", domain.SeverityError},
		{"GRAPH_ORPHAN_NODE", graphOrphansSource, "notes[2].name", domain.SeverityWarning},
		{"GRAPH_DUPLICATE_NOTE_NAME", graphDuplicateNoteNamesSource, "notes[name=\"dup\"]", domain.SeverityWarning},
		{"GRAPH_CROSS_REPORT_REFERENCE", graphCrossReportSource, "relationships[1]", domain.SeverityError},
	}
	for _, e := range expected {
		found := false
		for _, d := range report.Diagnostics {
			if d.Code == e.code && d.Source == e.source && d.Path == e.location && d.Severity == e.severity {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected graph diagnostic %s at %s (%s), got %#v", e.code, e.location, e.source, report.Diagnostics)
		}
	}
}

func TestGraphIntegritySeverityFollowsPolicy(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source:               "app",
		Reports:              []domain.RawReport{{Title: "R", Filepath: "reports/r.md", Sections: []domain.RawReportSection{{Title: "S"}}}},
		Notes:                []domain.RawNote{{Name: "n1", Title: "N1"}},
		Relationships:        []domain.RawRelationship{{FromID: "n1", ToID: "missing", Label: "depends_on"}},
		GraphIntegrityPolicy: domain.RawGraphIntegrityPolicy{MissingNode: "warning"},
	}

	_, report, err := AppDataValidator{}.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	for _, d := range report.Diagnostics {
		if d.Code == "GRAPH_MISSING_NODE" {
			if d.Severity != domain.SeverityWarning {
				t.Fatalf("expected warning severity for missing node, got %s", d.Severity)
			}
			return
		}
	}
	t.Fatalf("missing GRAPH_MISSING_NODE diagnostic: %#v", report.Diagnostics)
}
