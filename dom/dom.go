package dom

import (
	"github.com/xhd2015/go-dom-tui/log"
)

// TODO: use memo pool to reuse VNodes

// ========================================
// Virtual DOM Core Types
// ========================================

type IVNode interface {
	GetProps() any
	SetProps(p any)

	GetChildren() INodeList
}

type INodeList interface {
	Len() int
	Get(index int) IVNode
	Append(node IVNode)
	Insert(index int, node IVNode)
	Remove(index int)
	Clear()
}

type IElement interface {
	IVNode
}

type Props interface {
	Get(key string) (any, bool)
}

type RangeableProps interface {
	Props
	Range(func(key string, value any) bool)
}

// Node represents a virtual DOM node
type Node struct {
	Type string

	Props    Props   // Element properties
	Children []*Node // Child nodes
	Key      string  // Unique key for efficient diffing

	Text string // Text content (for text nodes)

	// Component-specific fields
	Component      Component
	State          interface{} // Component state
	Hooks          []Hook      // Component hooks
	EffectCleanups []func()    // Effect cleanup functions

	// DOM extension fields
	Parent *Node // Parent node for event bubbling
}

// Component represents a React-like component function
type Component func(props Props) *Node

// Hook represents a React hook
type Hook interface {
	GetType() string
	GetValue() interface{}
	SetValue(interface{})
}

// ========================================
// Virtual DOM Tree Operations
// ========================================

// DOM represents a DOM-like tree with event handling
type DOM struct {
	Root *Node
}

// NewDOM creates a new DOM from a VNode tree
func NewDOM(root *Node) *DOM {
	dom := &DOM{}
	dom.Root = dom.setupVNode(root, nil)
	return dom
}

// setupVNode recursively sets up VNodes with DOM functionality
func (d *DOM) setupVNode(vnode *Node, parent *Node) *Node {
	if vnode == nil {
		return nil
	}

	vnode.Parent = parent

	// Setup children
	for _, child := range vnode.Children {
		d.setupVNode(child, vnode)
	}

	return vnode
}

// DispatchEvent dispatches an event to the focused node and bubbles it up
func (d *DOM) DispatchEvent(eventType string, key string, model interface{}) interface{} {
	focusedNode := d.Root.FindFocused()
	if focusedNode == nil {
		log.Logf("DOM: DispatchEvent - no focused node")
		return nil
	}

	log.Logf("DOM: DispatchEvent %s key='%s' to focused node %s", eventType, key, focusedNode.Type)

	// Create the event
	event := &DOMEvent{
		Type:          eventType,
		Target:        focusedNode,
		CurrentTarget: focusedNode,
		Key:           key,
		BubblePhase:   false,
	}

	// Handle the event at the target and bubble up
	result := d.handleEventBubbling(focusedNode, event)
	log.Logf("DOM: DispatchEvent result: %v", result)
	return result
}

// handleEventBubbling handles event bubbling up the DOM tree
func (d *DOM) handleEventBubbling(node *Node, event *DOMEvent) interface{} {
	if node == nil || event.StopPropagation {
		log.Logf("DOM: handleEventBubbling - node is nil or event stopped")
		return nil
	}

	// Set current target
	event.CurrentTarget = node

	handler := node.GetEventHandler(event.Type)

	var handleResult interface{}
	var preventDefault bool
	if handler != nil {
		result := handler(event)
		if result != nil {
			handleResult = result
		}
		if event.DefaultPrevented {
			preventDefault = true
		}
	}

	if !preventDefault {
		if event.Type == "keydown" {
			switch event.Key {
			case "up", "down":
				// handle focus navigation
				direction := 1
				if event.Key == "up" {
					direction = -1
				}
				handled := d.HandleFocusNavigation(event, direction)
				if handled {
					return handleResult
				}
			case "left", "right":
				// move inside the input
				if node.Type == "input" {
					props := ExtractProps[InputComponentProps](node.Props)
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
					props := ExtractProps[InputComponentProps](node.Props)

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

	// If event wasn't stopped, bubble to parent
	log.Logf("DOM: handleEventBubbling - checking bubbling conditions: stopPropagation=%t, hasParent=%t",
		event.StopPropagation, node.Parent != nil)
	if !event.StopPropagation && node.Parent != nil {
		log.Logf("DOM: handleEventBubbling - bubbling to parent %s", node.Parent.Type)
		event.BubblePhase = true
		return d.handleEventBubbling(node.Parent, event)
	}

	log.Logf("DOM: handleEventBubbling - reached end of bubbling chain (stopPropagation=%t, hasParent=%t)",
		event.StopPropagation, node.Parent != nil)
	return handleResult
}

// SetFocus sets focus to a specific node
func (d *DOM) SetFocus(node *Node) {
	if node.IsFocused() {
		return
	}

	prev := d.Root.FindFocused()

	// Clear previous focus
	if prev != nil {
		if prev == node {
			return
		}
		log.Logf("DOM: SetFocus - clearing focus from %s", prev.Type)
		prev.SetFocused(false)
	}

	node.SetFocused(true)
}

// MoveFocus moves focus in the tab order
func (d *DOM) MoveFocus(direction int) bool {
	allFocusable := d.Root.FindAllFocusable()
	if len(allFocusable) == 0 {
		return false
	}

	n := len(allFocusable)
	j := -1
	for i := 0; i < n; i++ {
		node := allFocusable[i]
		if node.IsFocused() {
			j = i
			break
		}
	}

	if j == -1 {
		allFocusable[0].SetFocused(true)
		return true
	}

	x := j + direction
	if x < 0 {
		x = n - 1
	} else if x >= n {
		x = 0
	}
	if x == j {
		return true
	}
	allFocusable[j].SetFocused(false)
	allFocusable[x].SetFocused(true)
	return true
}

// HandleFocusNavigation handles focus navigation events through bubbling
func (d *DOM) HandleFocusNavigation(event *DOMEvent, direction int) bool {
	log.Logf("DOM: HandleFocusNavigation called with key='%s', direction=%d", event.Key, direction)

	// Only handle navigation keys
	switch event.Key {
	case "down":
		if direction == 0 {
			direction = 1
		}
		log.Logf("DOM: HandleFocusNavigation - handling down, direction=%d", direction)
	case "up":
		if direction == 0 {
			direction = -1
		}
		log.Logf("DOM: HandleFocusNavigation - handling up, direction=%d", direction)
	default:
		log.Logf("DOM: HandleFocusNavigation - key '%s' not a navigation key", event.Key)
		return false
	}

	// Move focus in the specified direction
	handled := d.MoveFocus(direction)
	if !handled {
		return false
	}

	// Prevent default and stop propagation since we handled the event
	event.PreventDefault()
	event.StopPropagationFunc()
	log.Logf("DOM: HandleFocusNavigation - event handled, prevented default and stopped propagation")

	return true
}
