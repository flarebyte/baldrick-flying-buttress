// purpose: Implements CLI behavior for severity.go so commands expose deterministic, machine-friendly output surfaces.
// responsibilities: parse command inputs; call pipeline/domain services; render structured outputs or diagnostics; enforce deterministic CLI behavior
// architecture_notes: CLI logic is split into focused files per command area to keep Cobra wiring thin and to isolate rendering from validation and domain logic.
package cli

import (
	"fmt"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func resolveSeverity(severity string) (domain.Severity, error) {
	switch strings.TrimSpace(severity) {
	case "", string(domain.SeverityWarning):
		return domain.SeverityWarning, nil
	case string(domain.SeverityError):
		return domain.SeverityError, nil
	default:
		return "", fmt.Errorf("invalid severity: %s", severity)
	}
}
