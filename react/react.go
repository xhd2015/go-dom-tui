package react

import (
	"reflect"
	"sync"

	"github.com/xhd2015/go-dom-tui/dom"
)

// ========================================
// React Infrastructure
// ========================================

// FunctionalComponent represents a functional component type
type FunctionalComponent func(props dom.Props, react *React) *dom.Node

// Renderer handles converting VNodes to actual output
type Renderer interface {
	Render(node *dom.Node) string
	ApplyPatch(patch Patch, currentOutput string) string
}

// PatchType represents the type of patch to apply
type PatchType int

const (
	PatchCreate PatchType = iota
	PatchUpdate
	PatchDelete
	PatchReplace
	PatchReorder
)

// Hook represents a React hook
type Hook interface {
	GetType() string
	GetValue() interface{}
	SetValue(interface{})
}

// Patch represents a change to apply to the DOM
type Patch struct {
	Type    PatchType
	Path    []int                  // Path to the node in the tree
	OldNode *dom.Node              // Old node (for updates/deletes)
	NewNode *dom.Node              // New node (for creates/updates)
	Props   map[string]interface{} // Property changes
}

// React represents the main React instance
type React struct {
	renderer    Renderer
	hookIndex   int
	hooks       []Hook
	isRendering bool
	updateQueue []func()
	mu          sync.RWMutex
}

// NewReact creates a new React instance
func NewReact(renderer Renderer) *React {
	return &React{
		renderer:    renderer,
		hooks:       make([]Hook, 0),
		updateQueue: make([]func(), 0),
	}
}

// ========================================
// Fiber Architecture (simplified React Fiber)
// ========================================

// Fiber represents a unit of work in the React reconciliation
type Fiber struct {
	Type        interface{}   // Component type
	Props       dom.Props     // Props
	State       interface{}   // State
	Hooks       []Hook        // Hooks
	VNode       *dom.Node     // Virtual node
	Parent      *Fiber        // Parent fiber
	Child       *Fiber        // First child
	Sibling     *Fiber        // Next sibling
	Alternate   *Fiber        // Previous version
	EffectTag   EffectTag     // What needs to be done
	UpdateQueue []StateUpdate // Pending state updates
}

// EffectTag represents what kind of work needs to be done
type EffectTag int

const (
	NoEffect EffectTag = iota
	Placement
	Update
	Deletion
)

// StateUpdate represents a state update
type StateUpdate struct {
	UpdaterFunc func(interface{}) interface{}
	Callback    func()
}

// ========================================
// React Hooks System
// ========================================

// HookType represents the type of hook
type HookType int

const (
	StateHook HookType = iota
	EffectHook
	ContextHook
	MemoHook
	CallbackHook
)

// StateHookValue represents a useState hook
type StateHookValue struct {
	State    interface{}
	SetState func(interface{})
}

// GetType returns the hook type
func (s *StateHookValue) GetType() string {
	return "state"
}

// GetValue returns the hook value
func (s *StateHookValue) GetValue() interface{} {
	return s.State
}

// SetValue sets the hook value
func (s *StateHookValue) SetValue(value interface{}) {
	s.State = value
}

// EffectHookValue represents a useEffect hook
type EffectHookValue struct {
	Effect  func()
	Cleanup func()
	Deps    []interface{}
	HasRun  bool
}

// GetType returns the hook type
func (e *EffectHookValue) GetType() string {
	return "effect"
}

// GetValue returns the hook value
func (e *EffectHookValue) GetValue() interface{} {
	return e
}

// SetValue sets the hook value
func (e *EffectHookValue) SetValue(value interface{}) {
	if effect, ok := value.(*EffectHookValue); ok {
		*e = *effect
	}
}

// ========================================
// React Hooks Implementation
// ========================================

