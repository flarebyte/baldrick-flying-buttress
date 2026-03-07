package cli

import (
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func TestLintNamesDefaultDotStyle(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.lint.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"lint", "names", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "lint_names_dot_output.golden"), "")
}

func TestLintNamesSnakeStyle(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.lint.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"lint", "names", "--config", configPath, "--style", "snake"}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "lint_names_snake_output.golden"), "")
}

func TestLintNamesRegexStyle(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.lint.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"lint", "names", "--config", configPath, "--style", "regex", "--pattern", "^cli\\..+$"}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "lint_names_regex_output.golden"), "")
}

func TestLintNamesRegexMissingPatternRuntimeFailure(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.lint.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"lint", "names", "--config", configPath, "--style", "regex"}, loaderFactory, validator)
	if code != outcome.ExitCodeRuntimeFailure {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeRuntimeFailure, code)
	}
	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if stderr == "" {
		t.Fatal("expected stderr")
	}
}

func TestLintNamesSeverityErrorProducesBlockingDiagnostics(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.lint.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"lint", "names", "--config", configPath, "--severity", "error"}, loaderFactory, validator)
	if code != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "lint_names_dot_error_output.golden"), "")
}

func TestLintOrphansSubjectLabel(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.orphans.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"lint", "orphans", "--config", configPath, "--subject-label", "ingredient"}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "lint_orphans_default_output.golden"), "")
}

func TestLintOrphansWithEdgeLabelFilter(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.orphans.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"lint", "orphans", "--config", configPath, "--subject-label", "ingredient", "--edge-label", "uses"}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "lint_orphans_edge_output.golden"), "")
}

func TestLintOrphansWithCounterpartLabelFilter(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.orphans.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"lint", "orphans", "--config", configPath, "--subject-label", "ingredient", "--counterpart-label", "tool"}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "lint_orphans_counterpart_output.golden"), "")
}

func TestLintOrphansWithDirectionFilter(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.orphans.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"lint", "orphans", "--config", configPath, "--subject-label", "ingredient", "--direction", "out"}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "lint_orphans_out_output.golden"), "")
}

func TestLintOrphansSeverityErrorProducesBlockingDiagnostics(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.orphans.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"lint", "orphans", "--config", configPath, "--subject-label", "ingredient", "--severity", "error"}, loaderFactory, validator)
	if code != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "lint_orphans_error_output.golden"), "")
}

func TestLintOrphansInvalidDirectionRuntimeFailure(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.orphans.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"lint", "orphans", "--config", configPath, "--subject-label", "ingredient", "--direction", "sideways"}, loaderFactory, validator)
	if code != outcome.ExitCodeRuntimeFailure {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeRuntimeFailure, code)
	}
	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if stderr == "" {
		t.Fatal("expected stderr")
	}
}

func TestLintOrphansBlockedOnValidationErrors(t *testing.T) {
	t.Parallel()

	code, stdout, stderr := runCommand(
		[]string{"lint", "orphans", "--subject-label", "ingredient"},
		stubLoader(),
		validatorWith(orphansValidatedApp(), errorOnlyReport(), nil),
	)
	if code != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, code)
	}
	assertOutput(t, stdout, stderr, "", "")
}

func TestLintOrphansMissingSubjectLabelRuntimeFailure(t *testing.T) {
	t.Parallel()

	code, stdout, stderr := runCommand(
		[]string{"lint", "orphans"},
		stubLoader(),
		validatorWith(orphansValidatedApp(), domain.ValidationReport{}, nil),
	)
	if code != outcome.ExitCodeRuntimeFailure {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeRuntimeFailure, code)
	}
	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if stderr == "" {
		t.Fatal("expected stderr")
	}
}
