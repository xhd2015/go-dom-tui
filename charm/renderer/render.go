package renderer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/xhd2015/go-dom-tui/dom"
	"github.com/xhd2015/go-dom-tui/log"
	"github.com/xhd2015/go-dom-tui/react"
	"github.com/xhd2015/go-dom-tui/styles"
)

// stripANSI removes ANSI escape sequences from a string to get visual width
func stripANSI(str string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(str, "")
}

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
	case dom.ElementTypeText:
		cr.renderTextNode(vnode)
	case dom.ElementTypeDiv:
		cr.renderContainer(vnode, depth)
	case dom.ElementTypeH1:
		cr.renderTitle(vnode)
	case dom.ElementTypeH2:
		cr.renderSubtitle(vnode)
	case dom.ElementTypeP:
		cr.renderText(vnode)
	case dom.ElementTypeButton:
		cr.renderButton(vnode)
	case dom.ElementTypeInput:
		cr.renderInput(vnode)
	case dom.ElementTypeUl:
		cr.renderList(vnode, depth)
	case dom.ElementTypeLi:
		cr.renderListItem(vnode)
	case dom.ElementTypeBr:
		cr.renderBr(vnode)
	case dom.ElementTypeSpacer:
		cr.renderSpacer(vnode, depth)
	case dom.ElementTypeFragment:
		cr.renderFragment(vnode)
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

func (cr *InteractiveCharmRenderer) renderFragment(vnode *dom.Node) {
	for _, child := range vnode.Children {
		cr.renderNode(child, 0)
	}
}

func (cr *InteractiveCharmRenderer) renderSpacer(vnode *dom.Node, depth int) {
	// Spacer rendering is now handled by the container's layout system
	// This method should only be called for spacers outside of containers
	// In that case, render minimum space
	props := dom.ExtractProps[dom.SpacerProps](vnode.Props)
	minSize := props.MinSize
	if minSize <= 0 {
		minSize = 1
	}

	for i := 0; i < minSize; i++ {
		cr.output += " "
	}
}

// renderContainer renders a container div with border
func (cr *InteractiveCharmRenderer) renderContainer(vnode *dom.Node, depth int) {
	var content strings.Builder

	// Check if container has spacers - if so, use layout calculation
	if cr.hasSpacers(vnode.Children) {
		cr.renderContainerWithLayout(vnode, depth, &content)
	} else {
		// Simple rendering for containers without spacers
		for _, child := range vnode.Children {
			childRenderer := &InteractiveCharmRenderer{styles: cr.styles}
			childRenderer.renderNode(child, depth+1)
			content.WriteString(childRenderer.output)
		}
	}

	style := cr.getNodeStyle(vnode)

	rendered := style.Render(content.String())
	cr.output += rendered + "\n"
}

// hasSpacers checks if any child nodes are spacers
func (cr *InteractiveCharmRenderer) hasSpacers(children []*dom.Node) bool {
	for _, child := range children {
		if child != nil && child.Type == dom.ElementTypeSpacer {
			return true
		}
	}
	return false
}

// renderContainerWithLayout renders container with spacer layout calculation
func (cr *InteractiveCharmRenderer) renderContainerWithLayout(vnode *dom.Node, depth int, content *strings.Builder) {
	// Get available width using guard clauses (early return pattern)
	availableWidth := 80 // Default fallback

	// Extract props width (0 if not specified or not DivProps)
	propsWidth := 0
	if vnode.Props != nil {
		if divProps, ok := vnode.Props.(dom.StructProps[dom.DivProps]); ok {
			propsWidth = divProps.Value.Width
		}
	}

	// If props specifies width, use it
	if propsWidth > 0 {
		availableWidth = propsWidth
	} else {
		// Props width is 0, use window width if available
		if vnode.Window != nil {
			windowWidth := vnode.Window.GetWidth()
			if windowWidth > 0 {
				availableWidth = windowWidth
			}
		}
		// Otherwise keep default (80)
	}

	// First pass: render non-spacer elements and calculate their total width
	var nonSpacerElements []string
	var spacerIndices []int
	totalNonSpacerWidth := 0

	for i, child := range vnode.Children {
		if child == nil {
			continue
		}

		if child.Type == dom.ElementTypeSpacer {
			spacerIndices = append(spacerIndices, i)
			nonSpacerElements = append(nonSpacerElements, "") // placeholder
		} else {
			childRenderer := &InteractiveCharmRenderer{styles: cr.styles}
			childRenderer.renderNode(child, depth+1)
			rendered := childRenderer.output
			nonSpacerElements = append(nonSpacerElements, rendered)

			// Calculate width (for horizontal layout, we sum up the content widths)
			lines := strings.Split(rendered, "\n")
			if len(lines) > 0 {
				// For horizontal layout, use the first non-empty line's visual width (without ANSI codes)
				for _, line := range lines {
					trimmed := strings.TrimSpace(line)
					if trimmed != "" {
						visualWidth := len(stripANSI(trimmed))
						totalNonSpacerWidth += visualWidth
						break
					}
				}
			}
		}
	}

	// Calculate space available for spacers
	spacerCount := len(spacerIndices)
	if spacerCount > 0 {
		availableSpacerWidth := availableWidth - totalNonSpacerWidth
		if availableSpacerWidth < spacerCount {
			availableSpacerWidth = spacerCount // Minimum 1 space per spacer
		}
		spacerWidth := availableSpacerWidth / spacerCount

		// Render all elements with calculated spacer widths
		for _, element := range nonSpacerElements {
			if element == "" { // This is a spacer
				for j := 0; j < spacerWidth; j++ {
					content.WriteString(" ")
				}
			} else {
				content.WriteString(element)
			}
		}
	} else {
		// No spacers, just render normally
		for _, element := range nonSpacerElements {
			content.WriteString(element)
		}
	}
}

