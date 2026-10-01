package markdown

import (
	"regexp"
	"strings"
)

// headerPattern matches markdown headers (# to ######).
var headerPattern = regexp.MustCompile(`^(#{1,6})\s+(.*)$`)

// ShiftHeaders adjusts all markdown headers by the specified offset.
// offset=1 converts # to ##, ## to ###, etc.
// offset=-1 converts ## to #, ### to ##, etc.
// Headers exceeding H6 are capped at H6.
// Headers that would go below H1 are capped at H1.
// Lines inside fenced code blocks (``` or ~~~) are not modified.
func ShiftHeaders(content string, offset int) string {
	if offset == 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	inCodeBlock := false

	for i, line := range lines {
		// Check for fenced code block toggle
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inCodeBlock = !inCodeBlock
			continue
		}

		// Skip lines inside code blocks
		if inCodeBlock {
			continue
		}

		// Check if this line is a header
		matches := headerPattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		currentLevel := len(matches[1])
		newLevel := currentLevel + offset

		// Cap at H1 minimum, H6 maximum
		if newLevel < 1 {
			newLevel = 1
		}
		if newLevel > 6 {
			newLevel = 6
		}

		// Rebuild the header line
		lines[i] = strings.Repeat("#", newLevel) + " " + matches[2]
	}

	return strings.Join(lines, "\n")
}

// NormalizeHeaders shifts headers so the highest level becomes targetLevel.
// For example, if content has # and ##, and targetLevel=3, they become ### and ####.
// If the content has no headers, it is returned unchanged.
// Lines inside fenced code blocks are not considered or modified.
func NormalizeHeaders(content string, targetLevel int) string {
	// Clamp targetLevel to valid range
	if targetLevel < 1 {
		targetLevel = 1
	}
	if targetLevel > 6 {
		targetLevel = 6
	}

	// Find the minimum header level in the content
	minLevel := findMinHeaderLevel(content)
	if minLevel == 0 {
		// No headers found
		return content
	}

	// Calculate offset needed
	offset := targetLevel - minLevel
	return ShiftHeaders(content, offset)
}

// findMinHeaderLevel returns the minimum header level (1-6) found in the content,
// or 0 if no headers are found. Lines inside fenced code blocks are skipped.
func findMinHeaderLevel(content string) int {
	lines := strings.Split(content, "\n")
	inCodeBlock := false
	minLevel := 0

	for _, line := range lines {
		// Check for fenced code block toggle
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inCodeBlock = !inCodeBlock
			continue
		}

		// Skip lines inside code blocks
		if inCodeBlock {
			continue
		}

		// Check if this line is a header
		matches := headerPattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		level := len(matches[1])
		if minLevel == 0 || level < minLevel {
			minLevel = level
		}
	}

	return minLevel
}
