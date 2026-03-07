package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
)

func assertDeterministic(t *testing.T, args []string, loader pipeline.AppLoader, validator pipeline.AppValidator, wantCode int) {
	t.Helper()

	code1, out1, err1 := runCommand(args, loader, validator)
	if code1 != wantCode {
		t.Fatalf("expected first exit code %d, got %d", wantCode, code1)
	}
	code2, out2, err2 := runCommand(args, loader, validator)
	if code2 != wantCode {
		t.Fatalf("expected second exit code %d, got %d", wantCode, code2)
	}
	if out1 != out2 {
		t.Fatalf("non-deterministic stdout\nfirst: %q\nsecond: %q", out1, out2)
	}
	if err1 != err2 {
		t.Fatalf("non-deterministic stderr\nfirst: %q\nsecond: %q", err1, err2)
	}
}

func runCommand(args []string, loader pipeline.AppLoader, validator pipeline.AppValidator) (int, string, string) {
	return runCommandWithFactory(args, func(string) pipeline.AppLoader {
		return loader
	}, validator)
}

func runCommandWithFactory(args []string, loaderFactory LoaderFactory, validator pipeline.AppValidator) (int, string, string) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	code := ExecuteWithFactory(args, &out, &errOut, loaderFactory, validator)
	return code, out.String(), errOut.String()
}

func stubLoader() pipeline.AppLoader {
	return pipeline.LoaderFunc(func() (domain.RawApp, error) {
		return domain.RawApp{Source: "in-memory-stub"}, nil
	})
}

func validatorWith(app domain.ValidatedApp, report domain.ValidationReport, err error) pipeline.AppValidator {
	return pipeline.ValidatorFunc(func(domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		return app, report, err
	})
}

func readGolden(t *testing.T, filename string) string {
	t.Helper()
	p := filepath.Join("testdata", filename)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read golden %s: %v", filename, err)
	}
	return string(b)
}

func writeFixtureConfig(t *testing.T, dir string, fixtureName string) string {
	t.Helper()
	src := filepath.Join("testdata", fixtureName)
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read fixture %s: %v", fixtureName, err)
	}
	dst := filepath.Join(dir, fmt.Sprintf("%s.config.json", fixtureName))
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		t.Fatalf("write fixture config %s: %v", fixtureName, err)
	}
	return dst
}

func writeFixtureBundle(t *testing.T, dir string, fixtureName string, relativePaths []string) string {
	t.Helper()
	configPath := writeFixtureConfig(t, dir, fixtureName)
	for _, rel := range relativePaths {
		src := filepath.Join("testdata", rel)
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("read fixture bundle file %s: %v", rel, err)
		}
		dst := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			t.Fatalf("mkdir fixture bundle file %s: %v", rel, err)
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			t.Fatalf("write fixture bundle file %s: %v", rel, err)
		}
	}
	return configPath
}

func assertOutput(t *testing.T, gotStdout, gotStderr, wantStdout, wantStderr string) {
	t.Helper()
	if gotStdout != wantStdout {
		t.Fatalf("stdout mismatch\nwant: %q\n got: %q", wantStdout, gotStdout)
	}
	if gotStderr != wantStderr {
		t.Fatalf("stderr mismatch\nwant: %q\n got: %q", wantStderr, gotStderr)
	}
}

func listValidatedApp() domain.ValidatedApp {
	return appFixture(
		[]string{"core", "edge"},
		[]domain.Report{{ID: "cpu-overview", Title: "CPU Overview"}, {ID: "memory-health", Title: "Memory Health"}},
		[]domain.Note{{ID: "n1", Label: "service.api"}, {ID: "n2", Label: "service.db"}},
		[]domain.Relationship{{FromID: "n1", ToID: "n2", Label: "depends_on"}},
	)
}

func unorderedValidatedApp() domain.ValidatedApp {
	return appFixture(
		[]string{"edge", "core"},
		[]domain.Report{{ID: "memory-health", Title: "Memory Health"}, {ID: "cpu-overview", Title: "CPU Overview"}},
		[]domain.Note{{ID: "n2", Label: "service.db"}, {ID: "n1", Label: "service.api"}},
		[]domain.Relationship{{FromID: "n1", ToID: "n2", Label: "depends_on"}, {FromID: "n1", ToID: "n2", Label: "owns"}},
	)
}

func orderedValidatedAppForOrdering() domain.ValidatedApp {
	return appFixture(
		[]string{"core", "edge"},
		[]domain.Report{{ID: "cpu-overview", Title: "CPU Overview"}, {ID: "memory-health", Title: "Memory Health"}},
		[]domain.Note{{ID: "n1", Label: "service.api"}, {ID: "n2", Label: "service.db"}},
		[]domain.Relationship{{FromID: "n1", ToID: "n2", Label: "depends_on"}, {FromID: "n1", ToID: "n2", Label: "owns"}},
	)
}

func appFixture(modules []string, reports []domain.Report, notes []domain.Note, relationships []domain.Relationship) domain.ValidatedApp {
	return domain.ValidatedApp{
		Name:          "stub-app",
		Modules:       modules,
		Reports:       reports,
		Notes:         notes,
		Relationships: relationships,
	}
}

func warningOnlyReport() domain.ValidationReport {
	return domain.ValidationReport{
		Diagnostics: []domain.Diagnostic{{
			Code:     "FBW01",
			Severity: domain.SeverityWarning,
			Message:  "warning only",
			Path:     "module.stub",
		}},
	}
}

func orphansValidatedApp() domain.ValidatedApp {
	return domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "n.ingredient.a", Label: "ingredient"},
			{ID: "n.ingredient.orphan", Label: "ingredient"},
			{ID: "n.tool.a", Label: "tool"},
		},
		Relationships: []domain.Relationship{
			{FromID: "n.ingredient.a", ToID: "n.tool.a", Label: "uses"},
		},
	}
}

func errorOnlyReport() domain.ValidationReport {
	return domain.ValidationReport{
		Diagnostics: []domain.Diagnostic{{
			Code:     "FBE01",
			Severity: domain.SeverityError,
			Message:  "error diagnostic",
			Path:     "module.stub",
		}},
	}
}
