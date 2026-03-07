package cli

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
)

func TestGenerateJSONGoldenOutput(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand([]string{"generate", "json"}, stubLoader(), validatorWith(listValidatedApp(), warningOnlyReport(), nil))
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "generate_json_output.golden"), "")
}

func TestGenerateJSONBlockedOnErrorDiagnostic(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand([]string{"generate", "json"}, stubLoader(), validatorWith(listValidatedApp(), errorOnlyReport(), nil))
	if exitCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, exitCode)
	}
	assertOutput(t, stdout, stderr, "", "")
}
