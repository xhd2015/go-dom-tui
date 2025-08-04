package dom

// ========================================
// DOM Event System
// ========================================

// EventHandler represents a DOM event handler function
type EventHandler func(event *DOMEvent) interface{}

// DOMEvent represents a DOM-like event
type DOMEvent struct {
	Type             string
	Target           *Node
	CurrentTarget    *Node
	Key              string
	DefaultPrevented bool
	StopPropagation  bool
	BubblePhase      bool
}

// PreventDefault prevents the default behavior of the event
func (e *DOMEvent) PreventDefault() {
	e.DefaultPrevented = true
}

// StopPropagationFunc stops the event from bubbling up
func (e *DOMEvent) StopPropagationFunc() {
	e.StopPropagation = true
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
