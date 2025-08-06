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
	Parent *Node   // Parent node for event bubbling
	Window *Window // Reference to global window state
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
	Root   *Node
	Window *Window
}

// NewDOM creates a new DOM from a VNode tree
func NewDOM(root *Node, window *Window) *DOM {
	dom := &DOM{
		Window: window,
	}

	dom.Root = dom.setupVNode(root, nil, window)

	return dom
}

// setupVNode recursively sets up VNodes with DOM functionality
func (d *DOM) setupVNode(vnode *Node, parent *Node, window *Window) *Node {
	if vnode == nil {
		return nil
	}

	vnode.Parent = parent
	vnode.Window = window // Set window reference on all nodes

	// Setup children
	for _, child := range vnode.Children {
		d.setupVNode(child, vnode, window)
	}

	return vnode
}

// DispatchEvent dispatches an event to the focused node and bubbles it up
func (d *DOM) DispatchKeyDownEvent(keyEvent *KeydownEvent) {
	eventNode := d.Root.FindFocused()
	if eventNode == nil {
		log.Logf("DOM: DispatchEvent - no focused node, fallback to root node")
		// if no focused node, just send to root node
		eventNode = d.Root
	}

	log.Logf("DOM: DispatchKeyDownEvent %s key='%s' to focused node %s", keyEvent.KeyType, eventNode.Type)

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
	keyEvent := event.KeydownEvent
	if keyEvent == nil {
		return false
	}
	log.Logf("DOM: HandleFocusNavigation called with key='%s', direction=%d", keyEvent.KeyType, direction)

	// Only handle navigation keys
	switch keyEvent.KeyType {
	case KeyTypeDown:
		if direction == 0 {
			direction = 1
		}
		log.Logf("DOM: HandleFocusNavigation - handling down, direction=%d", direction)
	case KeyTypeUp:
		if direction == 0 {
			direction = -1
		}
		log.Logf("DOM: HandleFocusNavigation - handling up, direction=%d", direction)
	default:
		log.Logf("DOM: HandleFocusNavigation - key '%s' not a navigation key", keyEvent.KeyType)
		return false
	}

	// Move focus in the specified direction
	handled := d.MoveFocus(direction)
	if !handled {
		return false
	}

	// Prevent default and stop propagation since we handled the event
	event.PreventDefault()
	event.StopPropagation()
	log.Logf("DOM: HandleFocusNavigation - event handled, prevented default and stopped propagation")

	return true
}

// DispatchWindowEvent dispatches window-level events (like resize) to the DOM tree
func (d *DOM) DispatchWindowEvent(eventType EventType, windowEvent *WindowResizeEvent) {
	log.Logf("DOM: DispatchWindowEvent %s - %dx%d", eventType, windowEvent.Width, windowEvent.Height)

	// Create window event that targets the root
	event := &DOMEvent{
		Type:          eventType,
		Target:        d.Root,
		CurrentTarget: d.Root,
		WindowEvent:   windowEvent,
		BubblePhase:   false,
	}

	// Handle at root level and propagate to interested components
	d.handleWindowEventPropagation(d.Root, event)
}

// handleWindowEventPropagation propagates window events through the DOM tree
func (d *DOM) handleWindowEventPropagation(node *Node, event *DOMEvent) {
	if node == nil {
		return
	}

	// Set current target
	event.CurrentTarget = node

	// Check if this node has a window event handler
	if handler := node.GetEventHandler(event.Type); handler != nil {
		handler(event)
	}

	// Propagate to all children (not bubbling, but broadcasting)
	for _, child := range node.Children {
		if child != nil {
			d.handleWindowEventPropagation(child, event)
		}
	}
}
