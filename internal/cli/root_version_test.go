package cli

import (
	"bytes"
	"context"
	"runtime"
	"strings"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/buildinfo"
	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
)

func TestRootVersionFlagOutputsJSON(t *testing.T) {
	t.Parallel()

	oldVersion := buildinfo.Version
	oldCommit := buildinfo.Commit
	oldDate := buildinfo.Date
	buildinfo.Version = "1.2.3"
	buildinfo.Commit = "abc123def456"
	buildinfo.Date = "2026-03-07T19:15:00Z"
	t.Cleanup(func() {
		buildinfo.Version = oldVersion
		buildinfo.Commit = oldCommit
		buildinfo.Date = oldDate
	})

	loaderFactory := func(string) pipeline.AppLoader {
		return pipeline.LoaderFunc(func(_ context.Context) (domain.RawApp, error) {
			t.Fatal("loader should not be called for --version")
			return domain.RawApp{}, nil
		})
	}
	validator := pipeline.ValidatorFunc(func(_ context.Context, _ domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		t.Fatal("validator should not be called for --version")
		return domain.ValidatedApp{}, domain.ValidationReport{}, nil
	})

	var out bytes.Buffer
	var errOut bytes.Buffer
	code := ExecuteWithFactory([]string{"--version"}, &out, &errOut, loaderFactory, validator)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	want := "{\"version\":\"1.2.3\",\"commit\":\"abc123def456\",\"date\":\"2026-03-07T19:15:00Z\",\"goVersion\":\"" + runtime.Version() + "\",\"goos\":\"" + runtime.GOOS + "\",\"goarch\":\"" + runtime.GOARCH + "\"}\n"
	if out.String() != want {
		t.Fatalf("stdout mismatch\nwant: %q\n got: %q", want, out.String())
	}
	if errOut.String() != "" {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}
}

func TestRootWithoutArgsShowsDescriptionAndUsage(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runCommand(
		[]string{},
		stubLoader(),
		validatorWith(domain.ValidatedApp{}, domain.ValidationReport{}, nil),
	)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
	if !strings.Contains(stdout, "flyb validates application graph configuration") {
		t.Fatalf("expected description in help output, got %q", stdout)
	}
	if !strings.Contains(stdout, "Available Commands:") {
		t.Fatalf("expected command usage output, got %q", stdout)
	}
}
