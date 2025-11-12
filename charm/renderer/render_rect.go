package renderer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/xhd2015/go-dom-tui/dom"
	"github.com/xhd2015/go-dom-tui/log"
)

// RenderToRect renders a DOM node into a Rectangle using pure rectangle-based rendering
// This is a completely new implementation that doesn't rely on string concatenation
// width and height define the container dimensions for the rendering
// Returns a Rectangle containing the rendered output
func (cr *InteractiveCharmRenderer) RenderToRect(vnode *dom.Node, width, height int) Rectangle {
	if vnode == nil {
		return Rectangle{Width: 0, Height: 0, Lines: []string{}}
	}

	return cr.renderNodeToRect(vnode, width, height)
}

// renderNodeToRect recursively renders a VNode into a Rectangle
// width and height define the container dimensions available for this node
func (cr *InteractiveCharmRenderer) renderNodeToRect(vnode *dom.Node, width, height int) Rectangle {
	if vnode == nil {
		return Rectangle{Width: 0, Height: 0, Lines: []string{}}
	}

	switch vnode.Type {
	case dom.ElementTypeText:
		return cr.renderTextNodeToRect(vnode, width, height)
	case dom.ElementTypeDiv:
		return cr.renderContainerToRect(vnode, width, height)
	case dom.ElementTypeHDiv:
		return cr.renderHDivToRect(vnode, width, height)
	case dom.ElementTypeZDiv:
		return cr.renderZDivToRect(vnode, width, height)
	case dom.ElementTypeSpan:
		return cr.renderSpanToRect(vnode, width, height)
	case dom.ElementTypeH1:
		return cr.renderTitleToRect(vnode, width, height)
	case dom.ElementTypeH2:
		return cr.renderSubtitleToRect(vnode, width, height)
	case dom.ElementTypeP:
		return cr.renderTextToRect(vnode, width, height)
	case dom.ElementTypeButton:
		return cr.renderButtonToRect(vnode, width, height)
	case dom.ElementTypeInput:
		return cr.renderInputToRect(vnode, width, height)
	case dom.ElementTypeUl:
		return cr.renderListToRect(vnode, width, height)
	case dom.ElementTypeLi:
		return cr.renderListItemToRect(vnode, width, height)
	case dom.ElementTypeBr:
		return cr.renderBrToRect(vnode, width, height)
	case dom.ElementTypeSpacer:
		return cr.renderSpacerToRect(vnode, width, height)
	case dom.ElementTypeFixedSpacer:
		return cr.renderFixedSpacerToRect(vnode, width, height)
	case dom.ElementTypeFragment:
		return cr.renderFragmentToRect(vnode, width, height)
	default:
		log.Logf("renderNodeToRect called for unknown type: %s", vnode.Type)
		return cr.renderDefaultToRect(vnode, width, height)
	}
}

// renderTextNodeToRect renders a text node to a Rectangle
func (cr *InteractiveCharmRenderer) renderTextNodeToRect(vnode *dom.Node, width, height int) Rectangle {
	text := vnode.Text
	if text == "" {
		return Rectangle{Width: 0, Height: 0, Lines: []string{}}
	}
	style := cr.getNodeStyle(vnode)
	rendered := style.Render(text)
	return NewRectangle(rendered)
}

// renderSpanToRect renders a span element to a Rectangle
func (cr *InteractiveCharmRenderer) renderSpanToRect(vnode *dom.Node, width, height int) Rectangle {
	text := cr.extractRenderedText(vnode)
	style := cr.getNodeStyle(vnode)
	rendered := style.Render(text)
	return NewRectangle(rendered)
}

// renderTitleToRect renders an h1 element to a Rectangle
func (cr *InteractiveCharmRenderer) renderTitleToRect(vnode *dom.Node, width, height int) Rectangle {
	text := cr.extractRenderedText(vnode)
	style := cr.getNodeStyle(vnode)
	rendered := style.Render(text)
	return NewRectangle(rendered)
}

// renderSubtitleToRect renders an h2 element to a Rectangle
func (cr *InteractiveCharmRenderer) renderSubtitleToRect(vnode *dom.Node, width, height int) Rectangle {
	text := cr.extractRenderedText(vnode)
	rendered := cr.styles.Subtitle.Render(text)
	return NewRectangle(rendered)
}

