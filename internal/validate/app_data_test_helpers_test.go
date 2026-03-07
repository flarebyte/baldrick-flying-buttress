package validate

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

type diagExpectation struct {
	code string
	path string
}

func assertHasDiagnostics(t *testing.T, diagnostics []domain.Diagnostic, want []diagExpectation, label string) {
	t.Helper()
	for _, item := range want {
		found := false
		for _, d := range diagnostics {
			if d.Code == item.code && d.Path == item.path {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected %s diagnostic %s at %s, got %#v", label, item.code, item.path, diagnostics)
		}
	}
}

func validateRaw(t *testing.T, raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport) {
	t.Helper()
	app, report, err := AppDataValidator{}.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	return app, report
}
