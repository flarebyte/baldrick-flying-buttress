package cli

import (
	"context"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
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

func TestValidateWithReportFilterTargetsSingleReport(t *testing.T) {
	t.Parallel()

	loader := pipeline.LoaderFunc(func(context.Context) (domain.RawApp, error) {
		return domain.RawApp{
			Source: "app",
			Reports: []domain.RawReport{
				{
					Title:    "Alpha Report",
					Filepath: "out/alpha.md",
					Sections: []domain.RawReportSection{{Title: "Overview", Arguments: []string{"unknown=x"}}},
				},
				{
					Title:    "Beta Report",
					Filepath: "out/beta.md",
					Sections: []domain.RawReportSection{{Title: "Overview"}},
				},
			},
			Notes:         []domain.RawNote{{Name: "n1", Title: "Service API"}},
			Relationships: []domain.RawRelationship{},
			Registry:      domain.RawArgumentRegistry{},
		}, nil
	})

	exitCode, stdout, stderr := runCommand([]string{"validate", "--report", "alpha"}, loader, validate.AppDataValidator{})
	if exitCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, exitCode)
	}
	want := "{\"diagnostics\":[{\"code\":\"FBC002\",\"severity\":\"error\",\"message\":\"unknown configured argument key [reportTitle=Alpha Report, sectionTitle=Overview, argumentName=unknown]\",\"path\":\"reports[0].sections[0].arguments[0]\"}]}\n"
	assertOutput(t, stdout, stderr, want, "")
}

func TestValidateWithUnknownReportFilterFailsAtRuntime(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand([]string{"validate", "--report", "missing"}, validate.StubAppLoader{}, validate.StubAppValidator{})
	if exitCode != outcome.ExitCodeRuntimeFailure {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeRuntimeFailure, exitCode)
	}
	assertOutput(t, stdout, stderr, "", "unknown report filter: missing (available: cpu-overview,memory-health)\n")
}
