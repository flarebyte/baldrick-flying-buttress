package validate

import (
	"context"

	"github.com/olivier/baldrick-flying-buttress/internal/diagnostics"
)

type Runner func(ctx context.Context) (diagnostics.Report, error)

func RunStub(ctx context.Context) (diagnostics.Report, error) {
	_ = ctx
	return diagnostics.Report{
		Diagnostics: []diagnostics.Diagnostic{
			{
				Code:     "FB001",
				Severity: diagnostics.SeverityWarning,
				Message:  "stub warning diagnostic",
				Path:     "module.stub",
			},
			{
				Code:     "FB002",
				Severity: diagnostics.SeverityError,
				Message:  "stub error diagnostic",
				Path:     "module.stub",
			},
		},
	}, nil
}
