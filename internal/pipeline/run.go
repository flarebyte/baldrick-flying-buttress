package pipeline

import (
	"context"
	"errors"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
)

type AppLoader interface {
	Load(context.Context) (domain.RawApp, error)
}

type AppValidator interface {
	Validate(context.Context, domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error)
}

type CommandAction interface {
	Execute(context.Context, domain.ValidatedApp, domain.ValidationReport) error
	AllowOnValidationErrors() bool
}

type LoaderFunc func(context.Context) (domain.RawApp, error)

func (f LoaderFunc) Load(ctx context.Context) (domain.RawApp, error) {
	return f(ctx)
}

type ValidatorFunc func(context.Context, domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error)

func (f ValidatorFunc) Validate(ctx context.Context, raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
	return f(ctx, raw)
}

func Run(ctx context.Context, loader AppLoader, validator AppValidator, action CommandAction) error {
	if loader == nil {
		return errors.New("loader is required")
	}
	if validator == nil {
		return errors.New("validator is required")
	}
	if action == nil {
		return errors.New("action is required")
	}
	if ctx == nil {
		return errors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	raw, err := loader.Load(ctx)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	validated, report, err := validator.Validate(ctx, raw)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	if report.HasErrors() && !action.AllowOnValidationErrors() {
		return outcome.ValidationBlockedError()
	}

	return action.Execute(ctx, validated, report)
}
