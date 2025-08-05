package dom

import (
	"fmt"
	"reflect"
	"strings"
)

func ExtractProps[T any](props Props) T {
	sv, ok := props.(StructProps[T])
	if !ok {
		var t T
		panic(fmt.Errorf("ExtractProps[%T]: props is not a StructProps[%T]", props, t))
	}
	return sv.Value
}

type StructProps[T any] struct {
	Value T
}

func NewStructProps[T any](value T) StructProps[T] {
	if reflect.TypeOf(value).Kind() != reflect.Struct {
		panic(fmt.Errorf("NewStructProps[%T]: value is not a struct", value))
	}
	return StructProps[T]{Value: value}
}

// func (p StructProps[T]) Get(key string) any {
// 	rv := reflect.ValueOf(p.Value)
// 	if rv.Kind() != reflect.Struct {
// 		panic(fmt.Errorf("StructProps[%T]: value is not a struct", p.Value))
// 	}
// 	fv := rv.FieldByName(capitalize(key))
// 	if fv.IsValid() {
// 		return fv.Interface()
// 	}
// 	return nil
// }

func (p StructProps[T]) Get(key string) (any, bool) {
	rv := reflect.ValueOf(p.Value)
	if rv.Kind() != reflect.Struct {
		panic(fmt.Errorf("StructProps[%T]: value is not a struct", p.Value))
	}
	fv := rv.FieldByName(capitalize(key))
	if fv.IsValid() {
		return fv.Interface(), true
	}
	return nil, false
}

func (p StructProps[T]) GetString(key string) string {
	v, _ := p.Get(key)
	s, _ := v.(string)
	return s
}

func (p StructProps[T]) Range(fn func(key string, value any) bool) {
	rv := reflect.ValueOf(p.Value)
	if rv.Kind() != reflect.Struct {
		panic(fmt.Errorf("StructProps[%T]: value is not a struct", p.Value))
	}
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		if fv.IsValid() {
			fn(decapitalize(rv.Type().Field(i).Name), fv.Interface())
		}
	}
}

func (p StructProps[T]) Clone() Props {
	rv := reflect.ValueOf(p.Value)
	if rv.Kind() != reflect.Struct {
		panic(fmt.Errorf("StructProps[%T]: value is not a struct", p.Value))
	}
	clone := reflect.New(rv.Type())
	clone.Elem().Set(rv)

	var v interface{} = clone.Interface()
	return NewStructProps(v.(T))
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func decapitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// GetStringProp safely gets a string property
func GetStringProp(props Props, key string) string {
	if props == nil {
		return ""
	}
	v, ok := props.Get(key)
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

type TextNodeProps struct {
	Style     Style
	Focused   bool
	Focusable bool

	OnKeyDown func(*DOMEvent)
}

// ButtonProps represents props for button elements
type ButtonProps struct {
	Text      string
	OnClick   func()
	Style     string
	Focusable *bool // Optional: nil = default, true/false = explicit
	TabIndex  *int  // Optional: nil = default, number = explicit
}

// DivProps represents props for div elements
type DivProps struct {
	Style Style

	OnKeyDown      func(*DOMEvent)
	OnWindowResize func(*DOMEvent)

	Focused   bool
	Focusable bool
}

// CounterProps represents props for Counter component
type CounterProps struct {
	InitialValue int
}

// TodoItemProps represents props for TodoItem component
type TodoItemProps struct {
	ID        int
	Text      string
	Completed bool
	OnToggle  func(int)
	OnSave    func(int, string)
	OnDelete  func(int)
}

// TodoListProps represents props for TodoList component
type TodoListProps struct {
	Todos    []map[string]interface{}
	OnToggle func(int)
	OnSave   func(int, string)
	OnDelete func(int)
}

// TodoAppProps represents props for TodoApp component
type TodoAppProps struct {
	Title        string
	InitialTodos []map[string]interface{}
}

// ElementProps represents props for basic HTML elements (h1, etc.)
type ElementProps struct {
	Style     string
	Focusable *bool // Optional: nil = default, true/false = explicit
	TabIndex  *int  // Optional: nil = default, number = explicit
}

// EmptyProps represents empty props
type EmptyProps struct{}

// InputProps represents props for input component
type InputProps struct {
	Placeholder string // Input placeholder text
	Value       string // Current input value

	CursorPosition int // Cursor position
	OnCursorMove   func(delta int, seek int)

	OnKeyDown func(e *DOMEvent) // Key down callback
	OnChange  func(string)      // Value change callback
	OnFocus   func()            // Focus callback
	OnBlur    func()            // Blur callback

	Focused bool // Whether the input is focused

	Focusable *bool // Optional: nil = default (true for input), true/false = explicit
}

// Focusable creates a boolean pointer for focusable property
func Focusable(value bool) *bool {
	return &value
}
func String(value string) *string {
	return &value
}

// ListItemProps represents props for focusable li elements
type ListItemProps struct {
	Style      Style
	Index      int
	Selected   bool
	ItemPrefix *string
	Focused    bool
	OnFocus    func()
	OnBlur     func()
	OnKeyDown  func(e *DOMEvent)
	Focusable  *bool
}
