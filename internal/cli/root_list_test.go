package cli

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
)

func TestListReportsGoldenOutput(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand([]string{"list", "reports"}, stubLoader(), validatorWith(listValidatedApp(), warningOnlyReport(), nil))
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "list_reports_output.golden"), "")
}

func TestListReportsBlockedOnErrorDiagnostic(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand([]string{"list", "reports"}, stubLoader(), validatorWith(listValidatedApp(), errorOnlyReport(), nil))
	if exitCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, exitCode)
	}
	assertOutput(t, stdout, stderr, "", "")
}

func TestListNamesDefaultTableOutput(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand(
		[]string{"list", "names", "--prefix", "cli."},
		stubLoader(),
		validatorWith(listNamesValidatedApp(), domain.ValidationReport{}, nil),
	)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "list_names_table_output.golden"), "")
}

func TestListNamesKindNotesOutput(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand(
		[]string{"list", "names", "--prefix", "cli.", "--kind", "notes"},
		stubLoader(),
		validatorWith(listNamesValidatedApp(), domain.ValidationReport{}, nil),
	)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	want := "KIND\tNAME\tFROM\tTO\nnote\tcli.root\t\t\nnote\tcli.worker\t\t\n"
	assertOutput(t, stdout, stderr, want, "")
}

func TestListNamesKindRelationshipsOutput(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand(
		[]string{"list", "names", "--prefix", "cli.", "--kind", "relationships"},
		stubLoader(),
		validatorWith(listNamesValidatedApp(), domain.ValidationReport{}, nil),
	)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	want := "KIND\tNAME\tFROM\tTO\nrelationship\t\tapp.db\tcli.worker\nrelationship\t\tcli.root\tapp.db\n"
	assertOutput(t, stdout, stderr, want, "")
}

func TestListNamesJSONOutput(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand(
		[]string{"list", "names", "--prefix", "cli.", "--format", "json"},
		stubLoader(),
		validatorWith(listNamesValidatedApp(), domain.ValidationReport{}, nil),
	)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "list_names_json_output.golden"), "")
}

func TestListNamesBlockedOnErrorDiagnostic(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand(
		[]string{"list", "names", "--prefix", "cli."},
		stubLoader(),
		validatorWith(listNamesValidatedApp(), errorOnlyReport(), nil),
	)
	if exitCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, exitCode)
	}
	assertOutput(t, stdout, stderr, "", "")
}

func TestListNamesMissingPrefixReturnsRuntimeFailure(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand(
		[]string{"list", "names"},
		stubLoader(),
		validatorWith(listNamesValidatedApp(), domain.ValidationReport{}, nil),
	)
	if exitCode != outcome.ExitCodeRuntimeFailure {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeRuntimeFailure, exitCode)
	}
	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if stderr == "" {
		t.Fatal("expected stderr for missing required prefix")
	}
}
