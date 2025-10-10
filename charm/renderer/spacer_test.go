package renderer

import (
	"strings"
	"testing"

	"github.com/xhd2015/go-dom-tui/dom"
)

func TestSpacerBasicRendering(t *testing.T) {
	renderer := NewInteractiveCharmRenderer()

	// Test spacer outside container (should render minimum size)
	spacer := dom.Spacer()
	output := renderer.Render(spacer)

	// Should render at least one space (don't trim spaces since that's what we're testing)
	if len(output) < 1 {
		t.Errorf("Spacer should render at least minimum size, got: '%s'", output)
	}

	// Should contain at least one space character
	if !strings.Contains(output, " ") {
		t.Errorf("Spacer should render space character, got: '%s'", output)
	}
}

func TestSpacerInContainer(t *testing.T) {
	tests := []struct {
		name          string
		terminalWidth int
		expectedMin   int // minimum expected content length
	}{
		{"Small terminal", 30, 25},
		{"Medium terminal", 50, 45},
		{"Large terminal", 80, 75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			window := &dom.Window{Width: tt.terminalWidth, Height: 10}

			layout := dom.Div(dom.DivProps{},
				dom.Text("Left"),
				dom.Spacer(),
				dom.Text("Right"),
			)
			layout.Window = window

			renderer := NewInteractiveCharmRenderer()
			output := renderer.Render(layout)

			// Extract the content line (should be line 0)
			lines := strings.Split(output, "\n")
			if len(lines) < 1 {
				t.Fatalf("Expected at least 1 line, got %d", len(lines))
			}

			contentLine := strings.TrimSpace(lines[0])

			// Verify content contains Left and Right
			if !strings.Contains(contentLine, "Left") {
				t.Errorf("Content should contain 'Left', got: '%s'", contentLine)
			}
			if !strings.Contains(contentLine, "Right") {
				t.Errorf("Content should contain 'Right', got: '%s'", contentLine)
			}

			// Verify content length is reasonable for terminal width
			if len(contentLine) < tt.expectedMin {
				t.Errorf("Content length %d should be at least %d for terminal width %d, got: '%s'",
					len(contentLine), tt.expectedMin, tt.terminalWidth, contentLine)
			}
		})
	}
}

func TestMultipleSpacers(t *testing.T) {
	window := &dom.Window{Width: 60, Height: 10}

	layout := dom.Div(dom.DivProps{},
		dom.Text("A"),
		dom.Spacer(),
		dom.Text("B"),
		dom.Spacer(),
		dom.Text("C"),
	)
	layout.Window = window

	renderer := NewInteractiveCharmRenderer()
	output := renderer.Render(layout)

	lines := strings.Split(output, "\n")
	if len(lines) < 1 {
		t.Fatalf("Expected at least 1 line, got %d", len(lines))
	}

	contentLine := strings.TrimSpace(lines[0])

	// Strip ANSI escape sequences to get the actual visual content
	visualContent := stripANSI(contentLine)

	// Calculate exact expected result
	terminalWidth := 60
	contentWidth := 3 // "A" + "B" + "C" = 3 characters
	spacerCount := 2  // Two spacers between A-B and B-C

	availableSpacerWidth := terminalWidth - contentWidth
	spacerWidth := availableSpacerWidth / spacerCount // Each spacer gets equal space

	// Construct exact expected string
	expected := "A" + strings.Repeat(" ", spacerWidth) + "B" + strings.Repeat(" ", spacerWidth) + "C"

	// Verify exact match using visual content (without ANSI codes)
	if visualContent != expected {
		t.Errorf("Content mismatch:\nExpected: '%s' (length: %d)\nActual:   '%s' (length: %d)",
			expected, len(expected), visualContent, len(visualContent))

		// Additional debugging info
		t.Errorf("Terminal width: %d, Content width: %d, Available spacer width: %d, Spacer width: %d",
			terminalWidth, contentWidth, availableSpacerWidth, spacerWidth)
		t.Errorf("Raw content with ANSI: '%s'", contentLine)
	}
}

