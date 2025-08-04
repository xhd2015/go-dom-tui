package charm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xhd2015/go-dom-tui/charm/renderer"
	"github.com/xhd2015/go-dom-tui/dom"
)

type CharmApp[T any] struct {
	State *T
	Root  func(state *T) *dom.Node

	renderer *renderer.InteractiveCharmRenderer
	dom      *dom.DOM // DOM tree with event handling
}

func NewCharmApp[T any](state *T, app func(state *T) *dom.Node) *CharmApp[T] {
	return &CharmApp[T]{
		State:    state,
		Root:     app,
		renderer: renderer.NewInteractiveCharmRenderer(),
	}
}

func (c *CharmApp[T]) Update(msg tea.Msg) interface{} {
	var res interface{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if c.dom != nil {
			res = c.dom.DispatchEvent("keydown", msg.String(), c)
		}
	}

	return res
}

// View renders the current view
func (c *CharmApp[T]) Render() string {
	c.dom = dom.NewDOM(c.Root(c.State))

	// Render the output
	return c.renderer.Render(c.dom.Root)
}
