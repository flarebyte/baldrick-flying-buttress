package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
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

func TestRuntimeFailureMapsToDistinctExitCodeAndStderr(t *testing.T) {
	t.Parallel()

	runtimeErr := errors.New("runtime exploded")
	loader := pipeline.LoaderFunc(func() (domain.RawApp, error) {
		return domain.RawApp{}, runtimeErr
	})
	validator := pipeline.ValidatorFunc(func(domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		return domain.ValidatedApp{}, domain.ValidationReport{}, nil
	})

	exitCode, stdout, stderr := runCommand([]string{"validate"}, loader, validator)
	if exitCode != outcome.ExitCodeRuntimeFailure {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeRuntimeFailure, exitCode)
	}
	assertOutput(t, stdout, stderr, "", "runtime exploded\n")
}

func TestCommandsWorkWithConfigFlagAndFilesystemLoader(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader {
		return load.FSAppLoader{ConfigPath: path}
	}
	validator := validate.AppDataValidator{}
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
	}{
		{
			name:       "validate",
			args:       []string{"validate", "--config", configPath},
			wantCode:   0,
			wantStdout: "{\"diagnostics\":[]}\n",
		},
		{
			name:       "list reports",
			args:       []string{"list", "reports", "--config", configPath},
			wantCode:   0,
			wantStdout: readGolden(t, "list_reports_output.golden"),
		},
		{
			name:       "generate json",
			args:       []string{"generate", "json", "--config", configPath},
			wantCode:   0,
			wantStdout: readGolden(t, "generate_json_output.golden"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var out bytes.Buffer
			var errOut bytes.Buffer
			exitCode := ExecuteWithFactory(tc.args, &out, &errOut, loaderFactory, validator)
			if exitCode != tc.wantCode {
				t.Fatalf("expected exit code %d, got %d", tc.wantCode, exitCode)
			}
			assertOutput(t, out.String(), errOut.String(), tc.wantStdout, "")
		})
	}
}

func TestCommandsWithInvalidStructureConfigProduceValidationDiagnostics(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.invalid.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader {
		return load.FSAppLoader{ConfigPath: path}
	}
	validator := validate.AppDataValidator{}

	validateCode, validateStdout, validateStderr := runCommandWithFactory(
		[]string{"validate", "--config", configPath},
		loaderFactory,
		validator,
	)
	if validateCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected validate exit code %d, got %d", outcome.ExitCodeValidationBlocked, validateCode)
	}
	if validateStderr != "" {
		t.Fatalf("expected empty validate stderr, got %q", validateStderr)
	}
	assertOutput(t, validateStdout, validateStderr, readGolden(t, "validate_invalid_output.golden"), "")

	listCode, listStdout, listStderr := runCommandWithFactory(
		[]string{"list", "reports", "--config", configPath},
		loaderFactory,
		validator,
	)
	if listCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected list exit code %d, got %d", outcome.ExitCodeValidationBlocked, listCode)
	}
	assertOutput(t, listStdout, listStderr, "", "")
}

func TestFilesystemLoaderDeterministicOutputAcrossRuns(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader {
		return load.FSAppLoader{ConfigPath: path}
	}
	validator := validate.AppDataValidator{}

	var out1 bytes.Buffer
	var err1 bytes.Buffer
	code1 := ExecuteWithFactory([]string{"list", "reports", "--config", configPath}, &out1, &err1, loaderFactory, validator)
	if code1 != 0 {
		t.Fatalf("expected first exit code 0, got %d", code1)
	}
	if err1.String() != "" {
		t.Fatalf("expected empty first stderr, got %q", err1.String())
	}

	var out2 bytes.Buffer
	var err2 bytes.Buffer
	code2 := ExecuteWithFactory([]string{"list", "reports", "--config", configPath}, &out2, &err2, loaderFactory, validator)
	if code2 != 0 {
		t.Fatalf("expected second exit code 0, got %d", code2)
	}
	if err2.String() != "" {
		t.Fatalf("expected empty second stderr, got %q", err2.String())
	}

	if out1.String() != out2.String() {
		t.Fatalf("non-deterministic stdout\nfirst: %q\nsecond: %q", out1.String(), out2.String())
	}
}

func TestCanonicalOrderingFromDifferentInputOrdersProducesIdenticalOutput(t *testing.T) {
	t.Parallel()

	orderedReport := domain.ValidationReport{Diagnostics: []domain.Diagnostic{
		{Code: "FB001", Severity: domain.SeverityWarning, Message: "w", Path: "p1"},
		{Code: "FB002", Severity: domain.SeverityError, Message: "e", Path: "p2"},
	}}
	unorderedReport := domain.ValidationReport{Diagnostics: []domain.Diagnostic{
		{Code: "FB002", Severity: domain.SeverityError, Message: "e", Path: "p2"},
		{Code: "FB001", Severity: domain.SeverityWarning, Message: "w", Path: "p1"},
	}}

	orderedValidator := validatorWith(orderedValidatedAppForOrdering(), orderedReport, nil)
	unorderedValidator := validatorWith(unorderedValidatedApp(), unorderedReport, nil)

	for _, args := range [][]string{{"validate"}, {"list", "reports"}, {"generate", "json"}} {
		orderedCode, orderedStdout, orderedStderr := runCommand(args, stubLoader(), orderedValidator)
		unorderedCode, unorderedStdout, unorderedStderr := runCommand(args, stubLoader(), unorderedValidator)

		if orderedCode != unorderedCode {
			t.Fatalf("exit code mismatch for %v: %d vs %d", args, orderedCode, unorderedCode)
		}
		if orderedStdout != unorderedStdout {
			t.Fatalf("stdout mismatch for %v\nordered: %q\nunordered: %q", args, orderedStdout, unorderedStdout)
		}
		if orderedStderr != unorderedStderr {
			t.Fatalf("stderr mismatch for %v\nordered: %q\nunordered: %q", args, orderedStderr, unorderedStderr)
		}
	}
}

