package dom

import "github.com/xhd2015/go-dom-tui/log"

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

// handleEventBubbling handles event bubbling up the DOM tree
func (d *DOM) handleEventBubbling(node *Node, event *DOMEvent) interface{} {
	if node == nil || event.PropagationStopped {
		log.Logf("DOM: handleEventBubbling - node is nil or event stopped")
		return nil
	}

	// Set current target
	event.CurrentTarget = node

	handler := node.GetEventHandler(event.Type)

	var handleResult interface{}
	if handler != nil {
		result := handler(event)
		if result != nil {
			handleResult = result
		}
	}

	// If event wasn't stopped, bubble to parent
	log.Logf("DOM: handleEventBubbling - checking bubbling conditions: stopPropagation=%t, hasParent=%t",
		event.PropagationStopped, node.Parent != nil)
	if !event.PropagationStopped && node.Parent != nil {
		log.Logf("DOM: handleEventBubbling - bubbling to parent %s", node.Parent.Type)
		event.BubblePhase = true
		return d.handleEventBubbling(node.Parent, event)
	}

	log.Logf("DOM: handleEventBubbling - reached end of bubbling chain (stopPropagation=%t, hasParent=%t)",
		event.PropagationStopped, node.Parent != nil)
	return handleResult
}

func (d *DOM) handleDefault(node *Node, event *DOMEvent) {
	if event.Type == "keydown" {
		switch event.Key {
		case "up", "down":
			// handle focus navigation
			direction := 1
			if event.Key == "up" {
				direction = -1
			}
			if d.HandleFocusNavigation(event, direction) {
				return
			}
		case "left", "right":
			// move inside the input
			if node.Type == "input" {
				props := ExtractProps[InputProps](node.Props)
				if props.OnCursorMove != nil {
					delta := 1
					if event.Key == "left" {
						delta = -1
					}
					props.OnCursorMove(delta, 0)
				}
			}
		default:
			// handle input
			if node.Type == "input" {
				props := ExtractProps[InputProps](node.Props)

				// Get current value from the model if available, otherwise from props
				currentValue := props.Value

				// Update value based on key input
				newValue, newPos := UpdateInputValue(currentValue, props.CursorPosition, event.Key)
				if newPos != props.CursorPosition {
					if props.OnCursorMove != nil {
						props.OnCursorMove(newPos-props.CursorPosition, 0)
					}
				}
				if newValue != currentValue {
					if props.OnChange != nil {
						props.OnChange(newValue)
					}
				}
			}
		}
	}
}
