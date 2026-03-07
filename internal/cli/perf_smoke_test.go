package cli

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func TestPerfSmokeValidateModerateFixture(t *testing.T) {
	if testing.Short() {
		t.Skip("skip perf smoke in short mode")
	}

	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}
	configPath := filepath.Join("testdata", "config.perf.raw.json")

	code, stdout, stderr := runCommandWithFactory([]string{"validate", "--config", configPath}, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%q", code, stderr)
	}
	if stdout != "{\"diagnostics\":[]}\n" {
		t.Fatalf("unexpected stdout: %q", stdout)
	}
	if stderr != "" {
		t.Fatalf("unexpected stderr: %q", stderr)
	}
}

func TestPerfSmokeGenerateJSONModerateFixture(t *testing.T) {
	if testing.Short() {
		t.Skip("skip perf smoke in short mode")
	}

	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}
	configPath := filepath.Join("testdata", "config.perf.raw.json")

	code1, stdout1, stderr1 := runCommandWithFactory([]string{"generate", "json", "--config", configPath}, loaderFactory, validator)
	if code1 != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%q", code1, stderr1)
	}
	code2, stdout2, stderr2 := runCommandWithFactory([]string{"generate", "json", "--config", configPath}, loaderFactory, validator)
	if code2 != code1 {
		t.Fatalf("exit code mismatch: got %d want %d", code2, code1)
	}
	if stdout2 != stdout1 {
		t.Fatal("non-deterministic json output")
	}
	if stderr1 != "" || stderr2 != "" {
		t.Fatalf("unexpected stderr values: first=%q second=%q", stderr1, stderr2)
	}
}

func TestPerfSmokeGenerateMarkdownModerateFixture(t *testing.T) {
	if testing.Short() {
		t.Skip("skip perf smoke in short mode")
	}

	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, "config.perf.raw.json")
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}

	code1, stdout1, stderr1 := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code1 != 0 {
		t.Fatalf("expected first exit code 0, got %d, stderr=%q", code1, stderr1)
	}
	if stdout1 != "" || stderr1 != "" {
		t.Fatalf("unexpected first command output stdout=%q stderr=%q", stdout1, stderr1)
	}

	firstBytes := make([][]byte, 0, 12)
	for i := 0; i < 12; i++ {
		b, err := os.ReadFile(filepath.Join(tmp, "out", "perf-"+strconv.Itoa(i)+".md"))
		if err != nil {
			t.Fatalf("read first report perf-%d.md failed: %v", i, err)
		}
		firstBytes = append(firstBytes, b)
	}

	code2, stdout2, stderr2 := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	if code2 != code1 {
		t.Fatalf("exit code mismatch: got %d want %d", code2, code1)
	}
	if stdout2 != stdout1 || stderr2 != stderr1 {
		t.Fatal("non-deterministic markdown command output")
	}

	for i := 0; i < 12; i++ {
		b, err := os.ReadFile(filepath.Join(tmp, "out", "perf-"+strconv.Itoa(i)+".md"))
		if err != nil {
			t.Fatalf("read second report perf-%d.md failed: %v", i, err)
		}
		if string(b) != string(firstBytes[i]) {
			t.Fatalf("non-deterministic report bytes for perf-%d.md", i)
		}
	}
}
