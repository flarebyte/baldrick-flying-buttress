package validate

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestValidateStubReturnsFullyPopulatedValidatedApp(t *testing.T) {
	t.Parallel()

	app, report, err := ValidateStub(domain.RawApp{Source: "in-memory-stub"})
	if err != nil {
		t.Fatalf("validate stub failed: %v", err)
	}

	if app.Name == "" {
		t.Fatal("expected validated app name")
	}
	if len(app.Modules) == 0 {
		t.Fatal("expected modules")
	}
	if len(app.Reports) == 0 {
		t.Fatal("expected reports")
	}
	if len(app.Notes) == 0 {
		t.Fatal("expected notes")
	}
	if len(app.Relationships) == 0 {
		t.Fatal("expected relationships")
	}
	if len(report.Diagnostics) == 0 {
		t.Fatal("expected diagnostics")
	}
}