// UseState implements React's useState hook
func (r *React) UseState(initialState interface{}) (interface{}, func(interface{})) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isRendering {
		panic("useState can only be called during rendering")
	}

	// Get or create hook
	var hook *StateHookValue
	if r.hookIndex < len(r.hooks) {
		// Use existing hook
		if existingHook, ok := r.hooks[r.hookIndex].(*StateHookValue); ok {
			hook = existingHook
		} else {
			panic("Hook type mismatch")
		}
	} else {
		// Create new hook
		hook = &StateHookValue{
			State: initialState,
		}

		// Create setState function
		hook.SetState = func(newState interface{}) {
			r.mu.Lock()
			defer r.mu.Unlock()

			// Handle functional updates
			if updateFunc, ok := newState.(func(interface{}) interface{}); ok {
				hook.State = updateFunc(hook.State)
			} else {
				hook.State = newState
			}

			// Schedule re-render
			r.scheduleUpdate()
		}

		r.hooks = append(r.hooks, hook)
	}

	r.hookIndex++
	return hook.State, hook.SetState
}

// UseEffect implements React's useEffect hook
func (r *React) UseEffect(effect func(), cleanup func(), deps []interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isRendering {
		panic("useEffect can only be called during rendering")
	}

	// Get or create hook
	var hook *EffectHookValue
	if r.hookIndex < len(r.hooks) {
		// Use existing hook
		if existingHook, ok := r.hooks[r.hookIndex].(*EffectHookValue); ok {
			hook = existingHook
		} else {
			panic("Hook type mismatch")
		}
	} else {
		// Create new hook
		hook = &EffectHookValue{
			Effect:  effect,
			Cleanup: cleanup,
			Deps:    deps,
			HasRun:  false,
		}
		r.hooks = append(r.hooks, hook)
	}

	// Check if effect should run
	shouldRun := !hook.HasRun || r.depsChanged(hook.Deps, deps)

	if shouldRun {
		// Run cleanup if it exists
		if hook.Cleanup != nil && hook.HasRun {
			hook.Cleanup()
		}

		// Schedule effect to run after render
		r.updateQueue = append(r.updateQueue, func() {
			effect()
			hook.HasRun = true
		})

		// Update deps
		hook.Deps = make([]interface{}, len(deps))
		copy(hook.Deps, deps)
		hook.Effect = effect
		hook.Cleanup = cleanup
	}

	r.hookIndex++
}

// depsChanged checks if dependencies have changed
func (r *React) depsChanged(oldDeps, newDeps []interface{}) bool {
	if len(oldDeps) != len(newDeps) {
		return true
	}

	for i, oldDep := range oldDeps {
		if oldDep != newDeps[i] {
			return true
		}
	}

	return false
}

// scheduleUpdate schedules a re-render
func (r *React) scheduleUpdate() {
	// In a real implementation, this would use a scheduler
	// For now, we'll just mark that an update is needed
	go func() {
		// This would trigger a re-render in the next tick
		r.forceUpdate()
	}()
}

// forceUpdate forces a re-render
func (r *React) forceUpdate() {
	// This would be called by the scheduler
	// Implementation depends on the specific use case
}

// ========================================
// React Context API
// ========================================

// Context represents a React context
type Context struct {
	DefaultValue interface{}
	Providers    []*ContextProvider
}

// ContextProvider provides context value to children
type ContextProvider struct {
	Context  *Context
	Value    interface{}
	Children []*dom.Node
}

func (c *ContextProvider) Render() *dom.Node {
	// TODO
	return nil
}

// CreateContext creates a new React context
func CreateContext(defaultValue interface{}) *Context {
	return &Context{
		DefaultValue: defaultValue,
		Providers:    make([]*ContextProvider, 0),
	}
}

// Provider creates a context provider element
func (c *Context) Provider(value interface{}, children ...*dom.Node) *dom.Node {
	provider := &ContextProvider{
		Context:  c,
		Value:    value,
		Children: make([]*dom.Node, len(children)),
	}

	copy(provider.Children, children)

	return dom.CreateComponent(func(props dom.Props) *dom.Node {
		panic("TODO Provider")
		return nil
	}, dom.NewStructProps(dom.ElementProps{}), children...)
}

// UseContext implements React's useContext hook
func (r *React) UseContext(context *Context) interface{} {
	// In a real implementation, this would walk up the fiber tree
	// to find the nearest provider
	if len(context.Providers) > 0 {
		return context.Providers[len(context.Providers)-1].Value
	}
	return context.DefaultValue
}

