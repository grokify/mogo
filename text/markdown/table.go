package markdown

import (
	"strings"
)

// TableAlign takes a markdown table string and aligns the pipes.
func TableAlign(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return ""
	}

	// Split rows into columns
	table := make([][]string, len(lines))
	colWidths := []int{}

	for i, line := range lines {
		// Remove leading/trailing pipe and split
		line = strings.Trim(line, "|")
		cols := strings.Split(line, "|")
		for j := range cols {
			cols[j] = strings.TrimSpace(cols[j])
		}
		table[i] = cols

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
	for _, row := range table {
		sb.WriteString("|")
		for j, col := range row {
			padding := colWidths[j] - len(col)
			sb.WriteString(" " + col + strings.Repeat(" ", padding) + " |")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
