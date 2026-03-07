package cli

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestResolveLintNamesPolicyDot(t *testing.T) {
	t.Parallel()
	policy, err := resolveLintNamesPolicy("dot", "", "warning")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if !policy.matcher("cli.root") || policy.matcher("cli_root") {
		t.Fatal("dot matcher mismatch")
	}
}

func TestResolveLintNamesPolicySnake(t *testing.T) {
	t.Parallel()
	policy, err := resolveLintNamesPolicy("snake", "", "warning")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if !policy.matcher("cli_root") || policy.matcher("cli.root") {
		t.Fatal("snake matcher mismatch")
	}
}

func TestResolveLintNamesPolicyRegex(t *testing.T) {
	t.Parallel()
	policy, err := resolveLintNamesPolicy("regex", "^cli\\.[a-z]+$", "warning")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if !policy.matcher("cli.root") || policy.matcher("app.db") {
		t.Fatal("regex matcher mismatch")
	}
}

func TestLintNamesPrefixScoping(t *testing.T) {
	t.Parallel()
	policy, _ := resolveLintNamesPolicy("dot", "", "warning")
	app := domain.ValidatedApp{
		Notes:         []domain.Note{{ID: "cli_root", Label: "x"}, {ID: "app_db", Label: "x"}},
		Relationships: []domain.Relationship{{FromID: "app_db", ToID: "cli_root", Label: "x"}},
	}
	diags := lintNames(app, "cli.", policy)
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics with prefix scope, got %#v", diags)
	}
}

func TestLintNamesSeverityOverride(t *testing.T) {
	t.Parallel()
	policy, _ := resolveLintNamesPolicy("dot", "", "error")
	app := domain.ValidatedApp{Notes: []domain.Note{{ID: "cli_root", Label: "x"}}}
	diags := lintNames(app, "", policy)
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %#v", diags)
	}
	if diags[0].Severity != domain.SeverityError {
		t.Fatalf("expected error severity, got %s", diags[0].Severity)
	}
}
