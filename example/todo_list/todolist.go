package main

import (
	"fmt"
	"strings"

	"github.com/xhd2015/go-dom-tui/dom"
)

// TodoItem represents a todo item
type TodoItem struct {
	Title       string
	Description string
	Priority    string
	Completed   bool
}

// TodoListDisplayProps represents props for TodoListDisplay component
type TodoListDisplayProps struct {
	Todos                     []TodoItem
	SelectedIndex             int
	OnChangeSelectedTodoIndex func(int)
}

// TodoListDisplay component
func TodoListDisplay(props TodoListDisplayProps) *dom.Node {
	todos := props.Todos
	selectedIndex := props.SelectedIndex
	onChangeSelectedTodoIndex := props.OnChangeSelectedTodoIndex

	todoElements := make([]*dom.Node, len(todos))
	for i, todo := range todos {
		status := "⭕"
		if todo.Completed {
			status = "✅"
		}

		// Don't add focus indicator here - it will be handled in the renderer based on actual DOM focus
		todoText := fmt.Sprintf("%s [%s] %s", status, strings.ToUpper(todo.Priority), todo.Title)
		if todo.Description != "" {
			todoText += " - " + todo.Description
		}

		// Create focusable li element with callbacks
		todoElements[i] = dom.Li(dom.ListItemProps{
			Text:    todoText,
			Index:   i,
			Focused: i == selectedIndex,
			OnFocus: func() {
				if onChangeSelectedTodoIndex != nil {
					onChangeSelectedTodoIndex(i)
				}
			},
			OnBlur: func() {
				if onChangeSelectedTodoIndex != nil {
					onChangeSelectedTodoIndex(-1)
				}
			},
			Focusable: dom.Focusable(true),
		})
	}

	return dom.Div(dom.DivProps{
		Style: dom.Style{}, // No border by default
	},
		dom.H2(dom.DivProps{}, dom.Text("📋 Current Tasks")),
		dom.Ul(dom.DivProps{}, todoElements...),
	)
}
