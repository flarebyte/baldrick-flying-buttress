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
	code, stdout, stderr := runGenerateMarkdownWithConfig(configPath)
	return tmp, code, stdout, stderr
}

func runGenerateMarkdownBundleFixture(t *testing.T, fixtureName string, relativePaths []string) (string, int, string, string) {
	t.Helper()
	tmp := t.TempDir()
	configPath := writeFixtureBundle(t, tmp, fixtureName, relativePaths)
	code, stdout, stderr := runGenerateMarkdownWithConfig(configPath)
	return tmp, code, stdout, stderr
}

func runGenerateMarkdownWithConfig(configPath string) (int, string, string) {
	return runGenerateMarkdownWithArgs([]string{"--config", configPath})
}

func runGenerateMarkdownWithArgs(args []string) (int, string, string) {
	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	validator := validate.AppDataValidator{}
	return runCommandWithFactory(append([]string{"generate", "markdown"}, args...), loaderFactory, validator)
}

func readGeneratedMarkdown(t *testing.T, tmpDir, outputPath string) string {
	t.Helper()
	output, err := os.ReadFile(filepath.Join(tmpDir, outputPath))
	if err != nil {
		t.Fatalf("read generated report failed: %v", err)
	}
	return string(output)
}

func assertGeneratedMarkdownGolden(t *testing.T, tmpDir, outputPath, goldenName string) {
	t.Helper()
	output := readGeneratedMarkdown(t, tmpDir, outputPath)
	want := readGolden(t, goldenName)
	if output != want {
		t.Fatalf("generated markdown mismatch\nwant: %q\n got: %q", want, output)
	}
}

type markdownGoldenExpectation struct {
	outputPath string
	goldenName string
}

func assertGeneratedMarkdownGoldens(t *testing.T, tmpDir string, expectations []markdownGoldenExpectation) {
	t.Helper()
	for _, expectation := range expectations {
		assertGeneratedMarkdownGolden(t, tmpDir, expectation.outputPath, expectation.goldenName)
	}
}
