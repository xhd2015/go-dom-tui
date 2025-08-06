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

	FirstFocusable *Node
	LastFocusable  *Node

	FocusedNode        *Node
	PreviousFocuseable *Node
	NextFocuseable     *Node
}

// NewDOM creates a new DOM from a VNode tree
func NewDOM(root *Node, window *Window) *DOM {
	dom := &DOM{
		Root:   root,
		Window: window,
	}

	dom.setupVNode(root, nil, window)

	return dom
}

// setupVNode recursively sets up VNodes with DOM functionality
func (d *DOM) setupVNode(vnode *Node, parent *Node, window *Window) {
	if vnode == nil {
		return
	}

	focusable := vnode.IsFocusable()
	focused := vnode.IsFocused()

	vnode.Parent = parent
	vnode.Window = window // Set window reference on all nodes

	// Track first focusable node
	if focusable && d.FirstFocusable == nil {
		d.FirstFocusable = vnode
	}

	// Track focused node and set previous focusable
	if focused && d.FocusedNode == nil {
		d.FocusedNode = vnode
		d.PreviousFocuseable = d.LastFocusable // Previous is the last focusable we've seen
	}

	// If we have a focused node but no next focusable yet, and this is focusable
	if d.FocusedNode != nil && d.NextFocuseable == nil && focusable && vnode != d.FocusedNode {
		d.NextFocuseable = vnode
	}

	// Update last focusable before processing children
	if focusable {
		d.LastFocusable = vnode
	}

	// Process children in depth-first order
	for _, child := range vnode.Children {
		if child == nil {
			continue
		}
		d.setupVNode(child, vnode, window)
	}
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
