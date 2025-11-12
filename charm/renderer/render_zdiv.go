package renderer

import (
	"github.com/xhd2015/go-dom-tui/dom"
)

// renderZDiv renders children stacked in z-order (like SwiftUI ZStack)
// Children are rendered on top of each other, with later children appearing above earlier ones
// Uses Rectangle-based overlay system for proper shadowing
func (cr *InteractiveCharmRenderer) renderZDiv(vnode *dom.Node, depth int) {
	if len(vnode.Children) == 0 {
		cr.updateRenderState(vnode.Type, true)
		return
	}

	// Get dimensions from window or use defaults
	width, height := 80, 24 // Default terminal size
	if vnode.Window != nil {
		width, height = vnode.Window.Get()
	}

	// Use the new RenderToRect method for proper rectangle-based rendering
	childRenderer := &InteractiveCharmRenderer{styles: cr.styles}
	result := childRenderer.RenderToRect(vnode, width, height)

	// Convert the Rectangle to string and append to output
	cr.output += result.String() + "\n"
	cr.updateRenderState(vnode.Type, true)
}
