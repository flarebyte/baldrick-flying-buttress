package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/app"
	"github.com/flarebyte/baldrick-flying-buttress/internal/diagnostics"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func TestValidateSuccessWarningsOnly(t *testing.T) {
	t.Parallel()

	loader := func() (app.RawApp, error) {
		return app.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw app.RawApp) (app.ValidatedApp, diagnostics.Report, error) {
		_ = raw
		return app.ValidatedApp{Name: "stub-app", Modules: []string{"core", "edge"}}, diagnostics.Report{
			Diagnostics: []diagnostics.Diagnostic{
				{
					Code:     "FBW01",
					Severity: diagnostics.SeverityWarning,
					Message:  "warning only",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"validate"}, &out, &errOut, loader, validator)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	want := "{\"diagnostics\":[{\"code\":\"FBW01\",\"severity\":\"warning\",\"message\":\"warning only\",\"path\":\"module.stub\"}]}\n"
	if out.String() != want {
		t.Fatalf("stdout mismatch\nwant: %q\n got: %q", want, out.String())
	}
	if errOut.String() != "" {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}
}

func TestValidateFailureErrorDiagnostic(t *testing.T) {
	t.Parallel()

	loader := func() (app.RawApp, error) {
		return app.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw app.RawApp) (app.ValidatedApp, diagnostics.Report, error) {
		_ = raw
		return app.ValidatedApp{Name: "stub-app", Modules: []string{"core", "edge"}}, diagnostics.Report{
			Diagnostics: []diagnostics.Diagnostic{
				{
					Code:     "FBE01",
					Severity: diagnostics.SeverityError,
					Message:  "error diagnostic",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"validate"}, &out, &errOut, loader, validator)

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	want := "{\"diagnostics\":[{\"code\":\"FBE01\",\"severity\":\"error\",\"message\":\"error diagnostic\",\"path\":\"module.stub\"}]}\n"
	if out.String() != want {
		t.Fatalf("stdout mismatch\nwant: %q\n got: %q", want, out.String())
	}
	if errOut.String() != "" {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}
}

func TestValidateGoldenOutput(t *testing.T) {
	t.Parallel()

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"validate"}, &out, &errOut, validate.LoadStub, validate.ValidateStub)

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if errOut.String() != "" {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}

	goldenPath := filepath.Join("testdata", "validate_output.golden")
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if out.String() != string(want) {
		t.Fatalf("golden mismatch\nwant: %q\n got: %q", string(want), out.String())
	}
}

func TestValidateDeterministicOutputAcrossRuns(t *testing.T) {
	t.Parallel()

	var out1, errOut1 bytesBuffer
	exitCode1 := Execute([]string{"validate"}, &out1, &errOut1, validate.LoadStub, validate.ValidateStub)
	if exitCode1 != 1 {
		t.Fatalf("expected first exit code 1, got %d", exitCode1)
	}
	if errOut1.String() != "" {
		t.Fatalf("expected empty first stderr, got %q", errOut1.String())
	}

	var out2, errOut2 bytesBuffer
	exitCode2 := Execute([]string{"validate"}, &out2, &errOut2, validate.LoadStub, validate.ValidateStub)
	if exitCode2 != 1 {
		t.Fatalf("expected second exit code 1, got %d", exitCode2)
	}
	if errOut2.String() != "" {
		t.Fatalf("expected empty second stderr, got %q", errOut2.String())
	}

	if out1.String() != out2.String() {
		t.Fatalf("non-deterministic output\nfirst: %q\nsecond: %q", out1.String(), out2.String())
	}
}

func TestListReportsGoldenOutput(t *testing.T) {
	t.Parallel()

	loader := func() (app.RawApp, error) {
		return app.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw app.RawApp) (app.ValidatedApp, diagnostics.Report, error) {
		_ = raw
		return listValidatedApp(), diagnostics.Report{
			Diagnostics: []diagnostics.Diagnostic{
				{
					Code:     "FBW01",
					Severity: diagnostics.SeverityWarning,
					Message:  "warning only",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"list", "reports"}, &out, &errOut, loader, validator)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if errOut.String() != "" {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}

	goldenPath := filepath.Join("testdata", "list_reports_output.golden")
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if out.String() != string(want) {
		t.Fatalf("golden mismatch\nwant: %q\n got: %q", string(want), out.String())
	}
}

func TestListReportsBlockedOnErrorDiagnostic(t *testing.T) {
	t.Parallel()

	loader := func() (app.RawApp, error) {
		return app.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw app.RawApp) (app.ValidatedApp, diagnostics.Report, error) {
		_ = raw
		return listValidatedApp(), diagnostics.Report{
			Diagnostics: []diagnostics.Diagnostic{
				{
					Code:     "FBE01",
					Severity: diagnostics.SeverityError,
					Message:  "error diagnostic",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"list", "reports"}, &out, &errOut, loader, validator)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if out.String() != "" {
		t.Fatalf("expected empty stdout, got %q", out.String())
	}
	if errOut.String() != "" {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}
}

func TestListReportsDeterministicOutputAcrossRuns(t *testing.T) {
	t.Parallel()

	loader := func() (app.RawApp, error) {
		return app.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw app.RawApp) (app.ValidatedApp, diagnostics.Report, error) {
		_ = raw
		return listValidatedApp(), diagnostics.Report{}, nil
	}

	var out1, errOut1 bytesBuffer
	exitCode1 := Execute([]string{"list", "reports"}, &out1, &errOut1, loader, validator)
	if exitCode1 != 0 {
		t.Fatalf("expected first exit code 0, got %d", exitCode1)
	}
	if errOut1.String() != "" {
		t.Fatalf("expected empty first stderr, got %q", errOut1.String())
	}

	var out2, errOut2 bytesBuffer
	exitCode2 := Execute([]string{"list", "reports"}, &out2, &errOut2, loader, validator)
	if exitCode2 != 0 {
		t.Fatalf("expected second exit code 0, got %d", exitCode2)
	}
	if errOut2.String() != "" {
		t.Fatalf("expected empty second stderr, got %q", errOut2.String())
	}

	if out1.String() != out2.String() {
		t.Fatalf("non-deterministic output\nfirst: %q\nsecond: %q", out1.String(), out2.String())
	}
}

func TestGenerateJSONGoldenOutput(t *testing.T) {
	t.Parallel()

	loader := func() (app.RawApp, error) {
		return app.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw app.RawApp) (app.ValidatedApp, diagnostics.Report, error) {
		_ = raw
		return listValidatedApp(), diagnostics.Report{
			Diagnostics: []diagnostics.Diagnostic{
				{
					Code:     "FBW01",
					Severity: diagnostics.SeverityWarning,
					Message:  "warning only",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"generate", "json"}, &out, &errOut, loader, validator)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if errOut.String() != "" {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}

	goldenPath := filepath.Join("testdata", "generate_json_output.golden")
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if out.String() != string(want) {
		t.Fatalf("golden mismatch\nwant: %q\n got: %q", string(want), out.String())
	}
}

func TestGenerateJSONBlockedOnErrorDiagnostic(t *testing.T) {
	t.Parallel()

	loader := func() (app.RawApp, error) {
		return app.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw app.RawApp) (app.ValidatedApp, diagnostics.Report, error) {
		_ = raw
		return listValidatedApp(), diagnostics.Report{
			Diagnostics: []diagnostics.Diagnostic{
				{
					Code:     "FBE01",
					Severity: diagnostics.SeverityError,
					Message:  "error diagnostic",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"generate", "json"}, &out, &errOut, loader, validator)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if out.String() != "" {
		t.Fatalf("expected empty stdout, got %q", out.String())
	}
	if errOut.String() != "" {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}
}

func TestGenerateJSONDeterministicOutputAcrossRuns(t *testing.T) {
	t.Parallel()

	loader := func() (app.RawApp, error) {
		return app.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw app.RawApp) (app.ValidatedApp, diagnostics.Report, error) {
		_ = raw
		return listValidatedApp(), diagnostics.Report{}, nil
	}

	var out1, errOut1 bytesBuffer
	exitCode1 := Execute([]string{"generate", "json"}, &out1, &errOut1, loader, validator)
	if exitCode1 != 0 {
		t.Fatalf("expected first exit code 0, got %d", exitCode1)
	}
	if errOut1.String() != "" {
		t.Fatalf("expected empty first stderr, got %q", errOut1.String())
	}

	var out2, errOut2 bytesBuffer
	exitCode2 := Execute([]string{"generate", "json"}, &out2, &errOut2, loader, validator)
	if exitCode2 != 0 {
		t.Fatalf("expected second exit code 0, got %d", exitCode2)
	}
	if errOut2.String() != "" {
		t.Fatalf("expected empty second stderr, got %q", errOut2.String())
	}

	if out1.String() != out2.String() {
		t.Fatalf("non-deterministic output\nfirst: %q\nsecond: %q", out1.String(), out2.String())
	}
}

func listValidatedApp() app.ValidatedApp {
	return app.ValidatedApp{
		Name:    "stub-app",
		Modules: []string{"core", "edge"},
		Reports: []app.Report{
			{
				ID:    "cpu-overview",
				Title: "CPU Overview",
			},
			{
				ID:    "memory-health",
				Title: "Memory Health",
			},
		},
		Notes: []app.Note{
			{
				ID:    "n1",
				Label: "service.api",
			},
			{
				ID:    "n2",
				Label: "service.db",
			},
		},
		Relationships: []app.Relationship{
			{
				FromID: "n1",
				ToID:   "n2",
				Label:  "depends_on",
			},
		},
	}
}

type bytesBuffer struct {
	b []byte
}

func (b *bytesBuffer) Write(p []byte) (int, error) {
	b.b = append(b.b, p...)
	return len(p), nil
}

func (b *bytesBuffer) String() string {
	return string(b.b)
}
