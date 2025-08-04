package renderer

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/xhd2015/go-dom-tui/dom"
)

// TestRenderInput tests the renderInput function
func TestRenderInput(t *testing.T) {
	renderer := NewInteractiveCharmRenderer()

	t.Run("RenderInputWithComponent", func(t *testing.T) {
		// Create a textinput component
		ti := textinput.New()
		ti.SetValue("test value")
		ti.Placeholder = "test placeholder"

		props := &dom.InputComponentProps{
			Placeholder: "test placeholder",
			Value:       "test value",
		}

		vnode := &dom.Node{
			Type:  "input",
			Props: props,
		}

		// Render the input
		renderer.output = ""
		renderer.renderInput(vnode)

		// Check that output was generated
		if renderer.output == "" {
			t.Error("Expected output to be generated")
		}

		// Output should contain a newline (from the render function)
		if !strings.Contains(renderer.output, "\n") {
			t.Error("Expected output to contain newline")
		}
	})

	t.Run("RenderInputWithoutComponent", func(t *testing.T) {
		props := &dom.InputComponentProps{
			Placeholder: "Enter text here",
			Value:       "current value",
		}

		vnode := &dom.Node{
			Type:  "input",
			Props: props,
		}

		// Render the input
		renderer.output = ""
		renderer.renderInput(vnode)

		// Check that output was generated
		if renderer.output == "" {
			t.Error("Expected output to be generated")
		}

		// Output should contain a newline
		if !strings.Contains(renderer.output, "\n") {
			t.Error("Expected output to contain newline")
		}
	})

	t.Run("RenderInputWithPasswordType", func(t *testing.T) {
		// Test with a custom props type that supports type
		customProps := &testInputProps{
			placeholder: "Enter password",
			value:       "secret",
			inputType:   "password",
		}

		vnode := &dom.Node{
			Type:  "input",
			Props: customProps,
		}

		// Render the input
		renderer.output = ""
		renderer.renderInput(vnode)

		// Check that output was generated
		if renderer.output == "" {
			t.Error("Expected output to be generated")
		}
	})

	t.Run("RenderInputWithEmptyProps", func(t *testing.T) {
		props := &dom.InputComponentProps{}

		vnode := &dom.Node{
			Type:  "input",
			Props: props,
		}

		// Render the input
		renderer.output = ""
		renderer.renderInput(vnode)

		// Check that output was generated (should use defaults)
		if renderer.output == "" {
			t.Error("Expected output to be generated with default values")
		}
	})

	t.Run("RenderInputWithNilProps", func(t *testing.T) {
		vnode := &dom.Node{
			Type:  "input",
			Props: nil,
		}

		// Should not panic with nil props
		renderer.output = ""
		renderer.renderInput(vnode)

		// Check that output was generated (should use defaults)
		if renderer.output == "" {
			t.Error("Expected output to be generated with default values")
		}
	})
}

// testInputProps is a custom props type for testing password input
type testInputProps struct {
	placeholder string
	value       string
	inputType   string
}

func (p *testInputProps) Get(key string) any {
	switch key {
	case "placeholder":
		return p.placeholder
	case "value":
		return p.value
	case "type":
		return p.inputType
	}
	return nil
}

func (p *testInputProps) GetString(key string) string {
	switch key {
	case "placeholder":
		return p.placeholder
	case "value":
		return p.value
	case "type":
		return p.inputType
	}
	return ""
}

func (p *testInputProps) GetOK(key string) (any, bool) {
	switch key {
	case "placeholder":
		return p.placeholder, p.placeholder != ""
	case "value":
		return p.value, true
	case "type":
		return p.inputType, p.inputType != ""
	}
	return nil, false
}

func (p *testInputProps) Range(fn func(key string, value any) bool) {
	if p.placeholder != "" {
		fn("placeholder", p.placeholder)
	}
	fn("value", p.value)
	if p.inputType != "" {
		fn("type", p.inputType)
	}
}

func (p *testInputProps) Clone() dom.Props {
	return &testInputProps{
		placeholder: p.placeholder,
		value:       p.value,
		inputType:   p.inputType,
	}
}

