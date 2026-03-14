package cli

import (
	"bytes"
	"context"
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
	loader := pipeline.LoaderFunc(func(context.Context) (domain.RawApp, error) {
		return domain.RawApp{}, runtimeErr
	})
	validator := pipeline.ValidatorFunc(func(context.Context, domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		return domain.ValidatedApp{}, domain.ValidationReport{}, nil
	})

	exitCode, stdout, stderr := runCommand([]string{"validate"}, loader, validator)
	if exitCode != outcome.ExitCodeRuntimeFailure {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeRuntimeFailure, exitCode)
	}
	assertOutput(t, stdout, stderr, "", "runtime exploded\n")
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

	for _, args := range [][]string{{"validate"}, {"list", "reports"}, {"generate", "json"}, {"export", "cue"}} {
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
		{name: "export cue", args: []string{"export", "cue"}, loader: stubLoader(), validator: validatorWith(domain.ValidatedApp{}, domain.ValidationReport{}, nil), exitCode: 0},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assertDeterministic(t, tc.args, tc.loader, tc.validator, tc.exitCode)
		})
	}
}

func TestExecuteContextNormalPathMatchesExecuteWithFactory(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader {
		return load.FSAppLoader{ConfigPath: path}
	}
	validator := validate.AppDataValidator{}

	codeA, stdoutA, stderrA := runCommandWithFactory([]string{"list", "reports", "--config", configPath}, loaderFactory, validator)

	var outB bytes.Buffer
	var errB bytes.Buffer
	codeB := ExecuteContextWithFactory(context.Background(), []string{"list", "reports", "--config", configPath}, &outB, &errB, loaderFactory, validator)

	if codeA != codeB {
		t.Fatalf("exit code mismatch: %d vs %d", codeA, codeB)
	}
	if stdoutA != outB.String() || stderrA != errB.String() {
		t.Fatalf("output mismatch between ExecuteWithFactory and ExecuteContextWithFactory")
	}
}
