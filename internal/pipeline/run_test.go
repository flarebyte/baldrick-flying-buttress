package pipeline

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
)

func TestRunCallsLoadValidateAction(t *testing.T) {
	t.Parallel()

	var calls []string
	wantRaw := domain.RawApp{Source: "raw-stub"}
	wantValidated := domain.ValidatedApp{Name: "validated", Modules: []string{"core"}}
	wantReport := domain.ValidationReport{Diagnostics: []domain.Diagnostic{{
		Code:     "FBW01",
		Severity: domain.SeverityWarning,
		Message:  "warning",
		Path:     "module.stub",
	}}}

	err := Run(
		LoaderFunc(func() (domain.RawApp, error) {
			calls = append(calls, "load")
			return wantRaw, nil
		}),
		ValidatorFunc(func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
			calls = append(calls, "validate")
			if raw != wantRaw {
				t.Fatalf("unexpected raw app: %#v", raw)
			}
			return wantValidated, wantReport, nil
		}),
		testAction{
			run: func(validated domain.ValidatedApp, report domain.ValidationReport) error {
				calls = append(calls, "action")
				if validated.Name != wantValidated.Name {
					t.Fatalf("unexpected app name: %s", validated.Name)
				}
				if len(validated.Modules) != len(wantValidated.Modules) || validated.Modules[0] != wantValidated.Modules[0] {
					t.Fatalf("unexpected modules: %#v", validated.Modules)
				}
				if len(report.Diagnostics) != len(wantReport.Diagnostics) || report.Diagnostics[0] != wantReport.Diagnostics[0] {
					t.Fatalf("unexpected report: %#v", report)
				}
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	if len(calls) != 3 || calls[0] != "load" || calls[1] != "validate" || calls[2] != "action" {
		t.Fatalf("unexpected call order: %#v", calls)
	}
}

func TestRunShortCircuitsActionOnValidationErrorDiagnostic(t *testing.T) {
	t.Parallel()

	actionCalled := false

	err := Run(
		LoaderFunc(func() (domain.RawApp, error) {
			return domain.RawApp{Source: "raw-stub"}, nil
		}),
		ValidatorFunc(func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
			_ = raw
			return domain.ValidatedApp{}, domain.ValidationReport{
				Diagnostics: []domain.Diagnostic{{
					Code:     "FBE01",
					Severity: domain.SeverityError,
					Message:  "error",
					Path:     "module.stub",
				}},
			}, nil
		}),
		testAction{
			run: func(validated domain.ValidatedApp, report domain.ValidationReport) error {
				actionCalled = true
				return nil
			},
			allowOnValidationErrors: false,
		},
	)
	if !outcome.IsValidationBlocked(err) {
		t.Fatalf("expected validation blocked error, got %v", err)
	}
	if actionCalled {
		t.Fatal("expected action to be skipped")
	}
}

type testAction struct {
	run                     func(domain.ValidatedApp, domain.ValidationReport) error
	allowOnValidationErrors bool
}

func (a testAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
	if a.run == nil {
		return nil
	}
	return a.run(validated, report)
}

func (a testAction) AllowOnValidationErrors() bool {
	return a.allowOnValidationErrors
}
