package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/safety"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

type commandExpectation struct {
	name       string
	args       []string
	wantCode   int
	wantStdout string
}

func fsLoaderFactory(path string) pipeline.AppLoader {
	return load.FSAppLoader{ConfigPath: path}
}

func runAndAssertCommand(t *testing.T, args []string, wantCode int, wantStdout string) {
	t.Helper()
	code, stdout, stderr := runCommandWithFactory(args, fsLoaderFactory, validate.AppDataValidator{})
	if code != wantCode {
		t.Fatalf("expected exit code %d, got %d", wantCode, code)
	}
	assertOutput(t, stdout, stderr, wantStdout, "")
}

func runAndAssertCommandSet(t *testing.T, tests []commandExpectation) {
	t.Helper()
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			runAndAssertCommand(t, tc.args, tc.wantCode, tc.wantStdout)
		})
	}
}

func assertValidateAndListBlocked(t *testing.T, configPath, validateGolden string) {
	t.Helper()
	validateCode, validateStdout, validateStderr := runCommandWithFactory(
		[]string{"validate", "--config", configPath},
		fsLoaderFactory,
		validate.AppDataValidator{},
	)
	if validateCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected validate exit code %d, got %d", outcome.ExitCodeValidationBlocked, validateCode)
	}
	if validateStderr != "" {
		t.Fatalf("expected empty validate stderr, got %q", validateStderr)
	}
	assertOutput(t, validateStdout, validateStderr, readGolden(t, validateGolden), "")

	listCode, listStdout, listStderr := runCommandWithFactory(
		[]string{"list", "reports", "--config", configPath},
		fsLoaderFactory,
		validate.AppDataValidator{},
	)
	if listCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected list exit code %d, got %d", outcome.ExitCodeValidationBlocked, listCode)
	}
	assertOutput(t, listStdout, listStderr, "", "")
}

func assertValidateOutput(t *testing.T, configPath string, wantCode int, golden string) {
	t.Helper()
	runAndAssertCommand(t, []string{"validate", "--config", configPath}, wantCode, readGolden(t, golden))
}

func assertCommandsWorkWithConfig(t *testing.T, configPath string) {
	t.Helper()
	runAndAssertCommandSet(t, []commandExpectation{
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
	})
}

func TestCommandsWorkWithConfigFlagAndFilesystemLoader(t *testing.T) {
	t.Parallel()

	assertCommandsWorkWithConfig(t, filepath.Join("testdata", "config.raw.json"))
}

func TestCommandsWorkWithCueConfigAndFilesystemLoader(t *testing.T) {
	t.Parallel()

	assertCommandsWorkWithConfig(t, filepath.Join("testdata", "config.cue"))
}

func TestCommandsWorkWithConfigDirectoryAndEquivalentAppCueFile(t *testing.T) {
	t.Parallel()

	dirPath := filepath.Join("..", "..", "doc", "design-meta")
	filePath := filepath.Join(dirPath, "app.cue")

	validateDirCode, validateDirStdout, validateDirStderr := runCommandWithFactory(
		[]string{"validate", "--config", dirPath},
		fsLoaderFactory,
		validate.AppDataValidator{},
	)
	validateFileCode, validateFileStdout, validateFileStderr := runCommandWithFactory(
		[]string{"validate", "--config", filePath},
		fsLoaderFactory,
		validate.AppDataValidator{},
	)
	if validateDirCode != 0 || validateFileCode != 0 {
		t.Fatalf("expected validate exit code 0, got %d and %d", validateDirCode, validateFileCode)
	}
	if validateDirStdout != validateFileStdout {
		t.Fatalf("expected matching validate stdout\nfrom dir: %q\nfrom file: %q", validateDirStdout, validateFileStdout)
	}
	if validateDirStderr != validateFileStderr {
		t.Fatalf("expected matching validate stderr\nfrom dir: %q\nfrom file: %q", validateDirStderr, validateFileStderr)
	}

	listDirCode, listDirStdout, listDirStderr := runCommandWithFactory(
		[]string{"list", "reports", "--config", dirPath},
		fsLoaderFactory,
		validate.AppDataValidator{},
	)
	listFileCode, listFileStdout, listFileStderr := runCommandWithFactory(
		[]string{"list", "reports", "--config", filePath},
		fsLoaderFactory,
		validate.AppDataValidator{},
	)
	if listDirCode != 0 || listFileCode != 0 {
		t.Fatalf("expected list exit code 0, got %d and %d", listDirCode, listFileCode)
	}
	if listDirStdout != listFileStdout {
		t.Fatalf("expected matching list stdout\nfrom dir: %q\nfrom file: %q", listDirStdout, listFileStdout)
	}
	if listDirStderr != listFileStderr {
		t.Fatalf("expected matching list stderr\nfrom dir: %q\nfrom file: %q", listDirStderr, listFileStderr)
	}
}

