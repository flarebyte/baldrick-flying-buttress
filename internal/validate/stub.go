package validate

import "github.com/flarebyte/baldrick-flying-buttress/internal/domain"

type StubAppLoader struct{}

type StubAppValidator struct{}

func (StubAppLoader) Load() (domain.RawApp, error) {
	return domain.RawApp{
		ConfigPath: "in-memory-stub",
		Source:     "in-memory-stub",
		Name:       "stub-app",
		Modules:    []string{"core", "edge"},
		Reports: []domain.RawReport{
			{ID: "cpu-overview", Title: "CPU Overview"},
			{ID: "memory-health", Title: "Memory Health"},
		},
		Notes: []domain.RawNote{
			{ID: "n1", Label: "service.api"},
			{ID: "n2", Label: "service.db"},
		},
		Relationships: []domain.RawRelationship{
			{FromID: "n1", ToID: "n2", Label: "depends_on"},
		},
	}, nil
}

func (StubAppValidator) Validate(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
	return AppDataValidator{}.Validate(raw)
}

func LoadStub() (domain.RawApp, error) {
	return StubAppLoader{}.Load()
}

func ValidateStub(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
	return StubAppValidator{}.Validate(raw)
}
