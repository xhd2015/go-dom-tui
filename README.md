# Go-DOM-TUI: React for Terminal UIs

Build terminal apps like React components.


## 🏃‍♂️ Getting Started

### 1. Install Dependencies
```bash
go mod init your-tui-app
cd your-tui-app
go get github.com/xhd2015/go-dom-tui
```

### 2. Create Your First Component
```go
package main

import (
    "github.com/xhd2015/go-dom-tui/dom"
    tea "github.com/charmbracelet/bubbletea"
)

func TodoApp(props TodoAppProps) *dom.Node {
    return dom.Div(dom.DivProps{},
        dom.H1(dom.DivProps{}, dom.Text("🎯 Quick Todo")),
        
        dom.Input(dom.InputComponentProps{
            Placeholder: "What needs to be done?",
            OnKeyDown:   props.HandleInputKeyDown,  // Add todo on Enter
        }),
        
        TodoList(TodoListProps{
            Todos: []TodoItem{
                {Title: "Learn Go-DOM", Completed: false},
                {Title: "Build awesome TUI", Completed: true},
            },
        }),
    )
}

func main() {
    // Integrate with Bubble Tea
    model := NewModel(TodoApp)
    p := tea.NewProgram(&model, tea.WithAltScreen())
    p.Run()
}
```

### 3. Run the Todo Example
```bash
go run .
```

**Terminal Output:**
```
🎯 Quick Todo

What needs to be done?
┌─────────────────────────┐
│ What needs to be done?  │
└─────────────────────────┘

📋 Current Tasks
⭕ [MEDIUM] Learn Go-DOM
✅ [MEDIUM] Build awesome TUI
```

Check [example/todo_list](example/todo_list/):
```sh
git clone https://github.com/xhd2015/go-dom-tui
cd go-dom-tui/example/todo_list
go run ./
```

## 🎯 The Problem

