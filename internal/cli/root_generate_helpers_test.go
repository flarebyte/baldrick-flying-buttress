package cli

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func setMarkdownReportWorkersForTest(workers int) func() {
	previous := int(atomic.LoadInt32(&markdownReportWorkers))
	atomic.StoreInt32(&markdownReportWorkers, int32(workers))
	return func() {
		atomic.StoreInt32(&markdownReportWorkers, int32(previous))
	}
}

func runGenerateMarkdownFixture(t *testing.T, fixtureName string) (string, int, string, string) {
	t.Helper()
	tmp := t.TempDir()
	configPath := writeFixtureConfig(t, tmp, fixtureName)
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}
	code, stdout, stderr := runCommandWithFactory([]string{"generate", "markdown", "--config", configPath}, loaderFactory, validator)
	return tmp, code, stdout, stderr
}

func assertGeneratedMarkdownGolden(t *testing.T, tmpDir, outputPath, goldenName string) {
	t.Helper()
	output, err := os.ReadFile(filepath.Join(tmpDir, outputPath))
	if err != nil {
		t.Fatalf("read generated report failed: %v", err)
	}
	want := readGolden(t, goldenName)
	if string(output) != want {
		t.Fatalf("generated markdown mismatch\nwant: %q\n got: %q", want, string(output))
	}
}