// ========================================
// React Component Lifecycle
// ========================================

// ComponentLifecycle represents component lifecycle methods
type ComponentLifecycle interface {
	ComponentDidMount()
	ComponentDidUpdate(prevProps map[string]interface{})
	ComponentWillUnmount()
}

// ClassComponent represents a class-based component (for compatibility)
type ClassComponent struct {
	Props    map[string]interface{}
	State    interface{}
	SetState func(interface{})
	React    *React
}

// NewClassComponent creates a new class component
func NewClassComponent(react *React, props map[string]interface{}) *ClassComponent {
	component := &ClassComponent{
		Props: props,
		React: react,
	}

	// Create setState method
	component.SetState = func(newState interface{}) {
		if updateFunc, ok := newState.(func(interface{}) interface{}); ok {
			component.State = updateFunc(component.State)
		} else {
			component.State = newState
		}
		react.scheduleUpdate()
	}

	return component
}

// ========================================
// React Reconciler
// ========================================

// Reconciler handles the reconciliation process
type Reconciler struct {
	react       *React
	rootFiber   *Fiber
	currentRoot *Fiber
	deletions   []*Fiber
}

// NewReconciler creates a new reconciler
func NewReconciler(react *React) *Reconciler {
	return &Reconciler{
		react:     react,
		deletions: make([]*Fiber, 0),
	}
}

// GetCurrentRootVNode returns the VNode from the current root fiber
func (rec *Reconciler) GetCurrentRootVNode() *dom.Node {
	if rec.currentRoot != nil {
		return rec.currentRoot.VNode
	}
	return nil
}

// Render starts the reconciliation process
func (rec *Reconciler) Render(element *dom.Node, container interface{}) {
	// Create or update root fiber
	if rec.rootFiber == nil {
		rec.rootFiber = &Fiber{
			Type:  "root",
			Props: dom.NewStructProps(dom.EmptyProps{}),
		}
	}

	// Create work-in-progress fiber
	wipRoot := &Fiber{
		Type:      rec.rootFiber.Type,
		Props:     rec.rootFiber.Props,
		Alternate: rec.rootFiber,
	}

	// Convert element to VNode once
	vnode := element
	wipRoot.VNode = vnode

	// Start reconciliation with pre-converted VNodes
	rec.reconcileChildrenWithVNodes(wipRoot, []*dom.Node{element}, []*dom.Node{vnode})

	// Commit the changes
	rec.commitRoot(wipRoot)

	// Update root reference
	rec.currentRoot = wipRoot
	rec.rootFiber = wipRoot
}

// reconcileChildrenWithVNodes reconciles child fibers using pre-converted VNodes
func (rec *Reconciler) reconcileChildrenWithVNodes(wipFiber *Fiber, elements []*dom.Node, vnodes []*dom.Node) {
	var oldFiber *Fiber
	if wipFiber.Alternate != nil {
		oldFiber = wipFiber.Alternate.Child
	}

	var prevSibling *Fiber

	for i, element := range elements {
		var newFiber *Fiber
		var vnode *dom.Node
		if i < len(vnodes) {
			vnode = vnodes[i]
		}

		// Check if we can reuse the old fiber
		sameType := oldFiber != nil && element != nil &&
			reflect.DeepEqual(oldFiber.Type, element.Type)

		if sameType {
			// Update existing fiber - reuse pre-converted VNode
			newFiber = &Fiber{
				Type:      oldFiber.Type,
				Props:     element.Props,
				VNode:     vnode, // Reuse pre-converted VNode
				Parent:    wipFiber,
				Alternate: oldFiber,
				EffectTag: Update,
			}
		} else {
			// Create new fiber - reuse pre-converted VNode
			if element != nil {
				newFiber = &Fiber{
					Type:      element.Type,
					Props:     element.Props,
					VNode:     vnode, // Reuse pre-converted VNode
					Parent:    wipFiber,
					Alternate: nil,
					EffectTag: Placement,
				}
			}

			// Delete old fiber
			if oldFiber != nil {
				oldFiber.EffectTag = Deletion
				rec.deletions = append(rec.deletions, oldFiber)
			}
		}

		// Move to next old fiber
		if oldFiber != nil {
			oldFiber = oldFiber.Sibling
		}

		// Link fibers
		if i == 0 {
			wipFiber.Child = newFiber
		} else if prevSibling != nil {
			prevSibling.Sibling = newFiber
		}

		prevSibling = newFiber
	}
}