func TestSpacerWithDifferentContent(t *testing.T) {
	tests := []struct {
		name      string
		leftText  string
		rightText string
	}{
		{"Short content", "A", "B"},
		{"Medium content", "Hello", "World"},
		{"Long content", "VeryLongLeftText", "VeryLongRightText"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			window := &dom.Window{Width: 80, Height: 10}

			layout := dom.Div(dom.DivProps{},
				dom.Text(tt.leftText),
				dom.Spacer(),
				dom.Text(tt.rightText),
			)
			layout.Window = window

			renderer := NewInteractiveCharmRenderer()
			output := renderer.Render(layout)

			lines := strings.Split(output, "\n")
			if len(lines) < 1 {
				t.Fatalf("Expected at least 1 line, got %d", len(lines))
			}

			contentLine := strings.TrimSpace(lines[0])

			// Verify both texts are present
			if !strings.Contains(contentLine, tt.leftText) {
				t.Errorf("Content should contain '%s', got: '%s'", tt.leftText, contentLine)
			}
			if !strings.Contains(contentLine, tt.rightText) {
				t.Errorf("Content should contain '%s', got: '%s'", tt.rightText, contentLine)
			}

			// Verify spacer creates separation
			expectedMinLength := len(tt.leftText) + len(tt.rightText) + 1 // +1 for at least one space
			if len(contentLine) < expectedMinLength {
				t.Errorf("Content should be at least %d chars, got %d: '%s'",
					expectedMinLength, len(contentLine), contentLine)
			}
		})
	}
}

func TestSpacerLayoutCalculation(t *testing.T) {
	// Test that spacer properly calculates available space
	window := &dom.Window{Width: 50, Height: 10}

	layout := dom.Div(dom.DivProps{},
		dom.Text("X"), // 1 char
		dom.Spacer(),
		dom.Text("Y"), // 1 char
	)
	layout.Window = window

	renderer := NewInteractiveCharmRenderer()
	output := renderer.Render(layout)

	lines := strings.Split(output, "\n")
	if len(lines) < 1 {
		t.Fatalf("Expected at least 1 line, got %d", len(lines))
	}

	contentLine := strings.TrimSpace(lines[0])

	// With terminal width 50, and content "X" + "Y" (2 chars),
	// the spacer should fill most of the remaining space
	// Content line should be close to terminal width
	if len(contentLine) < 40 { // Allow some margin for styling
		t.Errorf("Spacer should expand to fill available space. Expected length ~50, got %d: '%s'",
			len(contentLine), contentLine)
	}

	// Verify X and Y are far apart (spacer working)
	xPos := strings.Index(contentLine, "X")
	yPos := strings.Index(contentLine, "Y")
	if xPos == -1 || yPos == -1 {
		t.Fatalf("Both X and Y should be present in content: '%s'", contentLine)
	}

	distance := yPos - xPos
	if distance < 10 { // Should have significant spacing
		t.Errorf("X and Y should be well separated by spacer, distance: %d, content: '%s'",
			distance, contentLine)
	}
}

func TestSpacerWithNoWindow(t *testing.T) {
	// Test spacer behavior when no window is set (should use default width)
	layout := dom.Div(dom.DivProps{},
		dom.Text("Left"),
		dom.Spacer(),
		dom.Text("Right"),
	)
	// No window set

	renderer := NewInteractiveCharmRenderer()
	output := renderer.Render(layout)

	// Should not crash and should produce some output
	if len(output) == 0 {
		t.Error("Spacer should render something even without window")
	}

	lines := strings.Split(output, "\n")
	if len(lines) < 1 {
		t.Fatalf("Expected at least 1 line, got %d", len(lines))
	}

	contentLine := strings.TrimSpace(lines[0])

	// Should still contain the text elements
	if !strings.Contains(contentLine, "Left") || !strings.Contains(contentLine, "Right") {
		t.Errorf("Content should contain both Left and Right, got: '%s'", contentLine)
	}
}
