package cli

import (
	"io"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	clioutput "github.com/flarebyte/baldrick-flying-buttress/internal/output"
)

func emitDiagnosticsOutcome(out io.Writer, diagnostics []domain.Diagnostic) error {
	lintReport := domain.ValidationReport{Diagnostics: diagnostics}
	if err := clioutput.EmitDiagnostics(out, lintReport); err != nil {
		return err
	}
	if lintReport.HasErrors() {
		return outcome.ValidationBlockedError()
	}
	return nil
}
