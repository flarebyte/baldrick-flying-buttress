package cli

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/orphans"
)

func TestResolveLintOrphansQueryDefaults(t *testing.T) {
	t.Parallel()

	query, severity, err := resolveLintOrphansQuery("ingredient", "", "", "", "")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if query.Direction != orphans.DirectionEither {
		t.Fatalf("expected default direction either, got %s", query.Direction)
	}
	if severity != domain.SeverityWarning {
		t.Fatalf("expected warning severity, got %s", severity)
	}
}

func TestResolveLintOrphansQueryInvalidDirection(t *testing.T) {
	t.Parallel()

	_, _, err := resolveLintOrphansQuery("ingredient", "", "", "sideways", "warning")
	if err == nil {
		t.Fatal("expected invalid direction error")
	}
}

func TestResolveLintOrphansQuerySeverityError(t *testing.T) {
	t.Parallel()

	_, severity, err := resolveLintOrphansQuery("ingredient", "", "", "either", "error")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if severity != domain.SeverityError {
		t.Fatalf("expected error severity, got %s", severity)
	}
}

func TestLintOrphansDiagnosticsIncludeContext(t *testing.T) {
	t.Parallel()

	diagnostics := lintOrphansAction{
		query:    orphans.Query{SubjectLabel: "ingredient", EdgeLabel: "uses", CounterpartLabel: "tool", Direction: orphans.DirectionEither},
		severity: domain.SeverityWarning,
	}.diagnostics(orphansValidatedApp())

	if len(diagnostics) != 1 {
		t.Fatalf("expected one orphan diagnostic, got %#v", diagnostics)
	}
	d := diagnostics[0]
	if d.Code != "ORPHAN_QUERY_MISSING_LINK" || d.NoteName == "" || d.SubjectLabel != "ingredient" || d.EdgeLabel != "uses" || d.CounterpartLabel != "tool" {
		t.Fatalf("unexpected diagnostic context: %#v", d)
	}
}
