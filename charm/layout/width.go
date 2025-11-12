package layout

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/xhd2015/go-dom-tui/dom"
)

// GetNodeRenderedWidth calculates the approximate rendered width of a DOM node
// This calculation accounts for ANSI escape sequences (colors, styles) and control characters
func GetNodeRenderedWidth(node *dom.Node) int {
	if node == nil {
		return 0
	}

	// For text nodes, use lipgloss.Width which correctly handles ANSI escape codes
	if node.Type == dom.ElementTypeText {
		return lipgloss.Width(node.Text)
	}

	// For hdiv (horizontal div), sum up children widths (placed horizontally)
	if node.Type == dom.ElementTypeHDiv {
		totalWidth := 0
		for _, child := range node.Children {
			totalWidth += GetNodeRenderedWidth(child)
		}
		return totalWidth
	}

	// For fragments and divs, sum up children widths
	if node.Type == dom.ElementTypeFragment || node.Type == dom.ElementTypeDiv {
		totalWidth := 0
		for _, child := range node.Children {
			totalWidth += GetNodeRenderedWidth(child)
		}
		return totalWidth
	}

	return 0
}

// GetMaxNodeWidth returns the maximum rendered width among a list of nodes
func GetMaxNodeWidth(nodes []*dom.Node) int {
	maxWidth := 0
	for _, node := range nodes {
		width := GetNodeRenderedWidth(node)
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}
