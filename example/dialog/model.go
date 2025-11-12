package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xhd2015/go-dom-tui/charm/renderer"
	"github.com/xhd2015/go-dom-tui/dom"
)

type Model struct {
	showDialog bool
	items      []string
	selected   int
	width      int
	height     int
}

func NewModel() *Model {
	return &Model{
		items: []string{
			"Item 1 - Press Enter to show dialog",
			"Item 2 - Use arrow keys to navigate",
			"Item 3 - Press Esc to close dialog",
			"Item 4 - Press q to quit",
		},
		selected: 0,
		width:    80, // Default terminal width
		height:   24, // Default terminal height
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.selected > 0 {
				m.selected--
			}
		case "down":
			if m.selected < len(m.items)-1 {
				m.selected++
			}
		case "enter":
			m.showDialog = !m.showDialog
		case "esc":
			m.showDialog = false
		}
	}
	return m, nil
}

func (m *Model) View() string {
	return m.renderView(false)
}

func (m *Model) ViewStripped() string {
	// Reuse the rendering logic with color stripping enabled
	return m.renderView(true)
}

func (m *Model) renderView(stripColor bool) string {
	// Create the list
	listItems := make([]*dom.Node, 0, len(m.items))
	for i, item := range m.items {
		text := item
		if i == m.selected {
			text = "> " + text
		} else {
			text = "  " + text
		}
		listItems = append(listItems, dom.Div(dom.DivProps{}, dom.Text(text)))
	}

	list := dom.Div(dom.DivProps{},
		dom.Text("Dialog Example - List with Overlay"),
		dom.Div(dom.DivProps{}, dom.Text("")), // Empty line
		dom.Fragment(listItems...),
	)

	// If dialog is shown, use ZDiv to overlay it
	var root *dom.Node
	if m.showDialog {
		dialog := dom.Div(dom.DivProps{},
			dom.Text("+---------------------------+"),
			dom.Text("|   Dialog Overlay Demo    |"),
			dom.Text("|                           |"),
			dom.Text("| This is an overlay dialog |"),
			dom.Text("| Press Esc to close        |"),
			dom.Text("+---------------------------+"),
		)

		root = dom.ZDiv(dom.DivProps{}, list, dialog)
	} else {
		root = list
	}

	// Use rectangle-based rendering
	r := renderer.NewInteractiveCharmRenderer()
	rect := r.RenderToRect(root, m.width, m.height)
	output := rect.String()

	if stripColor {
		return renderer.StripColor(output)
	}
	return output
}
