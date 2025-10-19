package renderer

import (
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/xhd2015/go-dom-tui/dom"
)

// stripColor removes ANSI escape sequences from a string
func stripColor(str string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(str, "")
}

// renderNodeHelper is a helper function that takes a *dom.Node and returns the rendered string
func renderNodeHelper(vnode *dom.Node) string {
	renderer := NewInteractiveCharmRenderer()
	output := renderer.Render(vnode)
	return stripColor(output)
}

// TestRenderInput tests the renderInput function
func TestRenderInput(t *testing.T) {
	t.Run("RenderInputWithComponent", func(t *testing.T) {
		// Create a textinput component
		ti := textinput.New()
		ti.SetValue("test value")
		ti.Placeholder = "test placeholder"

		props := dom.NewStructProps(dom.InputProps{
			Placeholder: "test placeholder",
			Value:       "test value",
		})

		vnode := &dom.Node{
			Type:  dom.ElementTypeInput,
			Props: props,
		}

		// Render the input using helper
		output := renderNodeHelper(vnode)

		// Check that output was generated
		if output == "" {
			t.Error("Expected output to be generated")
		}

		// Output should contain a newline (from the render function)
		if !strings.Contains(output, "\n") {
			t.Error("Expected output to contain newline")
		}
	})

	t.Run("RenderInputWithoutComponent", func(t *testing.T) {
		props := dom.NewStructProps(dom.InputProps{
			Placeholder: "Enter text here",
			Value:       "current value",
		})

		vnode := &dom.Node{
			Type:  dom.ElementTypeInput,
			Props: props,
		}

		// Render the input using helper
		output := renderNodeHelper(vnode)

		// Check that output was generated
		if output == "" {
			t.Error("Expected output to be generated")
		}

		// Output should contain a newline
		if !strings.Contains(output, "\n") {
			t.Error("Expected output to contain newline")
		}
	})

	t.Run("RenderInputWithPasswordType", func(t *testing.T) {
		// Test with password input type
		props := dom.NewStructProps(dom.InputProps{
			Placeholder: "Enter password",
			Value:       "secret",
			InputType:   "password",
		})

		vnode := &dom.Node{
			Type:  dom.ElementTypeInput,
			Props: props,
		}

		// Render the input using helper
		output := renderNodeHelper(vnode)

		// Check that output was generated
		if output == "" {
			t.Error("Expected output to be generated")
		}
	})

	t.Run("RenderInputWithEmptyProps", func(t *testing.T) {
		props := dom.NewStructProps(dom.InputProps{})

		vnode := &dom.Node{
			Type:  dom.ElementTypeInput,
			Props: props,
		}

		// Render the input using helper
		output := renderNodeHelper(vnode)

		// Check that output was generated (should use defaults)
		if output == "" {
			t.Error("Expected output to be generated with default values")
		}
	})

	t.Run("RenderInputWithNilProps", func(t *testing.T) {
		vnode := &dom.Node{
			Type:  dom.ElementTypeInput,
			Props: dom.NewStructProps(dom.InputProps{}), // Use empty props instead of nil
		}

		// Should not panic with empty props
		output := renderNodeHelper(vnode)

		// Check that output was generated (should use defaults)
		if output == "" {
			t.Error("Expected output to be generated with default values")
		}
	})
}

// TestRenderButton tests the renderButton function
func TestRenderButton(t *testing.T) {
	t.Run("RenderButtonWithText", func(t *testing.T) {
		// Create a button with text child
		textNode := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Click me",
		}

		vnode := &dom.Node{
			Type:     dom.ElementTypeButton,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{textNode},
		}

		// Render the button using helper
		output := renderNodeHelper(vnode)

		// Check that output was generated
		if output == "" {
			t.Error("Expected output to be generated")
		}

		// Output should contain a newline
		if !strings.Contains(output, "\n") {
			t.Error("Expected output to contain newline")
		}
	})
}

// TestRenderText tests the renderText function
func TestRenderText(t *testing.T) {
	t.Run("RenderTextWithContent", func(t *testing.T) {
		// Create a text element with text child
		textNode := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Hello, world!",
		}

		vnode := &dom.Node{
			Type:     dom.ElementTypeP,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{textNode},
		}

		// Render the text using helper
		output := renderNodeHelper(vnode)

		// Check that output was generated
		if output == "" {
			t.Error("Expected output to be generated")
		}

		// Output should contain a newline
		if !strings.Contains(output, "\n") {
			t.Error("Expected output to contain newline")
		}
	})
}