func TestDeterministicOutputAcrossRuns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		loader    pipeline.AppLoader
		validator pipeline.AppValidator
		exitCode  int
	}{
		{name: "validate", args: []string{"validate"}, loader: validate.StubAppLoader{}, validator: validate.StubAppValidator{}, exitCode: 0},
		{name: "list reports", args: []string{"list", "reports"}, loader: stubLoader(), validator: validatorWith(listValidatedApp(), domain.ValidationReport{}, nil), exitCode: 0},
		{name: "generate json", args: []string{"generate", "json"}, loader: stubLoader(), validator: validatorWith(listValidatedApp(), domain.ValidationReport{}, nil), exitCode: 0},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assertDeterministic(t, tc.args, tc.loader, tc.validator, tc.exitCode)
		})
	}
}

func assertDeterministic(t *testing.T, args []string, loader pipeline.AppLoader, validator pipeline.AppValidator, wantCode int) {
	t.Helper()

	code1, out1, err1 := runCommand(args, loader, validator)
	if code1 != wantCode {
		t.Fatalf("expected first exit code %d, got %d", wantCode, code1)
	}
	code2, out2, err2 := runCommand(args, loader, validator)
	if code2 != wantCode {
		t.Fatalf("expected second exit code %d, got %d", wantCode, code2)
	}
	if out1 != out2 {
		t.Fatalf("non-deterministic stdout\nfirst: %q\nsecond: %q", out1, out2)
	}
	if err1 != err2 {
		t.Fatalf("non-deterministic stderr\nfirst: %q\nsecond: %q", err1, err2)
	}
}

func runCommand(args []string, loader pipeline.AppLoader, validator pipeline.AppValidator) (int, string, string) {
	return runCommandWithFactory(args, func(string) pipeline.AppLoader {
		return loader
	}, validator)
}

func runCommandWithFactory(args []string, loaderFactory LoaderFactory, validator pipeline.AppValidator) (int, string, string) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	code := ExecuteWithFactory(args, &out, &errOut, loaderFactory, validator)
	return code, out.String(), errOut.String()
}

func stubLoader() pipeline.AppLoader {
	return pipeline.LoaderFunc(func() (domain.RawApp, error) {
		return domain.RawApp{Source: "in-memory-stub"}, nil
	})
}

func validatorWith(app domain.ValidatedApp, report domain.ValidationReport, err error) pipeline.AppValidator {
	return pipeline.ValidatorFunc(func(domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		return app, report, err
	})
}

func readGolden(t *testing.T, filename string) string {
	t.Helper()
	p := filepath.Join("testdata", filename)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read golden %s: %v", filename, err)
	}
	return string(b)
}

func assertOutput(t *testing.T, gotStdout, gotStderr, wantStdout, wantStderr string) {
	t.Helper()
	if gotStdout != wantStdout {
		t.Fatalf("stdout mismatch\nwant: %q\n got: %q", wantStdout, gotStdout)
	}
	if gotStderr != wantStderr {
		t.Fatalf("stderr mismatch\nwant: %q\n got: %q", wantStderr, gotStderr)
	}
}

func listValidatedApp() domain.ValidatedApp {
	return appFixture(
		[]string{"core", "edge"},
		[]domain.Report{{ID: "cpu-overview", Title: "CPU Overview"}, {ID: "memory-health", Title: "Memory Health"}},
		[]domain.Note{{ID: "n1", Label: "service.api"}, {ID: "n2", Label: "service.db"}},
		[]domain.Relationship{{FromID: "n1", ToID: "n2", Label: "depends_on"}},
	)
}

func unorderedValidatedApp() domain.ValidatedApp {
	return appFixture(
		[]string{"edge", "core"},
		[]domain.Report{{ID: "memory-health", Title: "Memory Health"}, {ID: "cpu-overview", Title: "CPU Overview"}},
		[]domain.Note{{ID: "n2", Label: "service.db"}, {ID: "n1", Label: "service.api"}},
		[]domain.Relationship{{FromID: "n1", ToID: "n2", Label: "depends_on"}, {FromID: "n1", ToID: "n2", Label: "owns"}},
	)
}

func orderedValidatedAppForOrdering() domain.ValidatedApp {
	return appFixture(
		[]string{"core", "edge"},
		[]domain.Report{{ID: "cpu-overview", Title: "CPU Overview"}, {ID: "memory-health", Title: "Memory Health"}},
		[]domain.Note{{ID: "n1", Label: "service.api"}, {ID: "n2", Label: "service.db"}},
		[]domain.Relationship{{FromID: "n1", ToID: "n2", Label: "depends_on"}, {FromID: "n1", ToID: "n2", Label: "owns"}},
	)
}

func appFixture(modules []string, reports []domain.Report, notes []domain.Note, relationships []domain.Relationship) domain.ValidatedApp {
	return domain.ValidatedApp{
		Name:          "stub-app",
		Modules:       modules,
		Reports:       reports,
		Notes:         notes,
		Relationships: relationships,
	}
}

func warningOnlyReport() domain.ValidationReport {
	return domain.ValidationReport{
		Diagnostics: []domain.Diagnostic{{
			Code:     "FBW01",
			Severity: domain.SeverityWarning,
			Message:  "warning only",
			Path:     "module.stub",
		}},
	}
}

func errorOnlyReport() domain.ValidationReport {
	return domain.ValidationReport{
		Diagnostics: []domain.Diagnostic{{
			Code:     "FBE01",
			Severity: domain.SeverityError,
			Message:  "error diagnostic",
			Path:     "module.stub",
		}},
	}
}
