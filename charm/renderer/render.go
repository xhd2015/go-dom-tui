package renderer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/xhd2015/go-dom-tui/dom"
	"github.com/xhd2015/go-dom-tui/log"
	"github.com/xhd2015/go-dom-tui/react"
)

// InteractiveCharmRenderer implements the Renderer interface with user interaction
type InteractiveCharmRenderer struct {
	output string
	styles CharmStyles
}

// NewInteractiveCharmRenderer creates a new interactive renderer with styled components
func NewInteractiveCharmRenderer() *InteractiveCharmRenderer {
	return &InteractiveCharmRenderer{
		styles: defaultStyles(),
	}
}

// Render renders a VNode using Charm CLI styling
func (cr *InteractiveCharmRenderer) Render(vnode *dom.Node) string {
	cr.output = ""

	// Update styles based on window size from VNode
	if vnode != nil && vnode.Window != nil {
		width, height := vnode.Window.Get()
		cr.updateStylesForWindowSize(width, height)
	}

	cr.renderNode(vnode, 0)
	return cr.output
}

// renderNode recursively renders a VNode with appropriate styling
func (cr *InteractiveCharmRenderer) renderNode(vnode *dom.Node, depth int) {
	if vnode == nil {
		return
	}

	switch vnode.Type {
	case "text":
		cr.renderTextNode(vnode)
	case "div":
		cr.renderContainer(vnode, depth)
	case "h1":
		cr.renderTitle(vnode)
	case "h2":
		cr.renderSubtitle(vnode)
	case "p":
		cr.renderText(vnode)
	case "button":
		cr.renderButton(vnode)
	case "input":
		cr.renderInput(vnode)
	case "ul":
		cr.renderList(vnode, depth)
	case "li":
		cr.renderListItem(vnode)
	case "br":
		cr.renderBr(vnode)
	case "component":
		panic("component is deprecated")
	default:
		log.Logf("renderNode called for unknown type: %s, depth: %d", vnode.Type, depth)
		cr.renderDefault(vnode, depth)
	}
}

func (cr *InteractiveCharmRenderer) renderBr(vnode *dom.Node) {
	cr.output += "\n"
}

// renderContainer renders a container div with border
func (cr *InteractiveCharmRenderer) renderContainer(vnode *dom.Node, depth int) {
	var content strings.Builder

	for _, child := range vnode.Children {
		childRenderer := &InteractiveCharmRenderer{styles: cr.styles}
		childRenderer.renderNode(child, depth+1)
		content.WriteString(childRenderer.output)
	}

	// Handle style object or fallback to old string-based styling
	var style lipgloss.Style
	if styleValue, ok := vnode.Props.Get("style"); ok {
		if divStyle, ok := styleValue.(dom.Style); ok {
			// New object-based styling
			if divStyle.BorderColor != "" {
				// Use container with custom border color
				style = cr.styles.Container.BorderForeground(lipgloss.Color(divStyle.BorderColor))
			} else {
				// No border specified, use no-border style
				style = cr.styles.NoBorderDiv
			}
		} else {
			style = cr.styles.NoBorderDiv // Default to no border
		}
	} else {
		style = cr.styles.NoBorderDiv // Default to no border
	}

	rendered := style.Render(content.String())
	cr.output += rendered + "\n"
}

// renderTitle renders an h1 element
func (cr *InteractiveCharmRenderer) renderTitle(vnode *dom.Node) {
	cr.renderNodeText(vnode)
}

func (cr *InteractiveCharmRenderer) getNodeStyle(vnode *dom.Node) lipgloss.Style {
	baseStyle := cr.styles.NoBorderDiv
	switch vnode.Type {
	case "h1":
		baseStyle = cr.styles.Title
	case "h2":
		baseStyle = cr.styles.Subtitle
	case "p":
		baseStyle = cr.styles.Text
	case "button":
		baseStyle = cr.styles.Button
	case "input":
		baseStyle = cr.styles.Input
	case "li":
		baseStyle = cr.styles.CompactText
	case "text":
		baseStyle = cr.styles.Text
	}
	if styleValue, ok := vnode.Props.Get("style"); ok {
		if divStyle, ok := styleValue.(dom.Style); ok {
			baseStyle = domStyleToCharmStyle(baseStyle, divStyle)
		}
	}
	return baseStyle
}

func (cr *InteractiveCharmRenderer) renderNodeText(vnode *dom.Node) {
	text := cr.extractText(vnode)
	style := cr.getNodeStyle(vnode)
	rendered := style.Render(text)
	cr.output += rendered + "\n"
}

// renderSubtitle renders an h2 element
func (cr *InteractiveCharmRenderer) renderSubtitle(vnode *dom.Node) {
	text := cr.extractText(vnode)
	rendered := cr.styles.Subtitle.Render(text)
	cr.output += rendered + "\n"
}

// renderText renders a p element
func (cr *InteractiveCharmRenderer) renderText(vnode *dom.Node) {
	cr.renderNodeText(vnode)
}

// renderButton renders a button element
func (cr *InteractiveCharmRenderer) renderButton(vnode *dom.Node) {
	text := cr.extractText(vnode)
	rendered := cr.styles.Button.Render(text)
	cr.output += rendered + "\n"
}

