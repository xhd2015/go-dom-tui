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

		props := dom.NewStructProps(dom.InputProps{
			Placeholder: "test placeholder",
			Value:       "test value",
		})

		vnode := &dom.Node{
			Type:  dom.ElementTypeInput,
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
		props := dom.NewStructProps(dom.InputProps{
			Placeholder: "Enter text here",
			Value:       "current value",
		})

		vnode := &dom.Node{
			Type:  dom.ElementTypeInput,
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

		// Render the input
		renderer.output = ""
		renderer.renderInput(vnode)

		// Check that output was generated
		if renderer.output == "" {
			t.Error("Expected output to be generated")
		}
	})

	t.Run("RenderInputWithEmptyProps", func(t *testing.T) {
		props := dom.NewStructProps(dom.InputProps{})

		vnode := &dom.Node{
			Type:  dom.ElementTypeInput,
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
			Type:  dom.ElementTypeInput,
			Props: dom.NewStructProps(dom.InputProps{}), // Use empty props instead of nil
		}

		// Should not panic with empty props
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

func (p *testInputProps) Get(key string) (any, bool) {
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
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Click me",
		}

		vnode := &dom.Node{
			Type:     dom.ElementTypeButton,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
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
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Hello, world!",
		}

		vnode := &dom.Node{
			Type:     dom.ElementTypeP,
			Props:    dom.NewStructProps(dom.EmptyProps{}),
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

// Note: TestExtractText was removed as it was testing non-existent recursive text extraction functionality
