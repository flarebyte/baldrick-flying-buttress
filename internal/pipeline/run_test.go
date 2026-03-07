package pipeline

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
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
			if raw.Source != wantRaw.Source {
				t.Fatalf("unexpected raw source: %#v", raw)
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

	err, actionCalled := runWithValidationError(t, false)
	assertValidationBlockedAndActionSkipped(t, err, actionCalled)
}
