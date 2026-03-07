package pipeline

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
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
		testAction{run: func(domain.ValidatedApp, domain.ValidationReport) error {
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

	err, actionCalled := runWithValidationError(t, false)
	assertValidationBlockedAndActionSkipped(t, err, actionCalled)
}
