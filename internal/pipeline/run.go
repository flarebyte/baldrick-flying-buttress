package pipeline

import (
	"errors"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

type Loader func() (domain.RawApp, error)
type Validator func(domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error)
type Action func(domain.ValidatedApp, domain.ValidationReport) error

func Run(loader Loader, validator Validator, action Action) error {
	if loader == nil {
		return errors.New("loader is required")
	}
	if validator == nil {
		return errors.New("validator is required")
	}
	if action == nil {
		return errors.New("action is required")
	}

	raw, err := loader()
	if err != nil {
		return err
	}

	validated, report, err := validator(raw)
	if err != nil {
		return err
	}

	return action(validated, report)
}
