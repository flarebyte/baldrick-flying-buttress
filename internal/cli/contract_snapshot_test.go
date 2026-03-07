package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func TestContractSnapshotValidateOutput(t *testing.T) {
	t.Parallel()
	assertContractCommandOutput(t,
		[]string{"validate", "--config", filepath.Join("testdata", "config.raw.json")},
		readGolden(t, "validate_output.golden"),
		"",
	)
}

func TestContractSnapshotListReportsOutput(t *testing.T) {
	t.Parallel()
	assertContractCommandOutput(t,
		[]string{"list", "reports", "--config", filepath.Join("testdata", "config.raw.json")},
		readGolden(t, "list_reports_output.golden"),
		"",
	)
}

func TestContractSnapshotListNamesOutput(t *testing.T) {
	t.Parallel()
	assertContractCommandOutput(t,
		[]string{"list", "names", "--config", filepath.Join("..", "..", "testdata", "names.raw.json"), "--prefix", "cli."},
		readGolden(t, "list_names_table_output.golden"),
		"",
	)
}

func TestContractSnapshotLintNamesOutput(t *testing.T) {
	t.Parallel()
	assertContractCommandOutput(t,
		[]string{"lint", "names", "--config", filepath.Join("testdata", "config.lint.raw.json")},
		readGolden(t, "lint_names_dot_output.golden"),
		"",
	)
}

func TestContractSnapshotLintOrphansOutput(t *testing.T) {
	t.Parallel()
	assertContractCommandOutput(t,
		[]string{"lint", "orphans", "--config", filepath.Join("testdata", "config.orphans.raw.json"), "--subject-label", "ingredient"},
		readGolden(t, "lint_orphans_default_output.golden"),
		"",
	)
}

func TestContractSnapshotGenerateJSONOutput(t *testing.T) {
	t.Parallel()
	assertContractCommandOutput(t,
		[]string{"generate", "json", "--config", filepath.Join("testdata", "config.raw.json")},
		readGolden(t, "generate_json_output.golden"),
		"",
	)
}

func TestContractSnapshotGenerateMarkdownFiles(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.markdown.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if stdout != "" || stderr != "" {
		t.Fatalf("unexpected output stdout=%q stderr=%q", stdout, stderr)
	}

	alpha, err := os.ReadFile(filepath.Join(tmp, "out", "alpha.md"))
	if err != nil {
		t.Fatalf("read alpha failed: %v", err)
	}
	beta, err := os.ReadFile(filepath.Join(tmp, "out", "beta.md"))
	if err != nil {
		t.Fatalf("read beta failed: %v", err)
	}
	if string(alpha) != readGolden(t, "generate_markdown_alpha_output.golden") {
		t.Fatalf("alpha snapshot mismatch")
	}
	if string(beta) != readGolden(t, "generate_markdown_beta_output.golden") {
		t.Fatalf("beta snapshot mismatch")
	}
}

func assertContractCommandOutput(t *testing.T, args []string, wantStdout string, wantStderr string) {
	t.Helper()
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code, stdout, stderr := runCommandWithFactory(args, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if stdout != wantStdout {
		t.Fatalf("stdout snapshot mismatch\nwant: %q\ngot: %q", wantStdout, stdout)
	}
	if stderr != wantStderr {
		t.Fatalf("stderr snapshot mismatch\nwant: %q\ngot: %q", wantStderr, stderr)
	}
}
