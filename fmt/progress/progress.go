// Package progress provides terminal progress bar rendering with stage support.
// It handles ANSI escape codes for line clearing and overwriting to show
// updating progress without scrolling the terminal.
package progress

import (
	"fmt"
	"io"
	"strings"
)

// StageInfo contains information about the current progress state.
type StageInfo struct {
	Stage       int    // Current stage number (1-based)
	TotalStages int    // Total number of stages
	Description string // What's happening in this stage
	Current     int    // Current item within stage (0 if not applicable)
	Total       int    // Total items in stage (0 if not applicable)
	Done        bool   // True if this stage is complete
	Text        string // Optional extra text (e.g., current item name)
}

// Percent returns the completion percentage for this stage.
// Returns 100 if Done is true, 0 if Total is 0, otherwise Current/Total * 100.
func (s StageInfo) Percent() int {
	if s.Done {
		return 100
	}
	if s.Total <= 0 || s.Current <= 0 {
		return 0
	}
	return (s.Current * 100) / s.Total
}

// MultiStageRenderer handles displaying progress to a terminal with stage support.
// It manages line clearing and overwriting for smooth progress updates.
type MultiStageRenderer struct {
	w         io.Writer
	barWidth  int
	descWidth int
	lastStage int
	lastDone  bool
}

// NewMultiStageRenderer creates a new multi-stage progress renderer that writes to w.
// Default bar width is 20 characters, description width is 34 characters.
func NewMultiStageRenderer(w io.Writer) *MultiStageRenderer {
	return &MultiStageRenderer{
		w:         w,
		barWidth:  20,
		descWidth: 34,
	}
}

// WithBarWidth sets the width of the progress bar in characters.
func (r *MultiStageRenderer) WithBarWidth(width int) *MultiStageRenderer {
	r.barWidth = width
	return r
}

// WithDescWidth sets the width of the description field in characters.
func (r *MultiStageRenderer) WithDescWidth(width int) *MultiStageRenderer {
	r.descWidth = width
	return r
}

// Update displays the current progress state.
// When Done is false, the line is overwritten on subsequent calls.
// When Done is true, a newline is printed and the next update starts on a new line.
func (r *MultiStageRenderer) Update(info StageInfo) {
	// Skip duplicate "done" messages for the same stage
	if info.Done && r.lastStage == info.Stage && r.lastDone {
		return
	}

	// Clear line and move cursor to beginning
	clearLine := "\r\033[K"

	pct := info.Percent()
	bar := RenderBar(pct, r.barWidth)

	// Format string with proper field widths
	format := fmt.Sprintf("%%s[%%d/%%d] %%-%ds %%s %%3d%%%%", r.descWidth)

	if info.Done {
		// Print completed stage with newline
		fmt.Fprintf(r.w, format+"\n",
			clearLine, info.Stage, info.TotalStages, info.Description, bar, pct)
		r.lastDone = true
	} else {
		// Print in-progress stage without newline (will be overwritten)
		fmt.Fprintf(r.w, format,
			clearLine, info.Stage, info.TotalStages, info.Description, bar, pct)
		// Add item count after percentage if available
		if info.Current > 0 && info.Total > 0 {
			fmt.Fprintf(r.w, " (%d/%d)", info.Current, info.Total)
		}
		// Add extra text if provided
		if info.Text != "" {
			fmt.Fprintf(r.w, " %s", info.Text)
		}
		r.lastDone = false
	}
	r.lastStage = info.Stage
}

// RenderBar creates a visual progress bar string.
// width is the number of characters for the bar (excluding brackets).
// Returns a string like "[████████░░░░░░░░░░░░]"
func RenderBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := (percent * width) / 100
	empty := width - filled

	// Use Unicode block characters for a smooth look
	// █ (full block) for filled, ░ (light shade) for empty
	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return "[" + bar + "]"
}

// RenderBarASCII creates a progress bar using ASCII characters only.
// Returns a string like "[========------------]"
func RenderBarASCII(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := (percent * width) / 100
	empty := width - filled

	bar := strings.Repeat("=", filled) + strings.Repeat("-", empty)
	return "[" + bar + "]"
}

// SingleStageRenderer handles displaying single-stage progress to a terminal.
// Use this for progress tracking without multiple stages.
type SingleStageRenderer struct {
	w         io.Writer
	barWidth  int
	textWidth int
}

// NewSingleStageRenderer creates a new single-stage progress renderer that writes to w.
// Default bar width is 40 characters, text width is 30 characters.
func NewSingleStageRenderer(w io.Writer) *SingleStageRenderer {
	return &SingleStageRenderer{
		w:         w,
		barWidth:  40,
		textWidth: 30,
	}
}

// WithBarWidth sets the width of the progress bar in characters.
func (r *SingleStageRenderer) WithBarWidth(width int) *SingleStageRenderer {
	r.barWidth = width
	return r
}

// WithTextWidth sets the width of the text field in characters.
func (r *SingleStageRenderer) WithTextWidth(width int) *SingleStageRenderer {
	r.textWidth = width
	return r
}

// Update displays the current progress state.
// The line is overwritten on subsequent calls.
// Call Done() when complete to print a newline.
func (r *SingleStageRenderer) Update(current, total int, text string) {
	// Clear line and move cursor to beginning
	clearLine := "\r\033[K"

	percent := 0
	if total > 0 {
		percent = (current * 100) / total
	}
	bar := RenderBar(percent, r.barWidth)

	// Truncate text if too long
	displayText := text
	if len(displayText) > r.textWidth {
		displayText = displayText[:r.textWidth-3] + "..."
	}

	// Format: [████░░░░] 20% (5/42) some-text
	format := fmt.Sprintf("%%s%%s %%3d%%%% (%%d/%%d) %%-%ds", r.textWidth)
	fmt.Fprintf(r.w, format, clearLine, bar, percent, current, total, displayText)
}

// Done clears the progress line and optionally prints a completion message.
// If message is empty, just clears the line.
func (r *SingleStageRenderer) Done(message string) {
	// Clear the line
	fmt.Fprintf(r.w, "\r\033[K")
	if message != "" {
		fmt.Fprintln(r.w, message)
	}
}
