// purpose: Implements CLI behavior for table.go so commands expose deterministic, machine-friendly output surfaces.
// responsibilities: parse command inputs; call pipeline/domain services; render structured outputs or diagnostics; enforce deterministic CLI behavior
// architecture_notes: CLI logic is split into focused files per command area to keep Cobra wiring thin and to isolate rendering from validation and domain logic.
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