Traditional TUI libraries like [Bubble Tea](https://github.com/charmbracelet/bubbletea) require imperative programming where you manually manage state, events, and rendering logic. This leads to:

- **Complex state management** scattered across event handlers
- **Tightly coupled components** that are hard to test and reuse
- **Manual DOM-like manipulations** for focus, events, and updates
- **Steep learning curve** for developers coming from web development

Compare traditional Bubble Tea code:
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            if m.focused == inputField {
                // Manual state updates
                m.todos = append(m.todos, m.input)
                m.input = ""
                m.cursor = 0
            }
        case "tab":
            // Manual focus management
            m.focused = (m.focused + 1) % totalFields
        }
    }
    return m, nil
}
```

With Go-DOM-TUI's declarative approach:
```go
func App(props dom.Props) *dom.Node {
    appProps := dom.ExtractProps[AppProps](props)
    
    return dom.Div(dom.DivProps{},
        dom.Input(dom.InputComponentProps{
            Value:     appProps.Data.InputValue,
            OnKeyDown: appProps.OnInputKeyDown,  // Declarative event handling
            OnChange:  appProps.OnInputChange,   // Automatic state binding
        }),
        TodoList(dom.NewStructProps(TodoListProps{
            Todos: appProps.Data.Todos,
        })),
    )
}
```

## 🚀 The Solution

Go-DOM-TUI implements a complete **Virtual DOM system** with **React-like infrastructure** specifically designed for terminal applications:

### 🏗️ Virtual DOM Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Components    │───▶│   Virtual DOM   │───▶│   TUI Output    │
│  (React-like)   │    │   (VNodes)      │    │   (Terminal)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Event System    │    │     Diffing     │    │    Renderer     │
│ (DOM Events)    │    │  (Reconciler)   │    │ (Charm/Lipgloss)│
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### ⚛️ React Infrastructure
- **Declarative Components**: Build UIs by describing what they should look like
- **Event-Driven Architecture**: DOM-style event bubbling and handling
- **State Management**: Props flow down, events flow up
- **Component Composition**: Reusable, testable components
- **Focus Management**: Automatic tab order and keyboard navigation

## 🌟 Key Benefits

### For Web Developers
- **Familiar Patterns**: Uses React-like component structure and props
- **Declarative Syntax**: Describe the UI state, not the mutations
- **Component Reusability**: Build once, use everywhere
- **Event Handling**: Standard DOM event patterns (onClick, onKeyDown, etc.)

### For TUI Development
- **Reduced Complexity**: No manual state synchronization
- **Better Separation of Concerns**: UI logic separate from business logic
- **Easier Testing**: Components are pure functions of their props
- **Maintainable Code**: Clear data flow and component boundaries

## 📋 Real-World Example: Todo Application

The `/example/todo_list` demonstrates a complete todo application built with Go-DOM:

### Component Structure
```go
// Main App Component
func App(props dom.Props) *dom.Node {
    return dom.Div(dom.DivProps{},
        dom.H1(dom.DivProps{}, dom.Text("🎯 Quick Todo")),
        
        // Reusable Input Component
        dom.Input(dom.InputComponentProps{
            Placeholder:    "What needs to be done?",
            Value:          appProps.Data.InputValue,
            OnKeyDown:      appProps.OnInputKeyDown,
            OnChange:       appProps.OnInputChange,
        }),
        
        // Reusable TodoList Component
        TodoListDisplay(dom.NewStructProps(TodoListDisplayProps{
            Todos:         appProps.Data.Todos,
            SelectedIndex: appProps.Data.SelectedTodoIndex,
        })),
    )
}
```

### Event-Driven Updates
```go
func (m *Model) onInputKeyDown(event *dom.DOMEvent) {
    switch event.Key {
    case "enter":
        if value := m.appData.InputValue; value != "" {
            // Clean separation: pure data operations
            newTodo := TodoItem{
                Title:    value,
                Priority: "medium",
                Completed: false,
            }
            m.appData.Todos = append(m.appData.Todos, newTodo)
            m.appData.InputValue = ""
        }
    }
}
```

### Automatic Focus Management
```go
dom.Li(dom.ListItemProps{
    Text:      todoText,
    Focused:   i == selectedIndex,  // Declarative focus state
    Focusable: dom.Focusable(true), // Automatic tab order
    TabIndex:  dom.TabIndex(i + 1),
    OnFocus:   func() { onChangeSelectedTodoIndex(i) },
})
```

## 🆚 Go-DOM-TUI vs Traditional TUI Libraries

| Aspect | Traditional TUI | Go-DOM-TUI |
|--------|----------------|---------|
| **Programming Model** | Imperative (manual updates) | Declarative (describe state) |
| **State Management** | Scattered across handlers | Centralized in props/data |
| **Event Handling** | Manual message routing | DOM-style event bubbling |
| **Focus Management** | Manual tracking | Automatic with tab order |
| **Component Reuse** | Copy-paste patterns | True component composition |
| **Testing** | Integration tests only | Unit test pure functions |
| **Learning Curve** | TUI-specific patterns | Familiar web dev patterns |
| **Code Maintenance** | Complex state synchronization | Clear data flow |

## 🔧 Core Features

### Virtual DOM Elements
- `dom.Div()`, `dom.Span()`, `dom.H1()` - Layout components  
- `dom.Input()` - Interactive text input with state
- `dom.Li()`, `dom.Ul()` - List components with focus
- `dom.Text()` - Text nodes

### Event System
- `OnKeyDown`, `OnChange`, `OnFocus`, `OnBlur`
- Event bubbling and delegation
- DOM-style event objects

### Props & State
- Type-safe props with `dom.ExtractProps[T]()`
- Automatic re-rendering on state changes
- Unidirectional data flow

## 🚧 Current Status

Go-DOM-TUI is in active development. Current features:
- ✅ Virtual DOM with efficient diffing
- ✅ Component system with props
- ✅ Event handling and bubbling  
- ✅ Focus management and tab order
- ✅ Integration with Charm ecosystem
- 🚧 React hooks (useState, useEffect)
- 🚧 Context API for global state
- 🚧 Performance optimizations

## 🤝 Contributing

Go-DOM-TUI aims to bring the best of web development to terminal UIs. We welcome contributions for:

- Additional components (forms, tables, menus)
- Performance optimizations
- React hooks implementation
- Documentation and examples
- Testing utilities

## 📄 License

MIT License - see LICENSE file for details.

---

**Go-DOM-TUI: Because terminal UIs deserve better developer experience.** 🚀