func TestCommandsWithInvalidStructureConfigProduceValidationDiagnostics(t *testing.T) {
	t.Parallel()

	assertValidateAndListBlocked(t, filepath.Join("testdata", "config.invalid.raw.json"), "validate_invalid_output.golden")
}

func TestCommandsWithInvalidStructureCueConfigProduceValidationDiagnostics(t *testing.T) {
	t.Parallel()

	assertValidateAndListBlocked(t, filepath.Join("testdata", "config.invalid.cue"), "validate_invalid_output.golden")
}

func TestValidateWithInvalidRegistryConfigCollectsMultipleDiagnostics(t *testing.T) {
	t.Parallel()

	assertValidateOutput(
		t,
		filepath.Join("testdata", "config.registry.invalid.raw.json"),
		outcome.ExitCodeValidationBlocked,
		"validate_registry_invalid_output.golden",
	)
}

func TestValidateWithInvalidConfiguredArgumentsCollectsMultipleDiagnostics(t *testing.T) {
	t.Parallel()

	assertValidateOutput(
		t,
		filepath.Join("testdata", "config.args.invalid.raw.json"),
		outcome.ExitCodeValidationBlocked,
		"validate_args_invalid_output.golden",
	)
}

func TestValidateWithUnknownLabelReferencesEmitsWarningsOnly(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.labels.invalid.raw.json")
	runAndAssertCommand(
		t,
		[]string{"validate", "--config", configPath},
		0,
		readGolden(t, "validate_labels_invalid_output.golden"),
	)
	runAndAssertCommand(
		t,
		[]string{"list", "reports", "--config", configPath},
		0,
		readGolden(t, "list_reports_output.golden"),
	)
}

func TestValidateWithGraphIntegrityIssuesCollectsDiagnostics(t *testing.T) {
	t.Parallel()

	assertValidateOutput(
		t,
		filepath.Join("testdata", "config.graph.invalid.raw.json"),
		outcome.ExitCodeValidationBlocked,
		"validate_graph_invalid_output.golden",
	)
}

func TestFilesystemLoaderDeterministicOutputAcrossRuns(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.raw.json")
	code1, out1, err1 := runCommandWithFactory(
		[]string{"list", "reports", "--config", configPath},
		fsLoaderFactory,
		validate.AppDataValidator{},
	)
	if code1 != 0 {
		t.Fatalf("expected first exit code 0, got %d", code1)
	}
	if err1 != "" {
		t.Fatalf("expected empty first stderr, got %q", err1)
	}

	code2, out2, err2 := runCommandWithFactory(
		[]string{"list", "reports", "--config", configPath},
		fsLoaderFactory,
		validate.AppDataValidator{},
	)
	if code2 != 0 {
		t.Fatalf("expected second exit code 0, got %d", code2)
	}
	if err2 != "" {
		t.Fatalf("expected empty second stderr, got %q", err2)
	}

	if out1 != out2 {
		t.Fatalf("non-deterministic stdout\nfirst: %q\nsecond: %q", out1, out2)
	}
}

func TestOversizedConfigRuntimeFailureDeterministic(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "oversized.cue")
	content := strings.Repeat("a", int(safety.MaxConfigFileBytes)+1)
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write oversized config failed: %v", err)
	}

	code1, stdout1, stderr1 := runCommandWithFactory(
		[]string{"validate", "--config", configPath},
		fsLoaderFactory,
		validate.AppDataValidator{},
	)
	code2, stdout2, stderr2 := runCommandWithFactory(
		[]string{"validate", "--config", configPath},
		fsLoaderFactory,
		validate.AppDataValidator{},
	)

	if code1 != outcome.ExitCodeRuntimeFailure || code2 != outcome.ExitCodeRuntimeFailure {
		t.Fatalf("expected runtime failure exit codes, got %d and %d", code1, code2)
	}
	if stdout1 != "" || stdout2 != "" {
		t.Fatalf("expected empty stdout, got %q and %q", stdout1, stdout2)
	}
	if stderr1 == "" || stderr2 == "" {
		t.Fatalf("expected non-empty stderr, got %q and %q", stderr1, stderr2)
	}
	if stderr1 != stderr2 {
		t.Fatalf("non-deterministic stderr\nfirst: %q\nsecond: %q", stderr1, stderr2)
	}
}
