package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func TestGenerateMarkdownSuccess(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	alpha, err := os.ReadFile(filepath.Join(tmp, "out", "alpha.md"))
	if err != nil {
		t.Fatalf("read alpha report failed: %v", err)
	}
	beta, err := os.ReadFile(filepath.Join(tmp, "out", "beta.md"))
	if err != nil {
		t.Fatalf("read beta report failed: %v", err)
	}
	if string(alpha) != readGolden(t, "generate_markdown_alpha_output.golden") {
		t.Fatalf("alpha markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_alpha_output.golden"), string(alpha))
	}
	if string(beta) != readGolden(t, "generate_markdown_beta_output.golden") {
		t.Fatalf("beta markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_beta_output.golden"), string(beta))
	}
}

func TestGenerateMarkdownSuccessWithCueConfig(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.cue")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	alpha, err := os.ReadFile(filepath.Join(tmp, "out", "alpha.md"))
	if err != nil {
		t.Fatalf("read alpha report failed: %v", err)
	}
	beta, err := os.ReadFile(filepath.Join(tmp, "out", "beta.md"))
	if err != nil {
		t.Fatalf("read beta report failed: %v", err)
	}
	if string(alpha) != readGolden(t, "generate_markdown_alpha_output.golden") {
		t.Fatalf("alpha markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_alpha_output.golden"), string(alpha))
	}
	if string(beta) != readGolden(t, "generate_markdown_beta_output.golden") {
		t.Fatalf("beta markdown mismatch\\nwant: %q\\n got: %q", readGolden(t, "generate_markdown_beta_output.golden"), string(beta))
	}
}

func TestGenerateMarkdownBlockedOnErrorDiagnostic(t *testing.T) {
	t.Parallel()

	app := domain.ValidatedApp{
		ConfigDir:       t.TempDir(),
		MarkdownReports: []domain.MarkdownReport{{Title: "Blocked", Filepath: "out/blocked.md"}},
	}
	code, stdout, stderr := runCommand([]string{"generate", "markdown"}, stubLoader(), validatorWith(app, errorOnlyReport(), nil))
	if code != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, code)
	}
	assertOutput(t, stdout, stderr, "", "")
}

func TestGenerateMarkdownRuntimeFailureOnUnwritablePath(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, "out"), 0o755); err != nil {
		t.Fatalf("create out dir failed: %v", err)
	}
	configPath := filepath.Join(tmp, "config.raw.json")
	content := `{"source":"x","name":"x","modules":[],"reports":[{"title":"R","filepath":"out","sections":[{"title":"H2","sections":[{"title":"H3","notes":["n1"]}]}]}],"notes":[{"name":"n1","title":"N1","markdown":"Body"}],"relationships":[]}`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
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

func TestGenerateMarkdownFailureLeavesExistingFileUnchanged(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	outDir := filepath.Join(tmp, "out")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("create out dir failed: %v", err)
	}
	reportPath := filepath.Join(outDir, "report.md")
	original := "ORIGINAL\n"
	if err := os.WriteFile(reportPath, []byte(original), 0o644); err != nil {
		t.Fatalf("seed report failed: %v", err)
	}

	configPath := filepath.Join(tmp, "config.raw.json")
	content := `{"source":"x","name":"x","modules":[],"reports":[{"title":"R","filepath":"out/report.md","sections":[{"title":"H2","sections":[{"title":"H3","notes":["n1"]}]}]}],"notes":[{"name":"n1","title":"N1","markdown":"Body"}],"relationships":[]}`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}

	if err := os.Chmod(outDir, 0o555); err != nil {
		t.Fatalf("chmod out dir failed: %v", err)
	}
	defer func() {
		_ = os.Chmod(outDir, 0o755)
	}()

	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != outcome.ExitCodeRuntimeFailure {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeRuntimeFailure, code)
	}
	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if stderr == "" {
		t.Fatal("expected stderr")
	}

	got, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("read report failed: %v", err)
	}
	if string(got) != original {
		t.Fatalf("expected original report to remain unchanged\\nwant: %q\\ngot: %q", original, string(got))
	}
}

