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

func (c *CharmApp[T]) Update(msg tea.Msg) {
	log.Logf("Update: %T", msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Logf("Key Msg %v: alt=%v, paste=%v, len(runes)=%v", msg.Type, msg.Alt, msg.Paste, len(msg.Runes))
		if c.dom != nil {
			var keyType dom.KeyType
			switch msg.Type {
			case tea.KeyUp:
				keyType = dom.KeyTypeUp
			case tea.KeyDown:
				keyType = dom.KeyTypeDown
			case tea.KeyLeft:
				keyType = dom.KeyTypeLeft
			case tea.KeyRight:
				keyType = dom.KeyTypeRight
			case tea.KeyEnter:
				keyType = dom.KeyTypeEnter
			case tea.KeyBackspace:
				keyType = dom.KeyTypeBackspace
			case tea.KeyDelete:
				keyType = dom.KeyTypeDelete
			case tea.KeyEscape:
				keyType = dom.KeyTypeEsc
			case tea.KeySpace:
				keyType = dom.KeyTypeSpace
			case tea.KeyTab:
				keyType = dom.KeyTypeTab
			case tea.KeyCtrlC:
				keyType = dom.KeyTypeCtrlC
			case tea.KeyCtrlV:
				keyType = dom.KeyTypeCtrlV
			case tea.KeyCtrlX:
				keyType = dom.KeyTypeCtrlX
			case tea.KeyCtrlW:
				keyType = dom.KeyTypeCtrlW
			case tea.KeyCtrlA:
				keyType = dom.KeyTypeCtrlA
			case tea.KeyCtrlE:
				keyType = dom.KeyTypeCtrlE
			case tea.KeyCtrlK:
				keyType = dom.KeyTypeCtrlK
			}
			c.dom.DispatchKeyDownEvent(&dom.KeydownEvent{
				KeyType: keyType,
				Runes:   msg.Runes,
				Alt:     msg.Alt,
				Paste:   msg.Paste,
			})
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
			c.dom.DispatchWindowEvent(dom.EventTypeResize, event)
		}
	}
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