// reconcileChildren reconciles child fibers
func (rec *Reconciler) reconcileChildren(wipFiber *Fiber, elements []*dom.Node) {
	var oldFiber *Fiber
	if wipFiber.Alternate != nil {
		oldFiber = wipFiber.Alternate.Child
	}

	var prevSibling *Fiber

	for i, element := range elements {
		var newFiber *Fiber

		// Check if we can reuse the old fiber
		sameType := oldFiber != nil && element != nil &&
			reflect.DeepEqual(oldFiber.Type, element.Type)

		if sameType {
			// Update existing fiber
			newFiber = &Fiber{
				Type:      oldFiber.Type,
				Props:     element.Props,
				VNode:     element,
				Parent:    wipFiber,
				Alternate: oldFiber,
				EffectTag: Update,
			}
		} else {
			// Create new fiber
			if element != nil {
				newFiber = &Fiber{
					Type:      element.Type,
					Props:     element.Props,
					VNode:     element,
					Parent:    wipFiber,
					Alternate: nil,
					EffectTag: Placement,
				}
			}

			// Delete old fiber
			if oldFiber != nil {
				oldFiber.EffectTag = Deletion
				rec.deletions = append(rec.deletions, oldFiber)
			}
		}

		// Move to next old fiber
		if oldFiber != nil {
			oldFiber = oldFiber.Sibling
		}

		// Link fibers
		if i == 0 {
			wipFiber.Child = newFiber
		} else if prevSibling != nil {
			prevSibling.Sibling = newFiber
		}

		prevSibling = newFiber
	}
}

// commitRoot commits all changes to the DOM
func (rec *Reconciler) commitRoot(root *Fiber) {
	// Delete removed nodes
	for _, fiber := range rec.deletions {
		rec.commitWork(fiber)
	}
	rec.deletions = rec.deletions[:0]

	// Commit all changes
	rec.commitWork(root.Child)
}

// commitWork commits a single fiber's changes
func (rec *Reconciler) commitWork(fiber *Fiber) {
	if fiber == nil {
		return
	}

	switch fiber.EffectTag {
	case Placement:
		// Add new node
		if fiber.VNode != nil {
			// In a real implementation, this would add to the actual DOM
			rec.react.renderer.Render(fiber.VNode)
		}

	case Update:
		// Update existing node
		if fiber.VNode != nil && fiber.Alternate != nil && fiber.Alternate.VNode != nil {
			// Diff and apply changes
			diff := Diff(fiber.Alternate.VNode, fiber.VNode)
			for _, patch := range diff.Patches {
				rec.react.renderer.ApplyPatch(patch, "")
			}
		}

	case Deletion:
		// Remove node
		// In a real implementation, this would remove from the actual DOM
		break
	}

	// Recursively commit children and siblings
	rec.commitWork(fiber.Child)
	rec.commitWork(fiber.Sibling)
}

// ========================================
// React Application
// ========================================

// Root represents a React application
type Root struct {
	React       *React
	Reconciler  *Reconciler
	RootElement *dom.Node
	Container   interface{}
}

// CreateRoot creates a new React application
func CreateRoot(renderer Renderer) *Root {
	react := NewReact(renderer)
	reconciler := NewReconciler(react)

	return &Root{
		React:      react,
		Reconciler: reconciler,
	}
}

// Render renders the application
func (app *Root) Render(element *dom.Node, container interface{}) {
	app.RootElement = element
	app.Container = container
	app.Reconciler.Render(element, container)
}

// Update triggers a re-render
func (app *Root) Update() {
	if app.RootElement != nil {
		app.Reconciler.Render(app.RootElement, app.Container)
	}
}