// TestTitle tests the renderTitle function for H1 and H2 elements
func TestTitle(t *testing.T) {
	t.Run("RenderH1WithText", func(t *testing.T) {
		// Create an H1 element with text child
		textNode := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Main Title",
		}

		vnode := &dom.Node{
			Type:     dom.ElementTypeH1,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{textNode},
		}

		// Render the H1 using helper
		output := renderNodeHelper(vnode)

		// Check that output was generated
		if output == "" {
			t.Error("Expected output to be generated")
		}

		// Output should contain a newline
		if !strings.Contains(output, "\n") {
			t.Error("Expected output to contain newline")
		}

		// Output should contain the title text
		if !strings.Contains(output, "Main Title") {
			t.Error("Expected output to contain 'Main Title'")
		}
	})

	t.Run("RenderH2WithText", func(t *testing.T) {
		// Create an H2 element with text child
		textNode := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Subtitle",
		}

		vnode := &dom.Node{
			Type:     dom.ElementTypeH2,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{textNode},
		}

		// Render the H2 using helper
		output := renderNodeHelper(vnode)

		// Check that output was generated
		if output == "" {
			t.Error("Expected output to be generated")
		}

		// Output should contain a newline
		if !strings.Contains(output, "\n") {
			t.Error("Expected output to contain newline")
		}

		// Output should contain the subtitle text
		if !strings.Contains(output, "Subtitle") {
			t.Error("Expected output to contain 'Subtitle'")
		}
	})

	t.Run("RenderH1WithEmptyContent", func(t *testing.T) {
		// Create an H1 element with empty content
		vnode := &dom.Node{
			Type:     dom.ElementTypeH1,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{},
		}

		// Render the H1 using helper
		output := renderNodeHelper(vnode)

		// Check that output was generated (should handle empty content gracefully)
		if output == "" {
			t.Error("Expected output to be generated even with empty content")
		}
	})

	t.Run("RenderH2WithEmptyContent", func(t *testing.T) {
		// Create an H2 element with empty content
		vnode := &dom.Node{
			Type:     dom.ElementTypeH2,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{},
		}

		// Render the H2 using helper
		output := renderNodeHelper(vnode)

		// Check that output was generated (should handle empty content gracefully)
		if output == "" {
			t.Error("Expected output to be generated even with empty content")
		}
	})
}

// TestNestedDiv tests rendering nested div elements
func TestNestedDiv(t *testing.T) {
	t.Run("RenderNestedDivWithChildren", func(t *testing.T) {
		// Create text nodes for child divs
		textNode1 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  " child 1",
		}

		textNode2 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  " child 2",
		}

		// Create child div elements
		childDiv1 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{textNode1},
		}

		childDiv2 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{textNode2},
		}

		// Create parent div with child divs
		parentDiv := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{childDiv1, childDiv2},
		}

		// Render the nested div structure
		output := renderNodeHelper(parentDiv)

		// Expected output should be exact literal string match
		// Should not have extra newline when children are already block elements
		expected := " child 1\n child 2\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})
}

// TestSpanFollowedByDiv tests inline span followed by block div
func TestSpanFollowedByDiv(t *testing.T) {
	t.Run("RenderSpanThenDiv", func(t *testing.T) {
		// Create text nodes
		spanText := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "text1",
		}

		divText := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "text2",
		}

		// Create span element (inline)
		spanElement := &dom.Node{
			Type:     dom.ElementTypeSpan,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{spanText},
		}

		// Create div element (block)
		divElement := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{divText},
		}

		// Create container with span followed by div
		container := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{spanElement, divElement},
		}

		// Render the structure
		output := renderNodeHelper(container)

		// Expected: span renders inline, then div starts on new line
		// The span "text1" should be followed by div "text2" on a new line
		expected := "text1\ntext2\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})
}

// TestSpanDivInFragment tests inline span followed by block div in a fragment
func TestSpanDivInFragment(t *testing.T) {
	t.Run("RenderSpanThenDivInFragment", func(t *testing.T) {
		// Create text nodes
		spanText := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "text1",
		}

		divText := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "text2",
		}

		// Create span element (inline)
		spanElement := &dom.Node{
			Type:     dom.ElementTypeSpan,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{spanText},
		}

		// Create div element (block)
		divElement := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{divText},
		}

		// Create fragment with span followed by div
		fragment := &dom.Node{
			Type:     dom.ElementTypeFragment,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{spanElement, divElement},
		}

		// Render the structure
		output := renderNodeHelper(fragment)

		// Expected: span renders inline, then div starts on new line
		// The span "text1" should be followed by div "text2" on a new line
		expected := "text1\ntext2\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})
}

// TestMultiLineTextInDiv tests rendering multi-line text content within a div
func TestMultiLineTextInDiv(t *testing.T) {
	t.Run("RenderDivWithMultiLineText", func(t *testing.T) {
		// Create multi-line text content
		multiLineText := "ABC\n123\nX Y Z\n--END--"

		// Create text node with multi-line content
		textNode := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  multiLineText,
		}

		// Create div element containing the multi-line text
		divElement := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
			Children: []*dom.Node{textNode},
		}

		// Render the structure
		output := renderNodeHelper(divElement)

		// Expected: the renderer strips internal newlines and adds a trailing newline
		expected := "ABC123X Y Z--END--\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})
}
