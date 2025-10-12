package markdown

import (
	"strings"
)

// TableAlign takes a markdown table string and aligns the pipes.
func TableAlign(input string, sepLineRowIdx int) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return ""
	}

	// Split rows into columns
	tbl := make([][]string, len(lines))
	colWidths := []int{}

	for i, line := range lines {
		// Remove leading/trailing pipe and split
		line = strings.Trim(line, "|")
		cols := strings.Split(line, "|")
		for j := range cols {
			cols[j] = strings.TrimSpace(cols[j])
		}
		tbl[i] = cols

		// Track max width per column
		for j, col := range cols {
			if len(colWidths) <= j {
				colWidths = append(colWidths, len(col))
			} else if len(col) > colWidths[j] {
				colWidths[j] = len(col)
			}
		}
	}

	// Rebuild table with aligned columns
	var sb strings.Builder
	for i, row := range tbl {
		sb.WriteString("|")
		paddingChar := " "
		if i == sepLineRowIdx {
			paddingChar = "-"
		}
		for j, cellVal := range row {
			paddingLen := colWidths[j] - len(cellVal)
			sb.WriteString(" " + cellVal + strings.Repeat(paddingChar, paddingLen) + " |")
		}
		// add newline except for last line
		if i < len(tbl)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
