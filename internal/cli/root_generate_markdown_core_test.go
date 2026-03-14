package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
)

func TestGenerateMarkdownSuccess(t *testing.T) {
	t.Parallel()
	assertGenerateMarkdownSuccessFixture(t, "config.markdown.raw.json")
}

func TestGenerateMarkdownSuccessWithCueConfig(t *testing.T) {
	t.Parallel()
	assertGenerateMarkdownSuccessFixture(t, "config.markdown.cue")
}

func TestGenerateMarkdownCanShowNoteLabelsWhenRequested(t *testing.T) {
	t.Parallel()

	tmp, code, stdout, stderr := runGenerateMarkdownFixture(t, "config.markdown.showlabels.raw.json")
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "alpha.md"), "generate_markdown_showlabels_output.golden")
}

func TestGenerateMarkdownWithReportFilterWritesOnlySelectedReport(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.raw.json")
	code, stdout, stderr := runGenerateMarkdownWithArgs([]string{"--config", configPath, "--report", "alpha"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGolden(t, tmp, filepath.Join("out", "alpha.md"), "generate_markdown_alpha_output.golden")
	if _, err := os.Stat(filepath.Join(tmp, "out", "beta.md")); !os.IsNotExist(err) {
		t.Fatalf("expected beta.md to be absent, got stat err %v", err)
	}
}

func assertGenerateMarkdownSuccessFixture(t *testing.T, fixtureName string) {
	t.Helper()
	tmp, code, stdout, stderr := runGenerateMarkdownFixture(t, fixtureName)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")
	assertGeneratedMarkdownGoldens(t, tmp, []markdownGoldenExpectation{
		{outputPath: filepath.Join("out", "alpha.md"), goldenName: "generate_markdown_alpha_output.golden"},
		{outputPath: filepath.Join("out", "beta.md"), goldenName: "generate_markdown_beta_output.golden"},
	})
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
	code, stdout, stderr := runGenerateMarkdownWithConfig(configPath)
	assertRuntimeFailureOutput(t, code, stdout, stderr)
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

	code, stdout, stderr := runGenerateMarkdownWithConfig(configPath)
	assertRuntimeFailureOutput(t, code, stdout, stderr)

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
	code1, stdout1, stderr1 := runGenerateMarkdownWithConfig(configPath)
	if code1 != 0 {
		t.Fatalf("expected first exit code 0, got %d", code1)
	}
	alpha1 := readGeneratedMarkdown(t, tmp, filepath.Join("out", "alpha.md"))
	beta1 := readGeneratedMarkdown(t, tmp, filepath.Join("out", "beta.md"))

	code2, stdout2, stderr2 := runGenerateMarkdownWithConfig(configPath)
	if code2 != 0 {
		t.Fatalf("expected second exit code 0, got %d", code2)
	}
	alpha2 := readGeneratedMarkdown(t, tmp, filepath.Join("out", "alpha.md"))
	beta2 := readGeneratedMarkdown(t, tmp, filepath.Join("out", "beta.md"))

	if stdout1 != stdout2 || stderr1 != stderr2 {
		t.Fatalf("non-deterministic command output")
	}
	if alpha1 != alpha2 {
		t.Fatalf("non-deterministic alpha markdown\\nfirst: %q\\nsecond: %q", alpha1, alpha2)
	}
	if beta1 != beta2 {
		t.Fatalf("non-deterministic beta markdown\\nfirst: %q\\nsecond: %q", beta1, beta2)
	}
}

func TestGenerateMarkdownFullyOverwritesExistingReportFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.raw.json")
	alphaPath := filepath.Join(tmp, "out", "alpha.md")
	if err := os.MkdirAll(filepath.Dir(alphaPath), 0o755); err != nil {
		t.Fatalf("create alpha output dir failed: %v", err)
	}
	if err := os.WriteFile(alphaPath, []byte("STALE CONTENT THAT SHOULD NOT SURVIVE\n\nextra trailing bytes\n"), 0o644); err != nil {
		t.Fatalf("seed alpha report failed: %v", err)
	}

	code, stdout, stderr := runGenerateMarkdownWithConfig(configPath)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, "", "")

	want := readGolden(t, "generate_markdown_alpha_output.golden")
	got := readGeneratedMarkdown(t, tmp, filepath.Join("out", "alpha.md"))
	if got != want {
		t.Fatalf("expected full overwrite without stale content\nwant: %q\n got: %q", want, got)
	}
}

func TestGenerateMarkdownSingleWorkerEqualsMultiWorker(t *testing.T) {
	tmpSingle := t.TempDir()
	configSingle := writeFixtureConfig(t, tmpSingle, "config.markdown.raw.json")
	var code1 int
	var stdout1 string
	var stderr1 string
	func() {
		restoreSingle := setMarkdownReportWorkersForTest(1)
		defer restoreSingle()
		code1, stdout1, stderr1 = runGenerateMarkdownWithConfig(configSingle)
	}()
	if code1 != 0 {
		t.Fatalf("expected single-worker exit code 0, got %d", code1)
	}
	alpha1 := readGeneratedMarkdown(t, tmpSingle, filepath.Join("out", "alpha.md"))
	beta1 := readGeneratedMarkdown(t, tmpSingle, filepath.Join("out", "beta.md"))

	tmpMulti := t.TempDir()
	configMulti := writeFixtureConfig(t, tmpMulti, "config.markdown.raw.json")
	var code2 int
	var stdout2 string
	var stderr2 string
	func() {
		restoreMulti := setMarkdownReportWorkersForTest(4)
		defer restoreMulti()
		code2, stdout2, stderr2 = runGenerateMarkdownWithConfig(configMulti)
	}()
	if code2 != 0 {
		t.Fatalf("expected multi-worker exit code 0, got %d", code2)
	}
	alpha2 := readGeneratedMarkdown(t, tmpMulti, filepath.Join("out", "alpha.md"))
	beta2 := readGeneratedMarkdown(t, tmpMulti, filepath.Join("out", "beta.md"))

	if stdout1 != stdout2 || stderr1 != stderr2 {
		t.Fatalf("single-worker and multi-worker command output mismatch")
	}
	if alpha1 != alpha2 {
		t.Fatalf("single-worker and multi-worker alpha mismatch")
	}
	if beta1 != beta2 {
		t.Fatalf("single-worker and multi-worker beta mismatch")
	}
}
