package layout

import (
	"strings"

	"github.com/xhd2015/go-dom-tui/dom"
	"github.com/xhd2015/go-dom-tui/styles"
)

// mergeAlignBottom is the base function that merges two DOM node arrays side by side, aligned from the bottom
// Returns a 2D array where each element is [nodeA, spacer, nodeB] (or subsets if nodes are missing)
// The caller can wrap these nodes in Div, Fragment, or other containers as needed
func mergeAlignBottom(a []*dom.Node, b []*dom.Node, spaceWidth int, useWidthAlignment bool) [][]*dom.Node {
	if len(a) == 0 && len(b) == 0 {
		return [][]*dom.Node{}
	}
	if len(a) == 0 {
		// Return b wrapped in individual arrays
		result := make([][]*dom.Node, len(b))
		for i, node := range b {
			result[i] = []*dom.Node{node}
		}
		return result
	}
	if len(b) == 0 {
		// Return a wrapped in individual arrays
		result := make([][]*dom.Node, len(a))
		for i, node := range a {
			result[i] = []*dom.Node{node}
		}
		return result
	}

	// Calculate the maximum width of elements in array 'a' if width alignment is enabled
	var maxWidthA int
	if useWidthAlignment {
		maxWidthA = GetMaxNodeWidth(a)
	}

	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	result := make([][]*dom.Node, maxLen)

	// Calculate padding needed for each array
	aPadding := maxLen - len(a)
	bPadding := maxLen - len(b)

	for i := 0; i < maxLen; i++ {
		var nodeA, nodeB *dom.Node
		var hasA, hasB bool

		// Get node from array a (with top padding)
		if i < aPadding {
			hasA = false
		} else {
			nodeA = a[i-aPadding]
			hasA = true
		}

		// Get node from array b (with top padding)
		if i < bPadding {
			hasB = false
		} else {
			nodeB = b[i-bPadding]
			hasB = true
		}

		// Build the node array for this row
		var rowNodes []*dom.Node
		if hasA && hasB {
			// Both nodes exist - calculate padding and add spacer
			var paddingNeeded int
			if useWidthAlignment {
				nodeWidth := GetNodeRenderedWidth(nodeA)
				paddingNeeded = maxWidthA - nodeWidth + spaceWidth
				if paddingNeeded < spaceWidth {
					paddingNeeded = spaceWidth
				}
			} else {
				paddingNeeded = spaceWidth
			}
			spacer := dom.Text(strings.Repeat(" ", paddingNeeded), styles.Style{})
			rowNodes = []*dom.Node{nodeA, spacer, nodeB}
		} else if hasA {
			// Only A exists
			rowNodes = []*dom.Node{nodeA}
		} else if hasB {
			// Only B exists - pad to align with max width if needed
			if useWidthAlignment {
				paddingNeeded := maxWidthA + spaceWidth
				spacer := dom.Text(strings.Repeat(" ", paddingNeeded), styles.Style{})
				rowNodes = []*dom.Node{spacer, nodeB}
			} else {
				rowNodes = []*dom.Node{nodeB}
			}
		}

		result[i] = rowNodes
	}

	return result
}

// MergeAlignBottom merges two DOM node arrays side by side, aligned from the bottom
// Handles cases where len(a) != len(b) by padding the shorter array at the top
// Uses Div to wrap elements (creates block-level layout)
func MergeAlignBottom(a []*dom.Node, b []*dom.Node, spaceWidth int) []*dom.Node {
	rows := mergeAlignBottom(a, b, spaceWidth, false)
	result := make([]*dom.Node, len(rows))
	for i, rowNodes := range rows {
		if len(rowNodes) == 0 {
			result[i] = dom.Div(dom.DivProps{})
		} else {
			result[i] = dom.Div(dom.DivProps{}, rowNodes...)
		}
	}
	return result
}

// MergeAlignBottomFragment merges two DOM node arrays side by side, aligned from the bottom
// Handles cases where len(a) != len(b) by padding the shorter array at the top
// Uses Fragment to wrap elements (creates inline layout - elements appear on same line)
// Pads elements in array 'a' to the max width for proper alignment
func MergeAlignBottomFragment(a []*dom.Node, b []*dom.Node, spaceWidth int) []*dom.Node {
	rows := mergeAlignBottom(a, b, spaceWidth, true)
	result := make([]*dom.Node, len(rows))
	for i, rowNodes := range rows {
		if len(rowNodes) == 0 {
			result[i] = dom.Fragment()
		} else {
			result[i] = dom.Fragment(rowNodes...)
		}
	}
	return result
}