// renderTextToRect renders a p element to a Rectangle
func (cr *InteractiveCharmRenderer) renderTextToRect(vnode *dom.Node, width, height int) Rectangle {
	text := cr.extractRenderedText(vnode)
	style := cr.getNodeStyle(vnode)
	rendered := style.Render(text)
	return NewRectangle(rendered)
}

// renderButtonToRect renders a button element to a Rectangle
func (cr *InteractiveCharmRenderer) renderButtonToRect(vnode *dom.Node, width, height int) Rectangle {
	text := cr.extractRenderedText(vnode)
	rendered := cr.styles.Button.Render(text)
	return NewRectangle(rendered)
}

// renderInputToRect renders an input element to a Rectangle
func (cr *InteractiveCharmRenderer) renderInputToRect(vnode *dom.Node, width, height int) Rectangle {
	props := dom.ExtractProps[dom.InputProps](vnode.Props)

	placeholder := "Enter text..."
	value := ""
	inputType := "text"
	if props.InputType != "" {
		inputType = props.InputType
	}

	if props.Placeholder != "" {
		placeholder = props.Placeholder
	}
	value = props.Value

	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(value)
	ti.CharLimit = 156

	if props.Width > 0 {
		ti.Width = props.Width
	} else {
		// Use the provided width parameter
		if width > 10 {
			ti.Width = width - 10
			if ti.Width < 20 {
				ti.Width = 20
			}
		} else {
			ti.Width = 50
		}
	}

	ti.SetCursor(props.CursorPosition)

	if inputType == "password" {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}

	ti.PromptStyle = cr.styles.Prompt
	ti.TextStyle = cr.styles.InputText
	ti.PlaceholderStyle = cr.styles.Text.Foreground(lipgloss.Color("#626262")).Italic(true)

	if props.Focused {
		ti.Focus()
	} else {
		ti.Blur()
	}

	rendered := cr.styles.Input.Render(ti.View())
	return NewRectangle(rendered)
}

// renderListToRect renders a ul element to a Rectangle
// Children are stacked vertically: width stays same, height decreases
func (cr *InteractiveCharmRenderer) renderListToRect(vnode *dom.Node, width, height int) Rectangle {
	var childRects []Rectangle
	remainingHeight := height

	for _, child := range vnode.Children {
		if remainingHeight <= 0 {
			break
		}
		// Handle FixedSpacer specially for vertical layout
		if child.Type == dom.ElementTypeFixedSpacer {
			childRect := cr.renderFixedSpacerForVertical(child)
			if childRect.Height > 0 {
				childRects = append(childRects, childRect)
				remainingHeight -= childRect.Height
			}
		} else {
			childRect := cr.renderNodeToRect(child, width, remainingHeight)
			if childRect.Height > 0 {
				childRects = append(childRects, childRect)
				remainingHeight -= childRect.Height
			}
		}
	}

	if len(childRects) == 0 {
		return Rectangle{Width: 0, Height: 0, Lines: []string{}}
	}

	return stackVertically(childRects)
}

// renderListItemToRect renders a li element to a Rectangle
func (cr *InteractiveCharmRenderer) renderListItemToRect(vnode *dom.Node, width, height int) Rectangle {
	props := dom.ExtractProps[dom.ListItemProps](vnode.Props)
	text := cr.extractRenderedText(vnode)

	var prefix string
	if props.ItemPrefix != nil {
		prefix = *props.ItemPrefix
	} else {
		prefix = "• "
		if props.Selected {
			prefix = "> "
		}
	}
	renderedText := prefix + text

	var rendered string
	if props.Selected {
		rendered = cr.styles.CompactSuccess.Render(renderedText)
	} else {
		rendered = cr.styles.CompactText.Render(renderedText)
	}
	return NewRectangle(rendered)
}

// renderBrToRect renders a br element to a Rectangle
func (cr *InteractiveCharmRenderer) renderBrToRect(vnode *dom.Node, width, height int) Rectangle {
	return Rectangle{Width: 0, Height: 1, Lines: []string{""}}
}

// renderSpacerToRect renders a spacer element to a Rectangle
func (cr *InteractiveCharmRenderer) renderSpacerToRect(vnode *dom.Node, width, height int) Rectangle {
	props := dom.ExtractProps[dom.SpacerProps](vnode.Props)
	minSize := props.MinSize
	if minSize <= 0 {
		minSize = 1
	}
	return Rectangle{Width: minSize, Height: 1, Lines: []string{strings.Repeat(" ", minSize)}}
}

