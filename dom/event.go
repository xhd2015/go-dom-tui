package dom

import "github.com/xhd2015/go-dom-tui/log"

// ========================================
// DOM Event System
// ========================================

type EventType string

const (
	// onKeydown
	EventTypeKeydown EventType = "keydown"
	EventTypeResize  EventType = "resize"
)

type KeyType string

const (
	KeyTypeEnter     KeyType = "enter"
	KeyTypeBackspace KeyType = "backspace"
	KeyTypeDelete    KeyType = "delete"
	KeyTypeTab       KeyType = "tab"
	KeyTypeEsc       KeyType = "esc"
	KeyTypeSpace     KeyType = "space"
	KeyTypeUp        KeyType = "up"
	KeyTypeDown      KeyType = "down"
	KeyTypeLeft      KeyType = "left"
	KeyTypeRight     KeyType = "right"
	KeyTypeCtrlC     KeyType = "ctrl+c"
	KeyTypeCtrlV     KeyType = "ctrl+v"
	KeyTypeCtrlX     KeyType = "ctrl+x"
	KeyTypeCtrlW     KeyType = "ctrl+w"
	KeyTypeCtrlA     KeyType = "ctrl+a"
	KeyTypeCtrlE     KeyType = "ctrl+e"
	KeyTypeCtrlK     KeyType = "ctrl+k"
)

// EventHandler represents a DOM event handler function
type EventHandler func(event *DOMEvent)

// WindowResizeEvent represents a window resize event
type WindowResizeEvent struct {
	Width  int
	Height int
}

// DOMEvent represents a DOM-like event
type DOMEvent struct {
	Type          EventType
	Target        *Node
	CurrentTarget *Node
	KeydownEvent  *KeydownEvent

	DefaultPrevented   bool
	PropagationStopped bool
	BubblePhase        bool
	WindowEvent        *WindowResizeEvent // For window-specific events
}

type KeydownEvent struct {
	KeyType KeyType
	Runes   []rune
	Alt     bool
	Paste   bool
}

// PreventDefault prevents the default behavior of the event
func (e *DOMEvent) PreventDefault() {
	e.DefaultPrevented = true
}

// StopPropagation stops the event from bubbling up
func (e *DOMEvent) StopPropagation() {
	e.PropagationStopped = true
}

// DispatchEvent dispatches an event to the focused node and bubbles it up
func (d *DOM) DispatchKeyDownEvent(keyEvent *KeydownEvent) {
	eventNode := d.Root.FindFocused()
	if eventNode == nil {
		log.Logf("DOM: DispatchEvent - no focused node, fallback to root node")
		// if no focused node, just send to root node
		eventNode = d.Root
	}

	log.Logf("DOM: DispatchKeyDownEvent keyType='%s' key='%s' to focused node %s", keyEvent.KeyType, keyEvent.Runes, eventNode.Type)

	// Create the event
	event := &DOMEvent{
		Type:          EventTypeKeydown,
		Target:        eventNode,
		CurrentTarget: eventNode,
		KeydownEvent:  keyEvent,
		BubblePhase:   false,
	}

	// Handle the event at the target and bubble up
	d.handleEventBubbling(eventNode, event)
	if !event.DefaultPrevented {
		// handle default event
		d.handleDefault(eventNode, event)
	}
}

// handleEventBubbling handles event bubbling up the DOM tree
func (d *DOM) handleEventBubbling(node *Node, event *DOMEvent) {
	if node == nil || event.PropagationStopped {
		log.Logf("DOM: handleEventBubbling - node is nil or event stopped")
		return
	}

	// Set current target
	event.CurrentTarget = node

	handler := node.GetEventHandler(event.Type)
	if handler != nil {
		log.Logf("DOM: handleEventBubbling - calling handler for node %s", node.Type)
		handler(event)
	} else {
		log.Logf("DOM: handleEventBubbling - no handler found for node %s", node.Type)
	}

	// If event wasn't stopped, bubble to parent
	log.Logf("DOM: handleEventBubbling - checking bubbling conditions: stopPropagation=%t, hasParent=%t",
		event.PropagationStopped, node.Parent != nil)
	if !event.PropagationStopped && node.Parent != nil {
		log.Logf("DOM: handleEventBubbling - bubbling to parent %s", node.Parent.Type)
		event.BubblePhase = true
		d.handleEventBubbling(node.Parent, event)
		return
	}

	log.Logf("DOM: handleEventBubbling - reached end of bubbling chain (stopPropagation=%t, hasParent=%t)",
		event.PropagationStopped, node.Parent != nil)
}

func (d *DOM) handleDefault(node *Node, event *DOMEvent) {
	if event.Type == EventTypeKeydown {
		keyEvent := event.KeydownEvent
		if keyEvent == nil {
			return
		}
		switch keyEvent.KeyType {
		case KeyTypeUp, KeyTypeDown:
			// handle focus navigation
			direction := 1
			if keyEvent.KeyType == KeyTypeUp {
				direction = -1
			}
			if d.HandleFocusNavigation(event, direction) {
				return
			}
		case KeyTypeLeft, KeyTypeRight:
			// move inside the input
			if node.Type == ElementTypeInput {
				props := ExtractProps[InputProps](node.Props)
				if props.OnCursorMove != nil {
					delta := 1
					if keyEvent.KeyType == KeyTypeLeft {
						delta = -1
					}
					props.OnCursorMove(props.CursorPosition + delta)
				}
			}
		default:
			// handle input
			if node.Type == ElementTypeInput {
				props := ExtractProps[InputProps](node.Props)

				// Get current value from the model if available, otherwise from props
				currentValue := props.Value

				// Update value based on key input
				newValue, newPos := UpdateInputValue(currentValue, props.CursorPosition, keyEvent)
				if newValue != currentValue {
					if props.OnChange != nil {
						props.OnChange(newValue)
					}
				}
				if newPos != props.CursorPosition {
					if props.OnCursorMove != nil {
						props.OnCursorMove(newPos)
					}
				}
			}
		}
	}
}