func TestGenerateMarkdownDeterministicAcrossRuns(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code1, stdout1, stderr1 := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code1 != 0 {
		t.Fatalf("expected first exit code 0, got %d", code1)
	}
	alpha1, err := os.ReadFile(filepath.Join(tmp, "out", "alpha.md"))
	if err != nil {
		t.Fatalf("read alpha first run failed: %v", err)
	}
	beta1, err := os.ReadFile(filepath.Join(tmp, "out", "beta.md"))
	if err != nil {
		t.Fatalf("read beta first run failed: %v", err)
	}

	code2, stdout2, stderr2 := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code2 != 0 {
		t.Fatalf("expected second exit code 0, got %d", code2)
	}
	alpha2, err := os.ReadFile(filepath.Join(tmp, "out", "alpha.md"))
	if err != nil {
		t.Fatalf("read alpha second run failed: %v", err)
	}
	beta2, err := os.ReadFile(filepath.Join(tmp, "out", "beta.md"))
	if err != nil {
		t.Fatalf("read beta second run failed: %v", err)
	}

	if stdout1 != stdout2 || stderr1 != stderr2 {
		t.Fatalf("non-deterministic command output")
	}
	if string(alpha1) != string(alpha2) {
		t.Fatalf("non-deterministic alpha markdown\\nfirst: %q\\nsecond: %q", string(alpha1), string(alpha2))
	}
	if string(beta1) != string(beta2) {
		t.Fatalf("non-deterministic beta markdown\\nfirst: %q\\nsecond: %q", string(beta1), string(beta2))
	}
}

func TestGenerateMarkdownSingleWorkerEqualsMultiWorker(t *testing.T) {
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	tmpSingle := t.TempDir()
	configSingle := writeFixtureConfig(t, tmpSingle, "config.markdown.raw.json")
	var code1 int
	var stdout1 string
	var stderr1 string
	func() {
		restoreSingle := setMarkdownReportWorkersForTest(1)
		defer restoreSingle()
		code1, stdout1, stderr1 = runCommandWithFactory([]string{"generate", "markdown", "--config", configSingle}, loaderFactory, validator)
	}()
	if code1 != 0 {
		t.Fatalf("expected single-worker exit code 0, got %d", code1)
	}
	alpha1, err := os.ReadFile(filepath.Join(tmpSingle, "out", "alpha.md"))
	if err != nil {
		t.Fatalf("read single-worker alpha failed: %v", err)
	}
	beta1, err := os.ReadFile(filepath.Join(tmpSingle, "out", "beta.md"))
	if err != nil {
		t.Fatalf("read single-worker beta failed: %v", err)
	}

	tmpMulti := t.TempDir()
	configMulti := writeFixtureConfig(t, tmpMulti, "config.markdown.raw.json")
	var code2 int
	var stdout2 string
	var stderr2 string
	func() {
		restoreMulti := setMarkdownReportWorkersForTest(4)
		defer restoreMulti()
		code2, stdout2, stderr2 = runCommandWithFactory([]string{"generate", "markdown", "--config", configMulti}, loaderFactory, validator)
	}()
	if code2 != 0 {
		t.Fatalf("expected multi-worker exit code 0, got %d", code2)
	}
	alpha2, err := os.ReadFile(filepath.Join(tmpMulti, "out", "alpha.md"))
	if err != nil {
		t.Fatalf("read multi-worker alpha failed: %v", err)
	}
	beta2, err := os.ReadFile(filepath.Join(tmpMulti, "out", "beta.md"))
	if err != nil {
		t.Fatalf("read multi-worker beta failed: %v", err)
	}

	if stdout1 != stdout2 || stderr1 != stderr2 {
		t.Fatalf("single-worker and multi-worker command output mismatch")
	}
	if string(alpha1) != string(alpha2) {
		t.Fatalf("single-worker and multi-worker alpha mismatch")
	}
	if string(beta1) != string(beta2) {
		t.Fatalf("single-worker and multi-worker beta mismatch")
	}
}
