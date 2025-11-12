package renderer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Rectangle represents a rendered box with content
// It stores the rendered output as lines of text
type Rectangle struct {
	Width  int      // Visual width (excluding ANSI codes)
	Height int      // Number of lines
	Lines  []string // Each line of rendered content (may contain ANSI codes)
}

// NewRectangle creates a Rectangle from a rendered string
func NewRectangle(content string) Rectangle {
	// Remove trailing newline if present
	content = strings.TrimSuffix(content, "\n")

	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return Rectangle{Width: 0, Height: 0, Lines: []string{}}
	}

	// Calculate width using lipgloss.Width to handle ANSI codes
	maxWidth := 0
	for _, line := range lines {
		width := lipgloss.Width(line)
		if width > maxWidth {
			maxWidth = width
		}
	}

	return Rectangle{
		Width:  maxWidth,
		Height: len(lines),
		Lines:  lines,
	}
}

// NewEmptyRectangle creates an empty Rectangle with the given dimensions
func NewEmptyRectangle(width, height int) Rectangle {
	lines := make([]string, height)
	for i := 0; i < height; i++ {
		lines[i] = strings.Repeat(" ", width)
	}

	return Rectangle{
		Width:  width,
		Height: height,
		Lines:  lines,
	}
}

// Overlay merges child Rectangle into parent Rectangle
// The child is overlaid on top of the parent at position (0, 0)
// Non-space characters from child replace characters in parent (proper shadowing)
// Returns a new Rectangle with the merged content
func Overlay(parent, child Rectangle) Rectangle {
	// Determine the dimensions of the result
	resultWidth := parent.Width
	if child.Width > resultWidth {
		resultWidth = child.Width
	}

	resultHeight := parent.Height
	if child.Height > resultHeight {
		resultHeight = child.Height
	}

	// Create result lines
	resultLines := make([]string, resultHeight)

	// Copy parent lines and extend if necessary
	for i := 0; i < resultHeight; i++ {
		if i < len(parent.Lines) {
			resultLines[i] = parent.Lines[i]
		} else {
			resultLines[i] = strings.Repeat(" ", resultWidth)
		}

		// Ensure line is at least resultWidth wide
		lineWidth := lipgloss.Width(resultLines[i])
		if lineWidth < resultWidth {
			resultLines[i] += strings.Repeat(" ", resultWidth-lineWidth)
		}
	}

	// Overlay child lines onto result
	for i := 0; i < len(child.Lines) && i < resultHeight; i++ {
		// Pad child line to resultWidth before overlaying
		childLine := child.Lines[i]
		childLineWidth := lipgloss.Width(childLine)
		if childLineWidth < resultWidth {
			childLine += strings.Repeat(" ", resultWidth-childLineWidth)
		}
		// Pass the actual rectangle widths to determine shadowing area
		resultLines[i] = overlayLine(resultLines[i], childLine, resultWidth, child.Width)
	}

	return Rectangle{
		Width:  resultWidth,
		Height: resultHeight,
		Lines:  resultLines,
	}
}

// overlayLine overlays the child line onto the parent line
// The child shadows the parent only within the child's rectangle width
// parentRectWidth: the width of the parent rectangle (for result sizing)
// childRectWidth: the actual width of the child rectangle (determines shadowing area)
func overlayLine(parent, child string, parentRectWidth, childRectWidth int) string {
	// If child rectangle has no width, return parent
	if childRectWidth == 0 {
		return parent
	}

	// Strip ANSI from both to work with visual characters
	parentStripped := stripANSI(parent)
	childStripped := stripANSI(child)

	// Convert to rune slices for proper character handling
	parentRunes := []rune(parentStripped)
	childRunes := []rune(childStripped)

	// Build result character by character
	// The child shadows the parent only within childRectWidth
	maxLen := len(parentRunes)
	if len(childRunes) > maxLen {
		maxLen = len(childRunes)
	}

	result := make([]rune, maxLen)
	for i := 0; i < maxLen; i++ {
		// If within the child's rectangle width, use child (shadowing)
		// Otherwise, use parent (visible through)
		if i < childRectWidth && i < len(childRunes) {
			result[i] = childRunes[i]
		} else if i < len(parentRunes) {
			result[i] = parentRunes[i]
		} else {
			result[i] = ' '
		}
	}

	// For now, return without ANSI codes to avoid corruption
	// TODO: Properly preserve ANSI codes in future iteration
	return string(result)
}

// String returns the Rectangle as a string (with newlines between lines)
// Renders exactly Height lines, each exactly Width characters wide
// - Lines longer than Width are truncated (ANSI codes preserved)
// - Lines shorter than Width are padded with spaces
// - Extra lines beyond Height are hidden
// - Missing lines are filled with spaces
func (r Rectangle) String() string {
	result := make([]string, r.Height)

	for i := 0; i < r.Height; i++ {
		var line string

		if i < len(r.Lines) {
			line = r.Lines[i]
		} else {
			// Missing line - fill with spaces
			line = ""
		}

		// Calculate the visual width of the line (handling ANSI codes)
		lineWidth := lipgloss.Width(line)

		if lineWidth > r.Width {
			// Line is too long - truncate it (preserving ANSI codes)
			line = lipgloss.NewStyle().Width(r.Width).Render(line)
		} else if lineWidth < r.Width {
			// Line is too short - pad with spaces
			line += strings.Repeat(" ", r.Width-lineWidth)
		}
		// else: line is exactly the right width

		result[i] = line
	}

	return strings.Join(result, "\n")
}
