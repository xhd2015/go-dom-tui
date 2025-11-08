package layout

import "github.com/xhd2015/go-dom-tui/dom"

// GetNodeRenderedWidth calculates the approximate rendered width of a DOM node
// This is a simplified calculation that works for text nodes
func GetNodeRenderedWidth(node *dom.Node) int {
	if node == nil {
		return 0
	}

	// For text nodes, count the runes
	if node.Type == dom.ElementTypeText {
		return len([]rune(node.Text))
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
