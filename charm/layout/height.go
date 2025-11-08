package layout

import "github.com/xhd2015/go-dom-tui/dom"

// GetNodeRenderedHeight calculates the approximate rendered height of a DOM node
// This is a simplified calculation that counts lines
func GetNodeRenderedHeight(node *dom.Node) int {
	if node == nil {
		return 0
	}

	// For text nodes, count as 1 line (unless it contains newlines)
	if node.Type == dom.ElementTypeText {
		// Text nodes are typically single line
		return 1
	}

	// Br creates a line break
	if node.Type == dom.ElementTypeBr {
		return 1
	}

	// For fragments, sum up children heights (inline layout)
	if node.Type == dom.ElementTypeFragment {
		// Fragment doesn't add height itself, just renders children inline
		// We take the max height of children since they're on the same line
		maxHeight := 0
		for _, child := range node.Children {
			height := GetNodeRenderedHeight(child)
			if height > maxHeight {
				maxHeight = height
			}
		}
		return maxHeight
	}

	// For divs and other block elements
	if node.Type == dom.ElementTypeDiv || node.Type == dom.ElementTypeSpan ||
		node.Type == dom.ElementTypeH1 || node.Type == dom.ElementTypeH2 ||
		node.Type == dom.ElementTypeP || node.Type == dom.ElementTypeUl ||
		node.Type == dom.ElementTypeLi || node.Type == dom.ElementTypeButton {

		// Check if all children are inline elements (Text nodes, Span, etc.)
		// If so, they render on the same line
		allInline := true
		for _, child := range node.Children {
			if child.Type == dom.ElementTypeDiv || child.Type == dom.ElementTypeBr ||
				child.Type == dom.ElementTypeH1 || child.Type == dom.ElementTypeH2 ||
				child.Type == dom.ElementTypeP || child.Type == dom.ElementTypeUl ||
				child.Type == dom.ElementTypeLi {
				allInline = false
				break
			}
		}

		// If all children are inline, the div is a single line
		if allInline && len(node.Children) > 0 {
			return 1
		}

		// Otherwise, sum up children heights (block elements stack vertically)
		totalHeight := 0
		for _, child := range node.Children {
			totalHeight += GetNodeRenderedHeight(child)
		}
		// Block elements themselves take at least 1 line if they have no children
		if totalHeight == 0 {
			return 1
		}
		return totalHeight
	}

	// Default: assume 1 line
	return 1
}

// GetTotalNodesHeight returns the total rendered height of a list of nodes
func GetTotalNodesHeight(nodes []*dom.Node) int {
	totalHeight := 0
	for _, node := range nodes {
		totalHeight += GetNodeRenderedHeight(node)
	}
	return totalHeight
}
