package pipeline

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
)

func TestContractPipelineStageOrder(t *testing.T) {
	t.Parallel()

	steps := ""
	err := Run(
		LoaderFunc(func() (domain.RawApp, error) {
			steps += "L"
			return domain.RawApp{Source: "s"}, nil
		}),
		ValidatorFunc(func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
			_ = raw
			steps += "V"
			return domain.ValidatedApp{}, domain.ValidationReport{}, nil
		}),
		testAction{run: func(app domain.ValidatedApp, report domain.ValidationReport) error {
			_ = app
			_ = report
			steps += "A"
			return nil
		}},
	)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if steps != "LVA" {
		t.Fatalf("expected LVA, got %q", steps)
	}
}

func TestContractPipelineValidationBlockedShortCircuit(t *testing.T) {
	t.Parallel()

	actionCalled := false
	err := Run(
		LoaderFunc(func() (domain.RawApp, error) {
			return domain.RawApp{Source: "s"}, nil
		}),
		ValidatorFunc(func(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
			_ = raw
			return domain.ValidatedApp{}, domain.ValidationReport{
				Diagnostics: []domain.Diagnostic{{Code: "E", Severity: domain.SeverityError, Message: "m", Path: "p"}},
			}, nil
		}),
		testAction{
			run: func(app domain.ValidatedApp, report domain.ValidationReport) error {
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