// TestRenderButton tests the renderButton function
func TestRenderButton(t *testing.T) {
	renderer := NewInteractiveCharmRenderer()

	t.Run("RenderButtonWithText", func(t *testing.T) {
		// Create a button with text child
		textNode := &dom.Node{
			Type:  "text",
			Props: &dom.TextNodeProps{Text: "Click me"},
		}

		vnode := &dom.Node{
			Type:     "button",
			Props:    &dom.EmptyProps{},
			Children: []*dom.Node{textNode},
		}

		// Render the button
		renderer.output = ""
		renderer.renderButton(vnode)

		// Check that output was generated
		if renderer.output == "" {
			t.Error("Expected output to be generated")
		}

		// Output should contain a newline
		if !strings.Contains(renderer.output, "\n") {
			t.Error("Expected output to contain newline")
		}
	})
}

// TestRenderText tests the renderText function
func TestRenderText(t *testing.T) {
	renderer := NewInteractiveCharmRenderer()

	t.Run("RenderTextWithContent", func(t *testing.T) {
		// Create a text element with text child
		textNode := &dom.Node{
			Type:  "text",
			Props: &dom.TextNodeProps{Text: "Hello, world!"},
		}

		vnode := &dom.Node{
			Type:     "p",
			Props:    &dom.EmptyProps{},
			Children: []*dom.Node{textNode},
		}

		// Render the text
		renderer.output = ""
		renderer.renderText(vnode)

		// Check that output was generated
		if renderer.output == "" {
			t.Error("Expected output to be generated")
		}

		// Output should contain a newline
		if !strings.Contains(renderer.output, "\n") {
			t.Error("Expected output to contain newline")
		}
	})
}

// TestExtractText tests the extractText helper function
func TestExtractText(t *testing.T) {
	renderer := NewInteractiveCharmRenderer()

	t.Run("ExtractTextFromSingleChild", func(t *testing.T) {
		textNode := &dom.Node{
			Type:  "text",
			Props: &dom.TextNodeProps{Text: "Hello"},
		}

		vnode := &dom.Node{
			Type:     "p",
			Props:    &dom.EmptyProps{},
			Children: []*dom.Node{textNode},
		}

		text := renderer.extractText(vnode)
		if text != "Hello" {
			t.Errorf("Expected 'Hello', got '%s'", text)
		}
	})

	t.Run("ExtractTextFromMultipleChildren", func(t *testing.T) {
		textNode1 := &dom.Node{
			Type:  "text",
			Props: &dom.TextNodeProps{Text: "Hello "},
		}
		textNode2 := &dom.Node{
			Type:  "text",
			Props: &dom.TextNodeProps{Text: "world!"},
		}

		vnode := &dom.Node{
			Type:     "p",
			Props:    &dom.EmptyProps{},
			Children: []*dom.Node{textNode1, textNode2},
		}

		text := renderer.extractText(vnode)
		if text != "Hello world!" {
			t.Errorf("Expected 'Hello world!', got '%s'", text)
		}
	})

	t.Run("ExtractTextFromNestedChildren", func(t *testing.T) {
		textNode := &dom.Node{
			Type:  "text",
			Props: &dom.TextNodeProps{Text: "nested"},
		}

		spanNode := &dom.Node{
			Type:     "span",
			Props:    &dom.EmptyProps{},
			Children: []*dom.Node{textNode},
		}

		vnode := &dom.Node{
			Type:     "p",
			Props:    &dom.EmptyProps{},
			Children: []*dom.Node{spanNode},
		}

		text := renderer.extractText(vnode)
		if text != "nested" {
			t.Errorf("Expected 'nested', got '%s'", text)
		}
	})

	t.Run("ExtractTextFromEmptyNode", func(t *testing.T) {
		vnode := &dom.Node{
			Type:     "p",
			Props:    &dom.EmptyProps{},
			Children: []*dom.Node{},
		}

		text := renderer.extractText(vnode)
		if text != "" {
			t.Errorf("Expected empty string, got '%s'", text)
		}
	})
}
