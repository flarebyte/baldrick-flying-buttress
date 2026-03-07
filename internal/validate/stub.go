package validate

import "github.com/flarebyte/baldrick-flying-buttress/internal/domain"

type StubAppLoader struct{}

type StubAppValidator struct{}

func (StubAppLoader) Load() (domain.RawApp, error) {
	return domain.RawApp{
		Source: "in-memory-stub",
	}, nil
}

func (StubAppValidator) Validate(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
	_ = raw
	return domain.ValidatedApp{
			Name:    "stub-app",
			Modules: []string{"core", "edge"},
			Reports: []domain.Report{
				{
					ID:    "cpu-overview",
					Title: "CPU Overview",
				},
				{
					ID:    "memory-health",
					Title: "Memory Health",
				},
			},
			Notes: []domain.Note{
				{
					ID:    "n1",
					Label: "service.api",
				},
				{
					ID:    "n2",
					Label: "service.db",
				},
			},
			Relationships: []domain.Relationship{
				{
					FromID: "n1",
					ToID:   "n2",
					Label:  "depends_on",
				},
			},
		}, domain.ValidationReport{
			Diagnostics: []domain.Diagnostic{
				{
					Code:     "FB001",
					Severity: domain.SeverityWarning,
					Message:  "stub warning diagnostic",
					Path:     "module.stub",
				},
				{
					Code:     "FB002",
					Severity: domain.SeverityError,
					Message:  "stub error diagnostic",
					Path:     "module.stub",
				},
			},
		}, nil
}

func LoadStub() (domain.RawApp, error) {
	return StubAppLoader{}.Load()
}

func ValidateStub(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
	return StubAppValidator{}.Validate(raw)
}