// renderTitle renders an h1 element
func (cr *InteractiveCharmRenderer) renderTitle(vnode *dom.Node) {
	cr.renderNodeText(vnode)
}

func (cr *InteractiveCharmRenderer) getNodeStyle(vnode *dom.Node) lipgloss.Style {
	var nodeStyle styles.Style
	var hasNodeStyle bool
	if styleValue, ok := vnode.Props.Get("style"); ok {
		if propStyle, ok := styleValue.(styles.Style); ok {
			if propStyle.NoDefault {
				return domStyleToCharmStyle(lipgloss.NewStyle(), propStyle)
			}
			nodeStyle = propStyle
			hasNodeStyle = true
		}
	}

	baseStyle := cr.styles.NoBorderDiv
	switch vnode.Type {
	case dom.ElementTypeH1:
		baseStyle = cr.styles.Title
	case dom.ElementTypeH2:
		baseStyle = cr.styles.Subtitle
	case dom.ElementTypeP:
		baseStyle = cr.styles.Text
	case dom.ElementTypeButton:
		baseStyle = cr.styles.Button
	case dom.ElementTypeInput:
		baseStyle = cr.styles.Input
	case dom.ElementTypeLi:
		baseStyle = cr.styles.CompactText
	case dom.ElementTypeText:
		baseStyle = cr.styles.Text
	}

	if hasNodeStyle {
		return domStyleToCharmStyle(baseStyle, nodeStyle)
	}

	return baseStyle
}

func (cr *InteractiveCharmRenderer) renderNodeText(vnode *dom.Node) {
	text := cr.extractRenderedText(vnode)
	style := cr.getNodeStyle(vnode)
	rendered := style.Render(text)
	cr.output += rendered + "\n"
}

// renderSubtitle renders an h2 element
func (cr *InteractiveCharmRenderer) renderSubtitle(vnode *dom.Node) {
	text := cr.extractRenderedText(vnode)
	rendered := cr.styles.Subtitle.Render(text)
	cr.output += rendered + "\n"
}

// renderText renders a p element
func (cr *InteractiveCharmRenderer) renderText(vnode *dom.Node) {
	cr.renderNodeText(vnode)
}

// renderButton renders a button element
func (cr *InteractiveCharmRenderer) renderButton(vnode *dom.Node) {
	text := cr.extractRenderedText(vnode)
	rendered := cr.styles.Button.Render(text)
	cr.output += rendered + "\n"
}

// renderInput renders an input element using Charm's textinput component
func (cr *InteractiveCharmRenderer) renderInput(vnode *dom.Node) {
	// Cast props to InputProps to avoid using GetOK, Get etc.
	props := dom.ExtractProps[dom.InputProps](vnode.Props)

	// Render input using static styling (no live textinput component)

	// Set default values
	placeholder := "Enter text..."
	value := ""
	inputType := "text"
	if props.InputType != "" {
		inputType = props.InputType
	}

	// Use typed props directly instead of GetOK/Get
	if props.Placeholder != "" {
		placeholder = props.Placeholder
	}
	value = props.Value
	// inputType remains "text" as there's no Type field in InputComponentProps

	// Create a textinput component
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(value)
	ti.CharLimit = 156

	// Set width based on props or window width
	if props.Width > 0 {
		ti.Width = props.Width
	} else {
		// Use window width if available, otherwise default to 50
		if vnode.Window != nil {
			windowWidth := vnode.Window.GetWidth()
			if windowWidth > 0 {
				// Leave some margin for styling and borders
				ti.Width = windowWidth - 10
				if ti.Width < 20 {
					ti.Width = 20 // Minimum width
				}
			} else {
				ti.Width = 50 // Default fallback
			}
		} else {
			ti.Width = 50 // Default fallback
		}
	}

	ti.SetCursor(props.CursorPosition)

	// Set password mode if it's a password field
	if inputType == "password" {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}

	// Style the textinput to match our theme
	ti.PromptStyle = cr.styles.Prompt
	ti.TextStyle = cr.styles.InputText
	ti.PlaceholderStyle = cr.styles.Text.Foreground(lipgloss.Color("#626262")).Italic(true)

	// Only call Focus() when the props indicate the element is focused, otherwise call Blur()
	if props.Focused {
		ti.Focus()
	} else {
		ti.Blur()
	}

	// Render the textinput component
	rendered := cr.styles.Input.Render(ti.View())
	cr.output += rendered + "\n"
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

	text := cr.extractRenderedText(vnode)

	var prefix string
	if props.ItemPrefix != nil {
		prefix = *props.ItemPrefix
	} else {
		prefix = "• "
		if props.Selected {
			// For focused items, use "> " prefix instead of bullet
			prefix = "> "
		}
	}
	renderedText := prefix + text

	// Style differently for selected/focused items
	// Use compact styles for list items (assuming they're part of todo list)
	if props.Selected {
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
func (cr *InteractiveCharmRenderer) extractRenderedText(vnode *dom.Node) string {
	var text strings.Builder

	for _, child := range vnode.Children {
		if child == nil {
			continue
		}
		if child.Type == "text" {
			text.WriteString(cr.extractTextNode(child))
		} else {
			text.WriteString(cr.extractRenderedText(child))
		}
	}

	return text.String()
}

func (cr *InteractiveCharmRenderer) extractTextNode(vnode *dom.Node) string {
	text := vnode.Text
	if text == "" {
		return ""
	}
	style := cr.getNodeStyle(vnode)
	return style.Render(text)

}
func (cr *InteractiveCharmRenderer) renderTextNode(vnode *dom.Node) {
	cr.output += cr.extractTextNode(vnode)
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
