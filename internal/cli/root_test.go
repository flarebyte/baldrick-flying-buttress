package cli

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/olivier/baldrick-flying-buttress/internal/diagnostics"
	"github.com/olivier/baldrick-flying-buttress/internal/validate"
)

func TestValidateSuccessWarningsOnly(t *testing.T) {
	t.Parallel()

	runner := func(ctx context.Context) (diagnostics.Report, error) {
		_ = ctx
		return diagnostics.Report{
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
	exitCode := Execute([]string{"validate"}, &out, &errOut, runner)

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

	runner := func(ctx context.Context) (diagnostics.Report, error) {
		_ = ctx
		return diagnostics.Report{
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
	exitCode := Execute([]string{"validate"}, &out, &errOut, runner)

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
	exitCode := Execute([]string{"validate"}, &out, &errOut, validate.RunStub)

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
