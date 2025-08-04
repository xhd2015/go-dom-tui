package charm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xhd2015/go-dom-tui/charm/renderer"
	"github.com/xhd2015/go-dom-tui/dom"
	"github.com/xhd2015/go-dom-tui/log"
)

type CharmApp[T any] struct {
	State *T
	Root  func(state *T, window *dom.Window) *dom.Node

	width  int
	height int

	renderer *renderer.InteractiveCharmRenderer
	dom      *dom.DOM // DOM tree with event handling
}

func NewCharmApp[T any](state *T, app func(state *T, window *dom.Window) *dom.Node) *CharmApp[T] {
	return &CharmApp[T]{
		State:    state,
		Root:     app,
		renderer: renderer.NewInteractiveCharmRenderer(),
	}
}

func (c *CharmApp[T]) Update(msg tea.Msg) interface{} {
	log.Logf("Update: %T", msg)
	var res interface{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if c.dom != nil {
			res = c.dom.DispatchEvent("keydown", msg.String(), c)
		}
	case tea.WindowSizeMsg:
		log.Logf("window size: %d x %d", msg.Width, msg.Height)

		// Dispatch window resize event to DOM
		if c.dom != nil {
			// Update window state through DOM root node
			if c.dom.Root != nil && c.dom.Root.Window != nil {
				c.width = msg.Width
				c.height = msg.Height
				c.dom.Root.Window.Update(c.width, c.height)
			}

			event := &dom.WindowResizeEvent{
				Width:  msg.Width,
				Height: msg.Height,
			}
			res = c.dom.DispatchWindowEvent("resize", event)
		}
	}

	return res
}

// View renders the current view
func (c *CharmApp[T]) Render() string {
	window := &dom.Window{
		Width:  c.width,
		Height: c.height,
	}
	c.dom = dom.NewDOM(c.Root(c.State, window), window)

	// Render the output
	return c.renderer.Render(c.dom.Root)
}
