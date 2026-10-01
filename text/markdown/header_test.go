package markdown

import (
	"testing"
)

func TestShiftHeaders(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		offset   int
		expected string
	}{
		{
			name:     "shift up by 1",
			content:  "# Heading 1\n\nSome text\n\n## Heading 2",
			offset:   1,
			expected: "## Heading 1\n\nSome text\n\n### Heading 2",
		},
		{
			name:     "shift up by 2",
			content:  "# Heading 1\n## Heading 2\n### Heading 3",
			offset:   2,
			expected: "### Heading 1\n#### Heading 2\n##### Heading 3",
		},
		{
			name:     "shift down by 1",
			content:  "## Heading 2\n### Heading 3",
			offset:   -1,
			expected: "# Heading 2\n## Heading 3",
		},
		{
			name:     "cap at H6",
			content:  "##### Heading 5\n###### Heading 6",
			offset:   2,
			expected: "###### Heading 5\n###### Heading 6",
		},
		{
			name:     "cap at H1",
			content:  "## Heading 2\n# Heading 1",
			offset:   -5,
			expected: "# Heading 2\n# Heading 1",
		},
		{
			name:     "zero offset",
			content:  "# Heading 1\n## Heading 2",
			offset:   0,
			expected: "# Heading 1\n## Heading 2",
		},
		{
			name:     "skip code blocks",
			content:  "# Real Header\n\n```\n# Not a header\n## Also not\n```\n\n## Another Header",
			offset:   1,
			expected: "## Real Header\n\n```\n# Not a header\n## Also not\n```\n\n### Another Header",
		},
		{
			name:     "skip tilde code blocks",
			content:  "# Real\n~~~\n# Fake\n~~~\n## Real2",
			offset:   1,
			expected: "## Real\n~~~\n# Fake\n~~~\n### Real2",
		},
		{
			name:     "no headers",
			content:  "Just some text\n\nMore text",
			offset:   1,
			expected: "Just some text\n\nMore text",
		},
		{
			name:     "preserve header content",
			content:  "# Title with `code` and **bold**",
			offset:   1,
			expected: "## Title with `code` and **bold**",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShiftHeaders(tt.content, tt.offset)
			if result != tt.expected {
				t.Errorf("ShiftHeaders(%q, %d) = %q, want %q", tt.content, tt.offset, result, tt.expected)
			}
		})
	}
}

func TestNormalizeHeaders(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		targetLevel int
		expected    string
	}{
		{
			name:        "normalize H1 to H2",
			content:     "# Heading 1\n## Heading 2\n### Heading 3",
			targetLevel: 2,
			expected:    "## Heading 1\n### Heading 2\n#### Heading 3",
		},
		{
			name:        "normalize H1 to H3",
			content:     "# Heading 1\n## Heading 2",
			targetLevel: 3,
			expected:    "### Heading 1\n#### Heading 2",
		},
		{
			name:        "normalize H2 to H3 (min is H2)",
			content:     "## Heading 2\n### Heading 3\n#### Heading 4",
			targetLevel: 3,
			expected:    "### Heading 2\n#### Heading 3\n##### Heading 4",
		},
		{
			name:        "normalize down",
			content:     "### Heading 3\n#### Heading 4",
			targetLevel: 1,
			expected:    "# Heading 3\n## Heading 4",
		},
		{
			name:        "no headers",
			content:     "Just text\nMore text",
			targetLevel: 3,
			expected:    "Just text\nMore text",
		},
		{
			name:        "already at target level",
			content:     "## Heading 2\n### Heading 3",
			targetLevel: 2,
			expected:    "## Heading 2\n### Heading 3",
		},
		{
			name:        "cap at H6",
			content:     "# H1\n## H2\n### H3\n#### H4\n##### H5",
			targetLevel: 4,
			expected:    "#### H1\n##### H2\n###### H3\n###### H4\n###### H5",
		},
		{
			name:        "clamp target level below 1",
			content:     "## Heading 2",
			targetLevel: 0,
			expected:    "# Heading 2",
		},
		{
			name:        "clamp target level above 6",
			content:     "# Heading 1",
			targetLevel: 10,
			expected:    "###### Heading 1",
		},
		{
			name:        "skip code blocks when finding min",
			content:     "```\n# Comment\n```\n## Real Heading\n### Another",
			targetLevel: 3,
			expected:    "```\n# Comment\n```\n### Real Heading\n#### Another",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeHeaders(tt.content, tt.targetLevel)
			if result != tt.expected {
				t.Errorf("NormalizeHeaders(%q, %d) = %q, want %q", tt.content, tt.targetLevel, result, tt.expected)
			}
		})
	}
}

func TestFindMinHeaderLevel(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "H1 is minimum",
			content:  "# H1\n## H2\n### H3",
			expected: 1,
		},
		{
			name:     "H2 is minimum",
			content:  "## H2\n### H3\n#### H4",
			expected: 2,
		},
		{
			name:     "H3 only",
			content:  "### Only H3",
			expected: 3,
		},
		{
			name:     "no headers",
			content:  "Just text",
			expected: 0,
		},
		{
			name:     "skip code blocks",
			content:  "```\n# In code\n```\n## Real",
			expected: 2,
		},
		{
			name:     "empty content",
			content:  "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findMinHeaderLevel(tt.content)
			if result != tt.expected {
				t.Errorf("findMinHeaderLevel(%q) = %d, want %d", tt.content, result, tt.expected)
			}
		})
	}
}
