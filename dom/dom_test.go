package dom

import (
	"testing"
)

// TestSetupEventHandlersFromProps tests the setupEventHandlersFromProps function
func TestSetupEventHandlersFromProps(t *testing.T) {
	// Create test callbacks
	var keyDownCalled bool
	var keyDownKey string
	onKeyDown := func(key string) {
		keyDownCalled = true
		keyDownKey = key
	}

	var focusCalled bool
	onFocus := func() {
		focusCalled = true
	}

	var blurCalled bool
	onBlur := func() {
		blurCalled = true
	}

	var changeCalled bool
	var changeValue string
	onChange := func(value string) {
		changeCalled = true
		changeValue = value
	}

	// Create props with event handlers
	props := &InputComponentProps{
		OnKeyDown: onKeyDown,
		OnFocus:   onFocus,
		OnBlur:    onBlur,
		OnChange:  onChange,
		Value:     "test-value",
	}

	// Create VNode
	vnode := &Node{
		Type:  "input",
		Props: props,
	}

	// Create DOM and setup event handlers
	dom := &DOM{}
	dom.setupEventHandlersFromProps(vnode)

	// Test that event handlers were set up
	t.Run("EventHandlersSetup", func(t *testing.T) {
		if vnode.EventHandlers == nil {
			t.Fatal("Expected EventHandlers to be initialized")
		}

		if _, exists := vnode.EventHandlers["keydown"]; !exists {
			t.Error("Expected keydown handler to be set up")
		}
		if _, exists := vnode.EventHandlers["focus"]; !exists {
			t.Error("Expected focus handler to be set up")
		}
		if _, exists := vnode.EventHandlers["blur"]; !exists {
			t.Error("Expected blur handler to be set up")
		}
		if _, exists := vnode.EventHandlers["change"]; !exists {
			t.Error("Expected change handler to be set up")
		}
	})

	// Test keydown handler execution
	t.Run("KeyDownHandler", func(t *testing.T) {
		handler := vnode.EventHandlers["keydown"]
		event := &DOMEvent{
			Type: "keydown",
			Key:  "enter",
		}

		keyDownCalled = false
		keyDownKey = ""

		handler(event, nil)

		if !keyDownCalled {
			t.Error("Expected onKeyDown callback to be called")
		}
		if keyDownKey != "enter" {
			t.Errorf("Expected key 'enter', got '%s'", keyDownKey)
		}
	})

	// Test focus handler execution
	t.Run("FocusHandler", func(t *testing.T) {
		handler := vnode.EventHandlers["focus"]
		event := &DOMEvent{
			Type: "focus",
		}

		focusCalled = false

		handler(event, nil)

		if !focusCalled {
			t.Error("Expected onFocus callback to be called")
		}
	})

	// Test blur handler execution
	t.Run("BlurHandler", func(t *testing.T) {
		handler := vnode.EventHandlers["blur"]
		event := &DOMEvent{
			Type: "blur",
		}

		blurCalled = false

		handler(event, nil)

		if !blurCalled {
			t.Error("Expected onBlur callback to be called")
		}
	})

	// Test change handler execution
	t.Run("ChangeHandler", func(t *testing.T) {
		handler := vnode.EventHandlers["change"]
		event := &DOMEvent{
			Type: "change",
		}

		changeCalled = false
		changeValue = ""

		handler(event, nil)

		if !changeCalled {
			t.Error("Expected onChange callback to be called")
		}
		if changeValue != "test-value" {
			t.Errorf("Expected value 'test-value', got '%s'", changeValue)
		}
	})
}

// TestSetupEventHandlersFromPropsWithNilProps tests with nil props
func TestSetupEventHandlersFromPropsWithNilProps(t *testing.T) {
	vnode := &Node{
		Type:  "input",
		Props: nil,
	}

	dom := &DOM{}

	// Should not panic with nil props
	dom.setupEventHandlersFromProps(vnode)

	// EventHandlers should remain nil or empty
	if vnode.EventHandlers != nil && len(vnode.EventHandlers) > 0 {
		t.Error("Expected no event handlers to be set up with nil props")
	}
}

// TestSetupEventHandlersFromPropsWithNoCallbacks tests with props but no callbacks
func TestSetupEventHandlersFromPropsWithNoCallbacks(t *testing.T) {
	props := &InputComponentProps{
		Placeholder: "test",
		Value:       "test-value",
		// No callbacks set
	}

	vnode := &Node{
		Type:  "input",
		Props: props,
	}

	dom := &DOM{}
	dom.setupEventHandlersFromProps(vnode)

	// Should not set up any event handlers
	if vnode.EventHandlers != nil && len(vnode.EventHandlers) > 0 {
		t.Error("Expected no event handlers to be set up with no callbacks")
	}
}

// TestSetupFocusFromProps tests the setupFocusFromProps function
func TestSetupFocusFromProps(t *testing.T) {
	t.Run("ExplicitFocusable", func(t *testing.T) {
		focusable := true
		tabIndex := 5
		props := &InputComponentProps{
			Focusable: &focusable,
			TabIndex:  &tabIndex,
		}

		vnode := &Node{
			Type:  "input",
			Props: props,
		}

		dom := &DOM{}
		dom.setupFocusFromProps(vnode)

		if !vnode.Focusable {
			t.Error("Expected vnode to be focusable")
		}
		if vnode.TabIndex != 5 {
			t.Errorf("Expected tabIndex 5, got %d", vnode.TabIndex)
		}
	})

	t.Run("DefaultFocusableForInput", func(t *testing.T) {
		props := &InputComponentProps{}

		vnode := &Node{
			Type:  "input",
			Props: props,
		}

		dom := &DOM{}
		dom.setupFocusFromProps(vnode)

		if !vnode.Focusable {
			t.Error("Expected input to be focusable by default")
		}
		if vnode.TabIndex != 0 {
			t.Errorf("Expected default tabIndex 0, got %d", vnode.TabIndex)
		}
	})

	t.Run("DefaultFocusableForButton", func(t *testing.T) {
		props := &EmptyProps{}

		vnode := &Node{
			Type:  "button",
			Props: props,
		}

		dom := &DOM{}
		dom.setupFocusFromProps(vnode)

		if !vnode.Focusable {
			t.Error("Expected button to be focusable by default")
		}
		if vnode.TabIndex != 0 {
			t.Errorf("Expected default tabIndex 0, got %d", vnode.TabIndex)
		}
	})

	t.Run("NonFocusableElement", func(t *testing.T) {
		props := &EmptyProps{}

		vnode := &Node{
			Type:  "div",
			Props: props,
		}

		dom := &DOM{}
		dom.setupFocusFromProps(vnode)

		if vnode.Focusable {
			t.Error("Expected div to not be focusable by default")
		}
		if vnode.TabIndex != 0 {
			t.Errorf("Expected default tabIndex 0, got %d", vnode.TabIndex)
		}
	})
}
