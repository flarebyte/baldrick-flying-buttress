// purpose: Implements CLI behavior for diagnostics_message.go so commands expose deterministic, machine-friendly output surfaces.
// responsibilities: parse command inputs; call pipeline/domain services; render structured outputs or diagnostics; enforce deterministic CLI behavior
// architecture_notes: CLI logic is split into focused files per command area to keep Cobra wiring thin and to isolate rendering from validation and domain logic.
package cli

import (
	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func withDiagnosticContextMessage(d domain.Diagnostic) domain.Diagnostic {
	d.Message = domain.FormatDiagnosticMessage(d.Message, d)
	return d
}
