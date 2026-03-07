package validate

import (
	"testing"
)

func TestValidateStubReturnsFullyPopulatedValidatedApp(t *testing.T) {
	t.Parallel()

	raw, err := LoadStub()
	if err != nil {
		t.Fatalf("load stub failed: %v", err)
	}

	app, report, err := ValidateStub(raw)
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
	if report.Diagnostics == nil {
		t.Fatal("expected non-nil diagnostics slice")
	}
}
