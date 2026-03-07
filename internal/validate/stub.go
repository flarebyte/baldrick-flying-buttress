package validate

import (
	"context"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

type StubAppLoader struct{}

type StubAppValidator struct{}

func (StubAppLoader) Load(context.Context) (domain.RawApp, error) {
	return domain.RawApp{
		Source:  "in-memory-stub",
		Name:    "stub-app",
		Modules: []string{"core", "edge"},
		Reports: []domain.RawReport{
			{Title: "CPU Overview", Filepath: "reports/cpu-overview.md", Sections: []domain.RawReportSection{{Title: "Overview"}}},
			{Title: "Memory Health", Filepath: "reports/memory-health.md", Sections: []domain.RawReportSection{{Title: "Overview"}}},
		},
		Notes: []domain.RawNote{
			{Name: "n1", Title: "service.api"},
			{Name: "n2", Title: "service.db"},
		},
		Relationships: []domain.RawRelationship{
			{FromID: "n1", ToID: "n2", Label: "depends_on"},
		},
	}, nil
}

func (StubAppValidator) Validate(ctx context.Context, raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
	return AppDataValidator{}.Validate(ctx, raw)
}

func LoadStub(ctx context.Context) (domain.RawApp, error) {
	return StubAppLoader{}.Load(ctx)
}

func ValidateStub(ctx context.Context, raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
	return StubAppValidator{}.Validate(ctx, raw)
}