// renderFixedSpacerToRect renders a fixed spacer element to a Rectangle
// This is called when the parent context is unknown (fallback)
func (cr *InteractiveCharmRenderer) renderFixedSpacerToRect(vnode *dom.Node, width, height int) Rectangle {
	// Fallback: render as horizontal spacer
	return cr.renderFixedSpacerForHorizontal(vnode)
}

// renderFixedSpacerForHorizontal renders a fixed spacer for horizontal layout (HDiv)
// Returns: width = space, height = 1
func (cr *InteractiveCharmRenderer) renderFixedSpacerForHorizontal(vnode *dom.Node) Rectangle {
	props := dom.ExtractProps[dom.FixedSpacerProps](vnode.Props)
	space := props.Space
	if space <= 0 {
		space = 1
	}
	return Rectangle{Width: space, Height: 1, Lines: []string{strings.Repeat(" ", space)}}
}

// renderFixedSpacerForHorizontal renders a fixed spacer for horizontal layout (HDiv)
// Returns: width = space, height = 1
func (cr *InteractiveCharmRenderer) renderSpacerForHorizontal(vnode *dom.Node) Rectangle {
	props := dom.ExtractProps[dom.SpacerProps](vnode.Props)
	minSize := props.MinSize
	if minSize <= 0 {
		minSize = 1
	}
	return Rectangle{Width: minSize, Height: 1, Lines: []string{strings.Repeat(" ", minSize)}}
}

// renderFixedSpacerForVertical renders a fixed spacer for vertical layout (Div, Fragment, Ul)
// Returns: width = 0, height = space
func (cr *InteractiveCharmRenderer) renderFixedSpacerForVertical(vnode *dom.Node) Rectangle {
	props := dom.ExtractProps[dom.FixedSpacerProps](vnode.Props)
	space := props.Space
	if space <= 0 {
		space = 1
	}
	lines := make([]string, space)
	for i := 0; i < space; i++ {
		lines[i] = ""
	}
	return Rectangle{Width: 0, Height: space, Lines: lines}
}

// renderFragmentToRect renders a fragment to a Rectangle
// Children are stacked vertically: width stays same, height decreases
func (cr *InteractiveCharmRenderer) renderFragmentToRect(vnode *dom.Node, width, height int) Rectangle {
	var childRects []Rectangle
	remainingHeight := height

	for _, child := range vnode.Children {
		if remainingHeight <= 0 {
			break
		}
		// Handle FixedSpacer specially for vertical layout
		if child.Type == dom.ElementTypeFixedSpacer {
			childRect := cr.renderFixedSpacerForVertical(child)
			if childRect.Height > 0 {
				childRects = append(childRects, childRect)
				remainingHeight -= childRect.Height
			}
		} else {
			childRect := cr.renderNodeToRect(child, width, remainingHeight)
			if childRect.Height > 0 {
				childRects = append(childRects, childRect)
				remainingHeight -= childRect.Height
			}
		}
	}

	if len(childRects) == 0 {
		return Rectangle{Width: 0, Height: 0, Lines: []string{}}
	}

	return stackVertically(childRects)
}

// renderContainerToRect renders a container div to a Rectangle
// Vertical layout: width stays same for all children, height decreases as children are rendered
func (cr *InteractiveCharmRenderer) renderContainerToRect(vnode *dom.Node, width, height int) Rectangle {
	var childRects []Rectangle
	remainingHeight := height

	for _, child := range vnode.Children {
		if remainingHeight <= 0 {
			break
		}
		// Handle FixedSpacer specially for vertical layout
		if child.Type == dom.ElementTypeFixedSpacer {
			childRect := cr.renderFixedSpacerForVertical(child)
			if childRect.Height > 0 {
				childRects = append(childRects, childRect)
				remainingHeight -= childRect.Height
			}
		} else {
			childRect := cr.renderNodeToRect(child, width, remainingHeight)
			if childRect.Height > 0 {
				childRects = append(childRects, childRect)
				remainingHeight -= childRect.Height
			}
		}
	}

	var contentRect Rectangle
	if len(childRects) == 0 {
		contentRect = Rectangle{Width: 0, Height: 0, Lines: []string{}}
	} else {
		contentRect = stackVertically(childRects)
	}

	// Apply style (border, padding, etc.)
	style := cr.getNodeStyle(vnode)
	rendered := style.Render(contentRect.String())
	return NewRectangle(rendered)
}

