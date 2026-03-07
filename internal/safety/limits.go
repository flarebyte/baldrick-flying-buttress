package safety

import (
	"errors"
	"fmt"
)

const (
	MaxConfigFileBytes            int64 = 1 << 20
	MaxReportsCount                     = 2000
	MaxNotesCount                       = 10000
	MaxRelationshipsCount               = 20000
	MaxCSVFileBytes                     = 1 << 20
	MaxCSVRowsRenderedPerNote           = 5000
	MaxGraphRenderNodesPerSection       = 10000
)

func CheckConfigFileSize(size int64) error {
	if size <= MaxConfigFileBytes {
		return nil
	}
	return newLimitError("config file too large: %d bytes exceeds limit %d bytes", size, MaxConfigFileBytes)
}

func CheckReportsCount(count int) error {
	if count <= MaxReportsCount {
		return nil
	}
	return newLimitError("reports count %d exceeds limit %d", count, MaxReportsCount)
}

func CheckNotesCount(count int) error {
	if count <= MaxNotesCount {
		return nil
	}
	return newLimitError("notes count %d exceeds limit %d", count, MaxNotesCount)
}

func CheckRelationshipsCount(count int) error {
	if count <= MaxRelationshipsCount {
		return nil
	}
	return newLimitError("relationships count %d exceeds limit %d", count, MaxRelationshipsCount)
}

func CheckCSVFileSize(size int) error {
	if size <= MaxCSVFileBytes {
		return nil
	}
	return newLimitError("csv file too large: %d bytes exceeds limit %d bytes", size, MaxCSVFileBytes)
}

func CheckCSVRenderedRows(rows int) error {
	if rows <= MaxCSVRowsRenderedPerNote {
		return nil
	}
	return newLimitError("csv rendered rows %d exceeds limit %d", rows, MaxCSVRowsRenderedPerNote)
}

func CheckGraphRenderNodeCount(nodes int) error {
	if nodes <= MaxGraphRenderNodesPerSection {
		return nil
	}
	return newLimitError("graph render node count %d exceeds limit %d", nodes, MaxGraphRenderNodesPerSection)
}

type LimitError struct {
	message string
}

func (e LimitError) Error() string {
	return e.message
}

func IsLimitError(err error) bool {
	var limitErr LimitError
	return errors.As(err, &limitErr)
}

func newLimitError(format string, args ...any) error {
	return LimitError{message: fmt.Sprintf(format, args...)}
}
