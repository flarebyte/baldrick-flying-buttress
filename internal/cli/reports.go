package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
	clioutput "github.com/flarebyte/baldrick-flying-buttress/internal/output"
)

const (
	reportsFormatJSON  = "json"
	reportsFormatTable = "table"
)

func validateReportsFormat(format string) error {
	switch format {
	case reportsFormatJSON, reportsFormatTable:
		return nil
	default:
		return fmt.Errorf("invalid format: %s", format)
	}
}

type listReportsAction struct {
	out    io.Writer
	format string
}

func (a listReportsAction) Execute(_ context.Context, validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = report
	switch a.format {
	case reportsFormatTable:
		return emitReportsTable(a.out, validated.Reports)
	case reportsFormatJSON:
		return clioutput.EmitReportList(a.out, validated)
	default:
		return fmt.Errorf("invalid format: %s", a.format)
	}
}

func (listReportsAction) AllowOnValidationErrors() bool {
	return false
}

func emitReportsTable(w io.Writer, reports []domain.Report) error {
	ordered := ordering.Reports(reports)
	headers := []string{"title", "filepath"}
	widths := []int{len(headers[0]), len(headers[1])}
	rows := make([][]string, 0, len(ordered))
	for _, report := range ordered {
		row := []string{report.Title, report.ID}
		if len(row[0]) > widths[0] {
			widths[0] = len(row[0])
		}
		if len(row[1]) > widths[1] {
			widths[1] = len(row[1])
		}
		rows = append(rows, row)
	}
	if _, err := io.WriteString(w, formatAlignedTableRow(headers, widths)); err != nil {
		return err
	}
	divider := []string{strings.Repeat("-", widths[0]), strings.Repeat("-", widths[1])}
	if _, err := io.WriteString(w, formatAlignedTableRow(divider, widths)); err != nil {
		return err
	}
	for _, row := range rows {
		if _, err := io.WriteString(w, formatAlignedTableRow(row, widths)); err != nil {
			return err
		}
	}
	return nil
}
