package main

import (
	"fmt"
	"time"

	"github.com/xhd2015/go-dom-tui/dom"
	"github.com/xhd2015/go-dom-tui/log"
)

const (
	CtrlCTimeoutMs = 1000
)

type AppState struct {
	Todos  []TodoItem
	Status string // Current status message

	InputValue          string // Current input field value
	InputFocused        bool   // Whether input field is focused
	InputCursorPosition int    // Cursor position

	SelectedTodoIndex int

	Quitting bool

	LastCtrlCTime time.Time
	Refresh       func()
}

func (m *AppState) OnChangeSelectedTodoIndex(index int) {
	m.SelectedTodoIndex = index
}

func (m *AppState) OnInputChange(value string) {
	m.InputValue = value
}

func (m *AppState) OnInputFocus() {
	m.InputFocused = true
}

func (m *AppState) OnInputBlur() {
	m.InputFocused = false
}

func (m *AppState) OnInputCursorMove(delta int, seek int) {
	newPos := m.InputCursorPosition + delta
	if newPos < 0 {
		newPos = 0
	}
	if newPos > len(m.InputValue)+1 {
		newPos = len(m.InputValue) + 1
	}
	m.InputCursorPosition = newPos
}

func (m *AppState) OnInputKeyDown(event *dom.DOMEvent) {
	log.Logf("Input key down: %s", event.Key)

	switch event.Key {
	case "enter":
		// Handle Enter key for input - get current value from model
		value := m.InputValue
		if value == "" {
			return
		}

		// Add new todo
		newTodo := TodoItem{
			Title:       value,
			Description: "",       // Empty as requested
			Priority:    "medium", // Default priority
			Completed:   false,
		}

		m.Todos = append(m.Todos, newTodo)
		m.Status = fmt.Sprintf("✅ Added: %s", newTodo.Title)

		// Clear the input
		m.InputValue = ""
		m.OnInputChange("")
	}
}

// AppComponent creates the main application component with live textinput
func App(props *AppState) *dom.Node {
	selectedTodoIndex := props.SelectedTodoIndex
	onChangeSelectedTodoIndex := props.OnChangeSelectedTodoIndex

	if props.Quitting {
		return dom.Div(dom.DivProps{}, dom.Text("Thanks for using Quick Todo!"))
	}

	return dom.Div(dom.DivProps{
		Style: "main-app",
		OnKeyDown: func(event *dom.DOMEvent) {
			if event.Key == "ctrl+c" {
				if time.Since(props.LastCtrlCTime) < CtrlCTimeoutMs*time.Millisecond {
					log.Logf("App: quitting by double ctrl-c ")
					props.Quitting = true
				} else {
					go func() {
						time.Sleep(CtrlCTimeoutMs * time.Millisecond)
						if props.Refresh != nil {
							props.Refresh()
						}
					}()
				}
				props.LastCtrlCTime = time.Now()
			}
		},
	},
		// Header
		dom.H1(dom.DivProps{}, dom.Text("🎯 Quick Todo")),

		// Status message (if any)
		func() *dom.Node {
			if props.Status != "" {
				return dom.P(dom.DivProps{}, dom.Text(props.Status))
			}
			return nil
		}(),

		// Quick input section with actual textinput component
		dom.Div(dom.DivProps{
			Style: "quick-input",
		},
			dom.P(dom.DivProps{}, dom.Text("Type a task and press Enter:")),
			// Use direct input element
			dom.Input(dom.InputComponentProps{
				Placeholder:    "What needs to be done?",
				Value:          props.InputValue,
				Focused:        props.InputFocused,
				CursorPosition: props.InputCursorPosition,
				OnCursorMove:   props.OnInputCursorMove,
				OnKeyDown: func(d *dom.DOMEvent) {
					if d.Key == "enter" {
						switch props.InputValue {
						case "quit", "exit", "q":
							log.Logf("App: quitting by input %s", props.InputValue)
							props.Quitting = true
						}
					}
					props.OnInputKeyDown(d)
				},
				OnChange: props.OnInputChange,
				OnFocus:  props.OnInputFocus,
				OnBlur:   props.OnInputBlur,
			}),
		),

		func() *dom.Node {
			if time.Since(props.LastCtrlCTime) < CtrlCTimeoutMs*time.Millisecond {

				// red color
				return dom.P(dom.DivProps{
					Style: "red",
				}, dom.Text("Ctrl-C again to quit."))
			}
			return nil
		}(),

		// Todo list section (only show if there are todos)
		func() *dom.Node {
			if len(props.Todos) > 0 {
				return TodoListDisplay(TodoListDisplayProps{
					Todos:                     props.Todos,
					SelectedIndex:             selectedTodoIndex,
					OnChangeSelectedTodoIndex: onChangeSelectedTodoIndex,
				})
			}
			return dom.Div(dom.DivProps{},
				dom.P(dom.DivProps{}, dom.Text("📝 No tasks yet. Start typing above!")))
		}(),
	)
}
