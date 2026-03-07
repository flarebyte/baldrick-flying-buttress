package pipeline

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/app"
	"github.com/flarebyte/baldrick-flying-buttress/internal/diagnostics"
)

func TestRunCallsLoadValidateAction(t *testing.T) {
	t.Parallel()

	var calls []string
	wantRaw := app.RawApp{Source: "raw-stub"}
	wantValidated := app.ValidatedApp{Name: "validated", Modules: []string{"core"}}
	wantReport := diagnostics.Report{Diagnostics: []diagnostics.Diagnostic{{
		Code:     "FBW01",
		Severity: diagnostics.SeverityWarning,
		Message:  "warning",
		Path:     "module.stub",
	}}}

	err := Run(
		func() (app.RawApp, error) {
			calls = append(calls, "load")
			return wantRaw, nil
		},
		func(raw app.RawApp) (app.ValidatedApp, diagnostics.Report, error) {
			calls = append(calls, "validate")
			if raw != wantRaw {
				t.Fatalf("unexpected raw app: %#v", raw)
			}
			return wantValidated, wantReport, nil
		},
		func(validated app.ValidatedApp, report diagnostics.Report) error {
			calls = append(calls, "action")
			if validated.Name != wantValidated.Name {
				t.Fatalf("unexpected app name: %s", validated.Name)
			}
			if len(validated.Modules) != len(wantValidated.Modules) || validated.Modules[0] != wantValidated.Modules[0] {
				t.Fatalf("unexpected modules: %#v", validated.Modules)
			}
			if len(report.Diagnostics) != len(wantReport.Diagnostics) || report.Diagnostics[0] != wantReport.Diagnostics[0] {
				t.Fatalf("unexpected report: %#v", report)
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	if len(calls) != 3 || calls[0] != "load" || calls[1] != "validate" || calls[2] != "action" {
		t.Fatalf("unexpected call order: %#v", calls)
	}
}
