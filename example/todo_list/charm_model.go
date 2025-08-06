package main

import (
	"os"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xhd2015/go-dom-tui/charm"
	"github.com/xhd2015/go-dom-tui/charm/renderer"
	"github.com/xhd2015/go-dom-tui/log"
	"github.com/xhd2015/go-dom-tui/react"
)

// Model represents the Bubble Tea model for the todo app
type Model struct {
	app *charm.CharmApp[AppState]

	renderer  *renderer.InteractiveCharmRenderer
	root      *react.Root
	cleanedUp bool     // Whether cleanup has been performed
	debugFile *os.File // Debug log file
	program   *tea.Program
}

// NewModel creates a new Bubble Tea model
func NewModel(debugFilePath string) *Model {
	renderer := renderer.NewInteractiveCharmRenderer()
	root := react.CreateRoot(renderer)

	var debugFile *os.File
	if debugFilePath != "" {
		var err error
		debugFile, err = os.OpenFile(debugFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			// Don't fail the program, just disable debug logging
			debugFile = nil
		} else {
			// Set the debug log file for the dom_based package
			log.SetLogger(log.NewFileLogger(debugFile))
		}
	}

	m := &Model{
		renderer:  renderer,
		root:      root,
		debugFile: debugFile,
	}

	appState := AppState{
		InputFocused:      true, // Start with input focused
		SelectedTodoIndex: -1,
		Refresh: func() {
			m.program.Send(cursor.Blink())
		},
	}
	m.app = charm.NewCharmApp(&appState, App)

	return m
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model using DOM-like event dispatching
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.app.Update(msg)
	if m.app.State.Quitting {
		m.cleanup()
		return m, tea.Quit
	}

	return m, nil
}

// cleanup cleans up resources before exiting
func (m *Model) cleanup() {
	// Only cleanup once
	if m.cleanedUp {
		return
	}

	log.Logf("Cleaning up resources...")
	m.cleanedUp = true

	// Close and remove debug log file if it exists
	if m.debugFile != nil {
		debugFilePath := m.debugFile.Name()
		log.Logf("Closing debug log file: %s", debugFilePath)

		// Close the file
		m.debugFile.Close()
		m.debugFile = nil
	}
}

// View renders the current view
func (m *Model) View() string {
	return m.app.Render()
}

// recreateDOMReact recreates the DOM structure (extracted for clarity)
// func (m *Model) recreateDOMReact() {
// 	// Create updated app component using the state field directly
// 	app := App(AppProps{
// 		State:                     m.appData,
// 		OnInputKeyDown:            m.onInputKeyDown,
// 		OnInputChange:             m.OnInputChange,
// 		OnInputFocus:              m.OnInputFocus,
// 		OnInputBlur:               m.OnInputBlur,
// 		OnChangeSelectedTodoIndex: m.OnChangeSelectedTodoIndex,
// 		OnInputCursorMove:         m.OnInputCursorMove,
// 	})
// 	m.root.Render(app, nil)

// 	// Get the VNode from the reconciler instead of creating a duplicate
// 	newVNode := m.root.Reconciler.GetCurrentRootVNode()
// 	if newVNode == nil {
// 		log.Logf("No reconciler VNode available, creating a new one")
// 		newVNode = app
// 	}

// 	// Recreate DOM with updated content - focus state is driven by app state
// 	m.dom = dom.NewDOM(newVNode)
// }