// renderHDivToRect renders an HDiv (horizontal layout) to a Rectangle
// Horizontal layout: height stays same for all children, width decreases as children are rendered
func (cr *InteractiveCharmRenderer) renderHDivToRect(vnode *dom.Node, width, height int) Rectangle {
	var childRects []Rectangle
	remainingWidth := width

	var spacerRectIndexes []int
	for _, child := range vnode.Children {
		if remainingWidth <= 0 {
			break
		}
		// Handle FixedSpacer specially for horizontal layout
		if child.Type == dom.ElementTypeFixedSpacer {
			childRect := cr.renderFixedSpacerForHorizontal(child)
			if childRect.Width > 0 {
				childRects = append(childRects, childRect)
				remainingWidth -= childRect.Width
			}
		} else if child.Type == dom.ElementTypeSpacer {
			spacerRect := Rectangle{
				Height: height,
				Lines:  make([]string, height),
			}
			idx := len(childRects)
			spacerRectIndexes = append(spacerRectIndexes, idx)
			childRects = append(childRects, spacerRect)
		} else {
			childRect := cr.renderNodeToRect(child, remainingWidth, height)
			if childRect.Width > 0 {
				childRects = append(childRects, childRect)
				remainingWidth -= childRect.Width
			}
		}
	}

	if len(childRects) == 0 {
		return Rectangle{Width: 0, Height: 0, Lines: []string{}}
	}

	if remainingWidth > 0 && len(spacerRectIndexes) > 0 {
		n := len(spacerRectIndexes)
		avgWidth := remainingWidth / n
		lastOneWidth := remainingWidth % n

		log.Logf("remainingWidth: %d, spacerRectIndexes: %v, avgWidth: %d, lastOneWidth: %d", remainingWidth, spacerRectIndexes, avgWidth, lastOneWidth)
		for i := 0; i < n-1; i++ {
			idx := spacerRectIndexes[i]
			childRects[idx].Width = avgWidth
		}
		idx := spacerRectIndexes[n-1]
		if lastOneWidth == 0 {
			log.Logf("lastOneWidth is 0, setting width to avgWidth: %d", avgWidth)
			childRects[idx].Width = avgWidth
		} else {
			childRects[idx].Width = lastOneWidth
		}
		for _, idx := range spacerRectIndexes {
			child := vnode.Children[idx]
			props := dom.ExtractProps[dom.SpacerProps](child.Props)
			maxSize := props.MaxSize
			width := childRects[idx].Width
			if maxSize > 0 && width > maxSize {
				width = maxSize
				childRects[idx].Width = width
			}

			lines := childRects[idx].Lines
			for i := 0; i < len(lines); i++ {
				lines[i] = strings.Repeat(" ", width)
			}
		}
	}

	props := dom.ExtractProps[dom.DivProps](vnode.Props)
	return stackHorizontally(childRects, props.Align)
}

// renderZDivToRect renders a ZDiv (z-order overlay) to a Rectangle
// Overlay layout: dimensions stay the same for all children
func (cr *InteractiveCharmRenderer) renderZDivToRect(vnode *dom.Node, width, height int) Rectangle {
	if len(vnode.Children) == 0 {
		return Rectangle{Width: 0, Height: 0, Lines: []string{}}
	}

	// Render each child to a Rectangle with the same dimensions
	childRects := make([]Rectangle, 0, len(vnode.Children))
	for _, child := range vnode.Children {
		// Skip FixedSpacer in ZDiv (no effect)
		if child.Type == dom.ElementTypeFixedSpacer {
			continue
		}
		childRect := cr.renderNodeToRect(child, width, height)
		if childRect.Height > 0 || childRect.Width > 0 {
			childRects = append(childRects, childRect)
		}
	}

	if len(childRects) == 0 {
		return Rectangle{Width: 0, Height: 0, Lines: []string{}}
	}

	// Start with the first child as the base
	result := childRects[0]

	// Overlay each subsequent child onto the result
	for i := 1; i < len(childRects); i++ {
		result = Overlay(result, childRects[i])
	}

	return result
}

