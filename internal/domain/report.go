// purpose: Declares core domain models in report.go used across validation, pipeline, and CLI layers.
// responsibilities: define canonical data structures; hold stable typed contracts; provide shared semantic primitives
// architecture_notes: Domain types are dependency-light and reused broadly to keep cross-layer contracts explicit and avoid cyclic package coupling.
package domain

import (
	"path/filepath"
	"strings"
)

func ReportIDFromFilepath(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