// renderInput renders an input element using Charm's textinput component
func (cr *InteractiveCharmRenderer) renderInput(vnode *dom.Node) {
	// Cast props to InputProps to avoid using GetOK, Get etc.
	inputProps := dom.ExtractProps[dom.InputProps](vnode.Props)

	// Render input using static styling (no live textinput component)

	// Set default values
	placeholder := "Enter text..."
	value := ""
	inputType := "text"

	// Use typed props directly instead of GetOK/Get
	if inputProps.Placeholder != "" {
		placeholder = inputProps.Placeholder
	}
	value = inputProps.Value
	// inputType remains "text" as there's no Type field in InputComponentProps

	// Create a textinput component
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(value)
	ti.CharLimit = 156
	ti.Width = 50
	ti.SetCursor(inputProps.CursorPosition)

	// Set password mode if it's a password field
	if inputType == "password" {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}

	// Style the textinput to match our theme
	ti.PromptStyle = cr.styles.Prompt
	ti.TextStyle = cr.styles.Text
	ti.PlaceholderStyle = cr.styles.Text.Foreground(lipgloss.Color("#626262"))

	// Only call Focus() when the props indicate the element is focused, otherwise call Blur()
	if inputProps.Focused {
		ti.Focus()
	} else {
		ti.Blur()
	}

	// Render the textinput component
	rendered := cr.styles.Input.Render(ti.View())
	cr.output += rendered + "\n"
}

// renderComponent renders a component
func (cr *InteractiveCharmRenderer) renderComponent(vnode *dom.Node) {
	log.Logf("renderComponent called, type: %s", vnode.Type)
	cr.renderNode(vnode.Children[0], 0)
}

// renderList renders a ul element
func (cr *InteractiveCharmRenderer) renderList(vnode *dom.Node, depth int) {
	for _, child := range vnode.Children {
		if child.Type == "li" {
			// Use the dedicated renderListItem function to properly handle li elements
			cr.renderListItem(child)
		} else {
			// Render other child elements
			cr.renderNode(child, depth+1)
		}
	}
}

// renderListItem renders a li element
func (cr *InteractiveCharmRenderer) renderListItem(vnode *dom.Node) {
	props := dom.ExtractProps[dom.ListItemProps](vnode.Props)

	text := props.Text
	if text == "" {
		text = cr.extractText(vnode)
	}

	prefix := "• "
	if props.Focused {
		// For focused items, use "> " prefix instead of bullet
		prefix = "> "
	}
	renderedText := prefix + text

	// Style differently for selected/focused items
	// Use compact styles for list items (assuming they're part of todo list)
	if props.Focused {
		cr.output += cr.styles.CompactSuccess.Render(renderedText) + "\n"
	} else {
		cr.output += cr.styles.CompactText.Render(renderedText) + "\n"
	}
}

// renderDefault renders unknown elements
func (cr *InteractiveCharmRenderer) renderDefault(vnode *dom.Node, depth int) {
	cr.output += fmt.Sprintf("<%s>\n", vnode.Type)
	for _, child := range vnode.Children {
		cr.renderNode(child, depth+1)
	}
	cr.output += fmt.Sprintf("</%s>\n", vnode.Type)
}

// extractText extracts text content from a VNode
func (cr *InteractiveCharmRenderer) extractText(vnode *dom.Node) string {
	var text strings.Builder

	for _, child := range vnode.Children {
		if child.Type == "text" {
			text.WriteString(child.Text)
		} else {
			text.WriteString(cr.extractText(child))
		}
	}

	return text.String()
}

func (cr *InteractiveCharmRenderer) renderTextNode(vnode *dom.Node) {
	text := vnode.Text
	if text == "" {
		return
	}
	style := cr.getNodeStyle(vnode)
	rendered := style.Render(text)
	cr.output += rendered
}

// ApplyPatch applies a patch to update the rendered output
func (cr *InteractiveCharmRenderer) ApplyPatch(patch react.Patch, currentOutput string) string {
	return currentOutput + fmt.Sprintf("<!-- Patch applied: %v -->\n", patch.Type)
}

// updateStylesForWindowSize adjusts styles based on provided terminal dimensions
func (cr *InteractiveCharmRenderer) updateStylesForWindowSize(width, height int) {
	// Only adjust if we have valid window dimensions
	if width > 0 {
		// Adjust input width to be responsive
		inputWidth := width - 8 // Leave some margin (reduced from 10 to 8)
		if inputWidth > 80 {
			inputWidth = 80 // Max width cap (increased from 50 to 80)
		}
		if inputWidth < 30 {
			inputWidth = 30 // Min width (increased from 20 to 30)
		}

		cr.styles.Input = cr.styles.Input.Width(inputWidth)

		// Set container to use full window width
		containerWidth := width - 2 // Leave minimal margin for border
		if containerWidth > 0 {
			cr.styles.Container = cr.styles.Container.Width(containerWidth).Margin(0, 0)

			// Adjust padding for smaller terminals
			if width < 80 {
				cr.styles.Container = cr.styles.Container.Padding(1, 2)
			} else {
				cr.styles.Container = cr.styles.Container.Padding(2, 3)
			}
		}
	}

	// Adjust for height if needed (future enhancement)
	_ = height
}
