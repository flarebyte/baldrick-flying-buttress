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
