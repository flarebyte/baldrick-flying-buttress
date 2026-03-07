package cli

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func TestValidateSuccessWarningsOnly(t *testing.T) {
	t.Parallel()

	loader := func() (domain.RawApp, error) {
		return domain.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		_ = raw
		return domain.ValidatedApp{Name: "stub-app", Modules: []string{"core", "edge"}}, domain.ValidationReport{
			Diagnostics: []domain.Diagnostic{
				{
					Code:     "FBW01",
					Severity: domain.SeverityWarning,
					Message:  "warning only",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"validate"}, &out, &errOut, pipeline.LoaderFunc(loader), pipeline.ValidatorFunc(validator))

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

	loader := func() (domain.RawApp, error) {
		return domain.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		_ = raw
		return domain.ValidatedApp{Name: "stub-app", Modules: []string{"core", "edge"}}, domain.ValidationReport{
			Diagnostics: []domain.Diagnostic{
				{
					Code:     "FBE01",
					Severity: domain.SeverityError,
					Message:  "error diagnostic",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"validate"}, &out, &errOut, pipeline.LoaderFunc(loader), pipeline.ValidatorFunc(validator))

	if exitCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, exitCode)
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
	exitCode := Execute([]string{"validate"}, &out, &errOut, validate.StubAppLoader{}, validate.StubAppValidator{})

	if exitCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, exitCode)
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
	exitCode1 := Execute([]string{"validate"}, &out1, &errOut1, validate.StubAppLoader{}, validate.StubAppValidator{})
	if exitCode1 != 1 {
		t.Fatalf("expected first exit code 1, got %d", exitCode1)
	}
	if errOut1.String() != "" {
		t.Fatalf("expected empty first stderr, got %q", errOut1.String())
	}

	var out2, errOut2 bytesBuffer
	exitCode2 := Execute([]string{"validate"}, &out2, &errOut2, validate.StubAppLoader{}, validate.StubAppValidator{})
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

	loader := func() (domain.RawApp, error) {
		return domain.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		_ = raw
		return listValidatedApp(), domain.ValidationReport{
			Diagnostics: []domain.Diagnostic{
				{
					Code:     "FBW01",
					Severity: domain.SeverityWarning,
					Message:  "warning only",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"list", "reports"}, &out, &errOut, pipeline.LoaderFunc(loader), pipeline.ValidatorFunc(validator))
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

	loader := func() (domain.RawApp, error) {
		return domain.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		_ = raw
		return listValidatedApp(), domain.ValidationReport{
			Diagnostics: []domain.Diagnostic{
				{
					Code:     "FBE01",
					Severity: domain.SeverityError,
					Message:  "error diagnostic",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"list", "reports"}, &out, &errOut, pipeline.LoaderFunc(loader), pipeline.ValidatorFunc(validator))
	if exitCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, exitCode)
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

	loader := func() (domain.RawApp, error) {
		return domain.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		_ = raw
		return listValidatedApp(), domain.ValidationReport{}, nil
	}

	var out1, errOut1 bytesBuffer
	exitCode1 := Execute([]string{"list", "reports"}, &out1, &errOut1, pipeline.LoaderFunc(loader), pipeline.ValidatorFunc(validator))
	if exitCode1 != 0 {
		t.Fatalf("expected first exit code 0, got %d", exitCode1)
	}
	if errOut1.String() != "" {
		t.Fatalf("expected empty first stderr, got %q", errOut1.String())
	}

	var out2, errOut2 bytesBuffer
	exitCode2 := Execute([]string{"list", "reports"}, &out2, &errOut2, pipeline.LoaderFunc(loader), pipeline.ValidatorFunc(validator))
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

	loader := func() (domain.RawApp, error) {
		return domain.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		_ = raw
		return listValidatedApp(), domain.ValidationReport{
			Diagnostics: []domain.Diagnostic{
				{
					Code:     "FBW01",
					Severity: domain.SeverityWarning,
					Message:  "warning only",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"generate", "json"}, &out, &errOut, pipeline.LoaderFunc(loader), pipeline.ValidatorFunc(validator))
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

	loader := func() (domain.RawApp, error) {
		return domain.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		_ = raw
		return listValidatedApp(), domain.ValidationReport{
			Diagnostics: []domain.Diagnostic{
				{
					Code:     "FBE01",
					Severity: domain.SeverityError,
					Message:  "error diagnostic",
					Path:     "module.stub",
				},
			},
		}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"generate", "json"}, &out, &errOut, pipeline.LoaderFunc(loader), pipeline.ValidatorFunc(validator))
	if exitCode != outcome.ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeValidationBlocked, exitCode)
	}
	if out.String() != "" {
		t.Fatalf("expected empty stdout, got %q", out.String())
	}
	if errOut.String() != "" {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}
}

func TestRuntimeFailureMapsToDistinctExitCodeAndStderr(t *testing.T) {
	t.Parallel()

	loader := func() (domain.RawApp, error) {
		return domain.RawApp{}, errors.New("runtime exploded")
	}
	validator := func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		_ = raw
		return domain.ValidatedApp{}, domain.ValidationReport{}, nil
	}

	var out, errOut bytesBuffer
	exitCode := Execute([]string{"validate"}, &out, &errOut, pipeline.LoaderFunc(loader), pipeline.ValidatorFunc(validator))

	if exitCode != outcome.ExitCodeRuntimeFailure {
		t.Fatalf("expected exit code %d, got %d", outcome.ExitCodeRuntimeFailure, exitCode)
	}
	if out.String() != "" {
		t.Fatalf("expected empty stdout, got %q", out.String())
	}
	if errOut.String() != "runtime exploded\n" {
		t.Fatalf("expected stderr runtime message, got %q", errOut.String())
	}
}

func TestGenerateJSONDeterministicOutputAcrossRuns(t *testing.T) {
	t.Parallel()

	loader := func() (domain.RawApp, error) {
		return domain.RawApp{Source: "in-memory-stub"}, nil
	}
	validator := func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
		_ = raw
		return listValidatedApp(), domain.ValidationReport{}, nil
	}

	var out1, errOut1 bytesBuffer
	exitCode1 := Execute([]string{"generate", "json"}, &out1, &errOut1, pipeline.LoaderFunc(loader), pipeline.ValidatorFunc(validator))
	if exitCode1 != 0 {
		t.Fatalf("expected first exit code 0, got %d", exitCode1)
	}
	if errOut1.String() != "" {
		t.Fatalf("expected empty first stderr, got %q", errOut1.String())
	}

	var out2, errOut2 bytesBuffer
	exitCode2 := Execute([]string{"generate", "json"}, &out2, &errOut2, pipeline.LoaderFunc(loader), pipeline.ValidatorFunc(validator))
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

func listValidatedApp() domain.ValidatedApp {
	return domain.ValidatedApp{
		Name:    "stub-app",
		Modules: []string{"core", "edge"},
		Reports: []domain.Report{
			{
				ID:    "cpu-overview",
				Title: "CPU Overview",
			},
			{
				ID:    "memory-health",
				Title: "Memory Health",
			},
		},
		Notes: []domain.Note{
			{
				ID:    "n1",
				Label: "service.api",
			},
			{
				ID:    "n2",
				Label: "service.db",
			},
		},
		Relationships: []domain.Relationship{
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
