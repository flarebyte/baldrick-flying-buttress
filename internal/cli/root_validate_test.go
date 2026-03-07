package cli

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func TestValidateSuccessWarningsOnly(t *testing.T) {
	t.Parallel()

	app := listValidatedApp()
	report := domain.ValidationReport{Diagnostics: []domain.Diagnostic{{
		Code:     "FBW01",
		Severity: domain.SeverityWarning,
		Message:  "warning only",
		Path:     "module.stub",
	}}}

	exitCode, stdout, stderr := runCommand([]string{"validate"}, stubLoader(), validatorWith(app, report, nil))
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	want := "{\"diagnostics\":[{\"code\":\"FBW01\",\"severity\":\"warning\",\"message\":\"warning only\",\"path\":\"module.stub\"}]}\n"
	assertOutput(t, stdout, stderr, want, "")
}

func TestValidateFailureErrorDiagnostic(t *testing.T) {
	t.Parallel()

	app := listValidatedApp()
	report := domain.ValidationReport{Diagnostics: []domain.Diagnostic{{
		Code:     "FBE01",
		Severity: domain.SeverityError,
		Message:  "error diagnostic",
		Path:     "module.stub",
	}}}

	exitCode, stdout, stderr := runCommand([]string{"validate"}, stubLoader(), validatorWith(app, report, nil))
	if exitCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, exitCode)
	}
	want := "{\"diagnostics\":[{\"code\":\"FBE01\",\"severity\":\"error\",\"message\":\"error diagnostic\",\"path\":\"module.stub\"}]}\n"
	assertOutput(t, stdout, stderr, want, "")
}

func TestValidateGoldenOutput(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand([]string{"validate"}, validate.StubAppLoader{}, validate.StubAppValidator{})
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "validate_output.golden"), "")
}
