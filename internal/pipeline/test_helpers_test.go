package pipeline

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
)

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

func runWithValidationError(t *testing.T, allowOnValidationErrors bool) (error, bool) {
	t.Helper()
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
			run: func(domain.ValidatedApp, domain.ValidationReport) error {
				actionCalled = true
				return nil
			},
			allowOnValidationErrors: allowOnValidationErrors,
		},
	)
	return err, actionCalled
}

func assertValidationBlockedAndActionSkipped(t *testing.T, err error, actionCalled bool) {
	t.Helper()
	if !outcome.IsValidationBlocked(err) {
		t.Fatalf("expected validation blocked error, got %v", err)
	}
	if actionCalled {
		t.Fatal("expected action to be skipped")
	}
}
