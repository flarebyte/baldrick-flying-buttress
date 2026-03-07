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
