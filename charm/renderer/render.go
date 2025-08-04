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

// CharmStyles holds various lipgloss styles for rendering
type CharmStyles struct {
	Title          lipgloss.Style
	Subtitle       lipgloss.Style
	Text           lipgloss.Style
	Button         lipgloss.Style
	Container      lipgloss.Style
	CompactDiv     lipgloss.Style
	CompactText    lipgloss.Style
	CompactSuccess lipgloss.Style
	Input          lipgloss.Style
	Prompt         lipgloss.Style
	Success        lipgloss.Style
	Error          lipgloss.Style
}

// NewInteractiveCharmRenderer creates a new interactive renderer with styled components
func NewInteractiveCharmRenderer() *InteractiveCharmRenderer {
	return &InteractiveCharmRenderer{
		styles: CharmStyles{
			Title: lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 2).
				Margin(1, 0),
			Subtitle: lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				Padding(0, 1).
				Margin(0, 0),
			Text: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262")).
				Margin(0, 1),
			Button: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFF")).
				Background(lipgloss.Color("#04B575")).
				Padding(0, 3).
				Margin(0, 1).
				Bold(true),
			Container: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#874BFD")).
				Padding(2, 3).
				Margin(1, 0),
			CompactDiv: lipgloss.NewStyle().
				Padding(0, 0).
				Margin(0, 0),
			CompactText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262")).
				Margin(0, 0),
			CompactSuccess: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF00")).
				Bold(true).
				Margin(0, 0),
			Input: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#04B575")).
				Padding(0, 1).
				Margin(0, 1).
				Width(40),
			Prompt: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFF00")).
				Bold(true).
				Margin(0, 1),
			Success: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF00")).
				Bold(true).
				Margin(0, 1),
			Error: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF0000")).
				Bold(true).
				Margin(0, 1),
		},
	}
}

// Render renders a VNode using Charm CLI styling
func (cr *InteractiveCharmRenderer) Render(vnode *dom.Node) string {
	cr.output = ""
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
		if value, ok := vnode.Props.Get("text"); ok {
			if textValue, ok := value.(string); ok {
				cr.output += textValue
			}
		}
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
	case "component":
		cr.renderComponent(vnode)
	default:
		log.Logf("renderNode called for unknown type: %s, depth: %d", vnode.Type, depth)
		cr.renderDefault(vnode, depth)
	}
}

// renderContainer renders a container div with border
func (cr *InteractiveCharmRenderer) renderContainer(vnode *dom.Node, depth int) {
	var content strings.Builder

	for _, child := range vnode.Children {
		childRenderer := &InteractiveCharmRenderer{styles: cr.styles}
		childRenderer.renderNode(child, depth+1)
		content.WriteString(childRenderer.output)
	}

	// Check for specific style names to use compact styling
	var style lipgloss.Style
	if styleValue, ok := vnode.Props.Get("style"); ok {
		if styleName, ok := styleValue.(string); ok {
			switch styleName {
			case "main-app", "quick-input", "todo-display":
				style = cr.styles.CompactDiv
			default:
				style = cr.styles.Container
			}
		} else {
			style = cr.styles.Container
		}
	} else {
		style = cr.styles.Container
	}

	rendered := style.Render(content.String())
	cr.output += rendered + "\n"
}

// renderTitle renders an h1 element
func (cr *InteractiveCharmRenderer) renderTitle(vnode *dom.Node) {
	text := cr.extractText(vnode)
	rendered := cr.styles.Title.Render(text)
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
	text := cr.extractText(vnode)
	rendered := cr.styles.Text.Render(text)
	cr.output += rendered + "\n"
}

// renderButton renders a button element
func (cr *InteractiveCharmRenderer) renderButton(vnode *dom.Node) {
	text := cr.extractText(vnode)
	rendered := cr.styles.Button.Render(text)
	cr.output += rendered + "\n"
}

// renderInput renders an input element using Charm's textinput component
func (cr *InteractiveCharmRenderer) renderInput(vnode *dom.Node) {
	// Cast props to InputComponentProps to avoid using GetOK, Get etc.
	inputProps := dom.ExtractProps[dom.InputComponentProps](vnode.Props)

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
	ti.Cursor.Style = cr.styles.Success

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
			if value, ok := child.Props.Get("text"); ok {
				if textValue, ok := value.(string); ok {
					text.WriteString(textValue)
				}
			}
		} else {
			text.WriteString(cr.extractText(child))
		}
	}

	return text.String()
}

// ApplyPatch applies a patch to update the rendered output
func (cr *InteractiveCharmRenderer) ApplyPatch(patch react.Patch, currentOutput string) string {
	return currentOutput + fmt.Sprintf("<!-- Patch applied: %v -->\n", patch.Type)
}
