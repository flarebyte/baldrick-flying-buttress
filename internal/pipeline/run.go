package pipeline

import (
	"errors"

	"github.com/flarebyte/baldrick-flying-buttress/internal/app"
	"github.com/flarebyte/baldrick-flying-buttress/internal/diagnostics"
)

type Loader func() (app.RawApp, error)
type Validator func(app.RawApp) (app.ValidatedApp, diagnostics.Report, error)
type Action func(app.ValidatedApp, diagnostics.Report) error

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
