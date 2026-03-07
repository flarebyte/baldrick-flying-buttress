package cli

import (
	"bytes"
	"errors"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

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

func TestValidateWithInvalidRegistryConfigCollectsMultipleDiagnostics(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.registry.invalid.raw.json")
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
	assertOutput(t, validateStdout, validateStderr, readGolden(t, "validate_registry_invalid_output.golden"), "")
}

func TestValidateWithInvalidConfiguredArgumentsCollectsMultipleDiagnostics(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.args.invalid.raw.json")
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
	assertOutput(t, validateStdout, validateStderr, readGolden(t, "validate_args_invalid_output.golden"), "")
}

func TestValidateWithUnknownLabelReferencesEmitsWarningsOnly(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.labels.invalid.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader {
		return load.FSAppLoader{ConfigPath: path}
	}
	validator := validate.AppDataValidator{}

	validateCode, validateStdout, validateStderr := runCommandWithFactory(
		[]string{"validate", "--config", configPath},
		loaderFactory,
		validator,
	)
	if validateCode != 0 {
		t.Fatalf("expected validate exit code 0, got %d", validateCode)
	}
	assertOutput(t, validateStdout, validateStderr, readGolden(t, "validate_labels_invalid_output.golden"), "")

	listCode, listStdout, listStderr := runCommandWithFactory(
		[]string{"list", "reports", "--config", configPath},
		loaderFactory,
		validator,
	)
	if listCode != 0 {
		t.Fatalf("expected list exit code 0, got %d", listCode)
	}
	assertOutput(t, listStdout, listStderr, readGolden(t, "list_reports_output.golden"), "")
}

func TestValidateWithGraphIntegrityIssuesCollectsDiagnostics(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.graph.invalid.raw.json")
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
	assertOutput(t, validateStdout, validateStderr, readGolden(t, "validate_graph_invalid_output.golden"), "")
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
		{name: "list names", args: []string{"list", "names", "--prefix", "cli."}, loader: stubLoader(), validator: validatorWith(listNamesValidatedApp(), domain.ValidationReport{}, nil), exitCode: 0},
		{name: "lint names", args: []string{"lint", "names", "--style", "dot"}, loader: stubLoader(), validator: validatorWith(listNamesValidatedApp(), domain.ValidationReport{}, nil), exitCode: 0},
		{name: "lint orphans", args: []string{"lint", "orphans", "--subject-label", "ingredient"}, loader: stubLoader(), validator: validatorWith(orphansValidatedApp(), domain.ValidationReport{}, nil), exitCode: 0},
		{name: "generate json", args: []string{"generate", "json"}, loader: stubLoader(), validator: validatorWith(listValidatedApp(), domain.ValidationReport{}, nil), exitCode: 0},
		{name: "generate markdown", args: []string{"generate", "markdown"}, loader: stubLoader(), validator: validatorWith(domain.ValidatedApp{}, domain.ValidationReport{}, nil), exitCode: 0},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assertDeterministic(t, tc.args, tc.loader, tc.validator, tc.exitCode)
		})
	}
}