// renderDefaultToRect renders unknown elements to a Rectangle
func (cr *InteractiveCharmRenderer) renderDefaultToRect(vnode *dom.Node, width, height int) Rectangle {
	lines := []string{fmt.Sprintf("<%s>", vnode.Type)}
	remainingHeight := height - 2 // Reserve space for opening and closing tags

	for _, child := range vnode.Children {
		if remainingHeight <= 0 {
			break
		}
		childRect := cr.renderNodeToRect(child, width, remainingHeight)
		lines = append(lines, childRect.Lines...)
		remainingHeight -= childRect.Height
	}
	lines = append(lines, fmt.Sprintf("</%s>", vnode.Type))

	content := strings.Join(lines, "\n")
	return NewRectangle(content)
}

// stackVertically stacks rectangles vertically
func stackVertically(rects []Rectangle) Rectangle {
	if len(rects) == 0 {
		return Rectangle{Width: 0, Height: 0, Lines: []string{}}
	}

	// Calculate total height and max width
	totalHeight := 0
	maxWidth := 0
	for _, rect := range rects {
		totalHeight += rect.Height
		if rect.Width > maxWidth {
			maxWidth = rect.Width
		}
	}

	// Combine all lines
	var allLines []string
	for _, rect := range rects {
		allLines = append(allLines, rect.Lines...)
	}

	return Rectangle{
		Width:  maxWidth,
		Height: totalHeight,
		Lines:  allLines,
	}
}

// stackHorizontally stacks rectangles horizontally
func stackHorizontally(rects []Rectangle, align dom.Align) Rectangle {
	if len(rects) == 0 {
		return Rectangle{Width: 0, Height: 0, Lines: []string{}}
	}

	// Calculate total width and max height
	totalWidth := 0
	maxHeight := 0
	for _, rect := range rects {
		totalWidth += rect.Width
		if rect.Height > maxHeight {
			maxHeight = rect.Height
		}
	}

	// Pad all rectangles to the same height with proper alignment
	paddedRects := make([]Rectangle, len(rects))
	for i, rect := range rects {
		if rect.Height < maxHeight {
			paddedRects[i] = padRectangleVertically(rect, maxHeight, align)
		} else {
			paddedRects[i] = rect
		}
	}

	// Combine lines horizontally
	resultLines := make([]string, maxHeight)
	for i := 0; i < maxHeight; i++ {
		var lineParts []string
		for _, rect := range paddedRects {
			lineParts = append(lineParts, rect.Lines[i])
		}
		resultLines[i] = strings.Join(lineParts, "")
	}

	return Rectangle{
		Width:  totalWidth,
		Height: maxHeight,
		Lines:  resultLines,
	}
}

// padRectangleVertically pads a rectangle to the target height with proper vertical alignment
func padRectangleVertically(rect Rectangle, targetHeight int, align dom.Align) Rectangle {
	if rect.Height >= targetHeight {
		return rect
	}

	paddingNeeded := targetHeight - rect.Height
	paddedLines := make([]string, targetHeight)
	emptyLine := strings.Repeat(" ", rect.Width)

	switch align {
	case dom.AlignTop, "": // Default to top alignment
		// Content at top, padding at bottom
		copy(paddedLines, rect.Lines)
		for j := rect.Height; j < targetHeight; j++ {
			paddedLines[j] = emptyLine
		}

	case dom.AlignBottom:
		// Padding at top, content at bottom
		for j := 0; j < paddingNeeded; j++ {
			paddedLines[j] = emptyLine
		}
		copy(paddedLines[paddingNeeded:], rect.Lines)

	case dom.AlignCenter:
		// Padding distributed top and bottom
		topPadding := paddingNeeded / 2
		for j := 0; j < topPadding; j++ {
			paddedLines[j] = emptyLine
		}
		copy(paddedLines[topPadding:], rect.Lines)
		for j := topPadding + rect.Height; j < targetHeight; j++ {
			paddedLines[j] = emptyLine
		}

	default:
		// Unknown alignment, default to top
		copy(paddedLines, rect.Lines)
		for j := rect.Height; j < targetHeight; j++ {
			paddedLines[j] = emptyLine
		}
	}

	return Rectangle{
		Width:  rect.Width,
		Height: targetHeight,
		Lines:  paddedLines,
	}
}

// RenderNodeToRect is a convenience function that creates a renderer,
// renders a node, and returns the result as a Rectangle
// Uses default dimensions of 80x24 (standard terminal size)
func RenderNodeToRect(vnode *dom.Node) Rectangle {
	renderer := NewInteractiveCharmRenderer()
	return renderer.RenderToRect(vnode, 80, 24)
}
