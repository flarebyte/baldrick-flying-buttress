package cli

import (
	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func withDiagnosticContextMessage(d domain.Diagnostic) domain.Diagnostic {
	d.Message = domain.FormatDiagnosticMessage(d.Message, d)
	return d
}
