package cli

import (
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func TestExportCueGoldenOutput(t *testing.T) {
	t.Parallel()

	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	exitCode, stdout, stderr := runCommandWithFactory(
		[]string{"export", "cue", "--config", filepath.Join("testdata", "config.export.raw.json")},
		loaderFactory,
		validate.AppDataValidator{},
	)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "export_cue_output.golden"), "")
}

func TestExportCueWithReportFilter(t *testing.T) {
	t.Parallel()

	loaderFactory := func(path string) pipeline.AppLoader { return load.FSAppLoader{ConfigPath: path} }
	code, stdout, stderr := runCommandWithFactory(
		[]string{"export", "cue", "--config", filepath.Join("testdata", "config.markdown.raw.json"), "--report", "alpha"},
		loaderFactory,
		validate.AppDataValidator{},
	)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	assertOutput(t, stdout, stderr, readGolden(t, "export_cue_alpha_output.golden"), "")
}
