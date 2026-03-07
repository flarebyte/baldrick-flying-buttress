package validate

import (
	"github.com/olivier/baldrick-flying-buttress/internal/app"
	"github.com/olivier/baldrick-flying-buttress/internal/diagnostics"
)

func LoadStub() (app.RawApp, error) {
	return app.RawApp{
		Source: "in-memory-stub",
	}, nil
}

func ValidateStub(raw app.RawApp) (app.ValidatedApp, diagnostics.Report, error) {
	_ = raw
	return app.ValidatedApp{
			Name:    "stub-app",
			Modules: []string{"core", "edge"},
			Reports: []app.Report{
				{
					ID:    "cpu-overview",
					Title: "CPU Overview",
				},
				{
					ID:    "memory-health",
					Title: "Memory Health",
				},
			},
		}, diagnostics.Report{
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
