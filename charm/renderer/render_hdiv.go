package renderer

import (
	"strings"

	"github.com/xhd2015/go-dom-tui/dom"
)

// renderHDiv renders a horizontal div that places children inline
func (cr *InteractiveCharmRenderer) renderHDiv(vnode *dom.Node, depth int) {
	// HDiv renders children horizontally (side by side)
	// Each child is rendered independently, then placed left to right

	if len(vnode.Children) == 0 {
		cr.updateRenderState(vnode.Type, true)
		return
	}

	// Extract alignment from props (default to top)
	props := dom.ExtractProps[dom.DivProps](vnode.Props)
	align := props.Align
	if align == "" {
		align = dom.AlignTop
	}

	// Render each child independently and collect their output lines
	type ChildLines struct {
		lines    []string
		maxWidth int
	}

	childrenLines := make([]ChildLines, 0, len(vnode.Children))
	maxLines := 0

	for _, child := range vnode.Children {
		childRenderer := &InteractiveCharmRenderer{styles: cr.styles}
		childRenderer.renderNode(child, depth+1)

		// Split the rendered output into lines
		lines := strings.Split(childRenderer.output, "\n")
		// Remove trailing empty line if exists
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}

		// Calculate max width of this child
		maxWidth := 0
		for _, line := range lines {
			// Count runes for proper width calculation (handles multi-byte chars)
			width := len([]rune(line))
			if width > maxWidth {
				maxWidth = width
			}
		}

		childrenLines = append(childrenLines, ChildLines{
			lines:    lines,
			maxWidth: maxWidth,
		})

		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}

	// Now compose the output by placing children side by side
	// For each line, concatenate the corresponding line from each child
	var result strings.Builder
	for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
		for childIdx, childLines := range childrenLines {
			var lineContent string

			// Determine which line to show based on alignment
			var childLineIdx int
			if align == dom.AlignBottom {
				// Bottom alignment: pad at the top
				childLineIdx = lineIdx - (maxLines - len(childLines.lines))
				if childLineIdx < 0 {
					// This line is padding (empty)
					lineContent = ""
				} else {
					lineContent = childLines.lines[childLineIdx]
				}
			} else {
				// Top alignment (default): pad at the bottom
				if lineIdx < len(childLines.lines) {
					lineContent = childLines.lines[lineIdx]
				} else {
					// This child doesn't have this line, use empty string
					lineContent = ""
				}
			}

			// Pad the line to match the child's max width (except for the last child)
			if childIdx < len(childrenLines)-1 {
				currentWidth := len([]rune(lineContent))
				padding := childLines.maxWidth - currentWidth
				if padding > 0 {
					lineContent += strings.Repeat(" ", padding)
				}
			}

			result.WriteString(lineContent)
		}

		// Add newline after each line
		result.WriteString("\n")
	}

	cr.output += result.String()
	cr.updateRenderState(vnode.Type, true)
}
