package pipeline

import (
	"errors"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
)

type AppLoader interface {
	Load() (domain.RawApp, error)
}

type AppValidator interface {
	Validate(domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error)
}

type CommandAction interface {
	Execute(domain.ValidatedApp, domain.ValidationReport) error
	AllowOnValidationErrors() bool
}

type LoaderFunc func() (domain.RawApp, error)

func (f LoaderFunc) Load() (domain.RawApp, error) {
	return f()
}

type ValidatorFunc func(domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error)

func (f ValidatorFunc) Validate(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
	return f(raw)
}

func Run(loader AppLoader, validator AppValidator, action CommandAction) error {
	if loader == nil {
		return errors.New("loader is required")
	}
	if validator == nil {
		return errors.New("validator is required")
	}
	if action == nil {
		return errors.New("action is required")
	}

	raw, err := loader.Load()
	if err != nil {
		return err
	}

	validated, report, err := validator.Validate(raw)
	if err != nil {
		return err
	}

	if report.HasErrors() && !action.AllowOnValidationErrors() {
		return outcome.ValidationBlockedError()
	}

	return action.Execute(validated, report)
}
