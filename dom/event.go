package dom

// ========================================
// DOM Event System
// ========================================

// EventHandler represents a DOM event handler function
type EventHandler func(event *DOMEvent) interface{}

// WindowResizeEvent represents a window resize event
type WindowResizeEvent struct {
	Width  int
	Height int
}

// DOMEvent represents a DOM-like event
type DOMEvent struct {
	Type               string
	Target             *Node
	CurrentTarget      *Node
	Key                string
	DefaultPrevented   bool
	PropagationStopped bool
	BubblePhase        bool
	WindowEvent        *WindowResizeEvent // For window-specific events
}

// PreventDefault prevents the default behavior of the event
func (e *DOMEvent) PreventDefault() {
	e.DefaultPrevented = true
}

// StopPropagation stops the event from bubbling up
func (e *DOMEvent) StopPropagation() {
	e.PropagationStopped = true
}

func CombineResults(a interface{}, b interface{}) interface{} {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}

	return a
}
