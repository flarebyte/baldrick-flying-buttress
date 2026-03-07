package cli

import "strings"

func formatAlignedTableRow(cells []string, widths []int) string {
	var b strings.Builder
	for i := range widths {
		if i > 0 {
			b.WriteString(" | ")
		}
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}
		b.WriteString(cell)
		padding := widths[i] - len(cell)
		if padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
		}
	}
	b.WriteByte('\n')
	return b.String()
}
