package dom

import (
	"testing"
)

// TestInputComponentProps tests the InputComponentProps methods
func TestInputComponentProps(t *testing.T) {
	// Create test callbacks
	onKeyDown := func(event *DOMEvent) { /* test callback */ }
	onChange := func(value string) { /* test callback */ }
	onFocus := func() { /* test callback */ }
	onBlur := func() { /* test callback */ }
	focusable := true
	tabIndex := 1

	props := &InputComponentProps{
		Placeholder: "Test placeholder",
		Value:       "test value",
		OnKeyDown:   onKeyDown,
		OnChange:    onChange,
		OnFocus:     onFocus,
		OnBlur:      onBlur,
		Focusable:   &focusable,
		TabIndex:    &tabIndex,
	}

	// Test Get method
	t.Run("Get", func(t *testing.T) {
		if placeholder := props.Get("placeholder"); placeholder != "Test placeholder" {
			t.Errorf("Expected placeholder 'Test placeholder', got %v", placeholder)
		}
		if value := props.Get("value"); value != "test value" {
			t.Errorf("Expected value 'test value', got %v", value)
		}
		if onKeyDownResult := props.Get("onKeyDown"); onKeyDownResult == nil {
			t.Error("Expected onKeyDown callback to be present")
		}
		if focusableResult := props.Get("focusable"); focusableResult != true {
			t.Errorf("Expected focusable true, got %v", focusableResult)
		}
		if tabIndexResult := props.Get("tabIndex"); tabIndexResult != 1 {
			t.Errorf("Expected tabIndex 1, got %v", tabIndexResult)
		}
		if unknown := props.Get("unknown"); unknown != nil {
			t.Errorf("Expected nil for unknown key, got %v", unknown)
		}
	})

	// Test GetString method
	t.Run("GetString", func(t *testing.T) {
		if placeholder := props.GetString("placeholder"); placeholder != "Test placeholder" {
			t.Errorf("Expected placeholder 'Test placeholder', got %s", placeholder)
		}
		if value := props.GetString("value"); value != "test value" {
			t.Errorf("Expected value 'test value', got %s", value)
		}
		if unknown := props.GetString("unknown"); unknown != "" {
			t.Errorf("Expected empty string for unknown key, got %s", unknown)
		}
	})

	// Test GetOK method
	t.Run("GetOK", func(t *testing.T) {
		if placeholder, ok := props.GetOK("placeholder"); !ok || placeholder != "Test placeholder" {
			t.Errorf("Expected placeholder 'Test placeholder' with ok=true, got %v, %v", placeholder, ok)
		}
		if value, ok := props.GetOK("value"); !ok || value != "test value" {
			t.Errorf("Expected value 'test value' with ok=true, got %v, %v", value, ok)
		}
		if _, ok := props.GetOK("onKeyDown"); !ok {
			t.Error("Expected onKeyDown callback to be present with ok=true")
		}
		if _, ok := props.GetOK("unknown"); ok {
			t.Error("Expected unknown key to return ok=false")
		}
	})

	// Test Range method
	t.Run("Range", func(t *testing.T) {
		keys := make(map[string]interface{})
		props.Range(func(key string, value interface{}) bool {
			keys[key] = value
			return true
		})

		expectedKeys := []string{"placeholder", "value", "onKeyDown", "onChange", "onFocus", "onBlur", "focusable", "tabIndex"}
		for _, key := range expectedKeys {
			if _, exists := keys[key]; !exists {
				t.Errorf("Expected key %s to be present in Range", key)
			}
		}
	})

	// Test Clone method
	t.Run("Clone", func(t *testing.T) {
		cloned := props.Clone().(*InputComponentProps)

		if cloned.Placeholder != props.Placeholder {
			t.Errorf("Expected cloned placeholder to match original")
		}
		if cloned.Value != props.Value {
			t.Errorf("Expected cloned value to match original")
		}
		if cloned.OnKeyDown == nil {
			t.Error("Expected cloned onKeyDown to be present")
		}
		if *cloned.Focusable != *props.Focusable {
			t.Errorf("Expected cloned focusable to match original")
		}
		if *cloned.TabIndex != *props.TabIndex {
			t.Errorf("Expected cloned tabIndex to match original")
		}
	})
}

// TestInputComponentPropsDefaults tests default values for InputComponentProps
func TestInputComponentPropsDefaults(t *testing.T) {
	props := &InputComponentProps{}

	// Test defaults for focusable and tabIndex
	t.Run("Defaults", func(t *testing.T) {
		if focusable := props.Get("focusable"); focusable != true {
			t.Errorf("Expected default focusable true, got %v", focusable)
		}
		if tabIndex := props.Get("tabIndex"); tabIndex != 0 {
			t.Errorf("Expected default tabIndex 0, got %v", tabIndex)
		}

		// Test GetOK returns defaults correctly
		if focusable, ok := props.GetOK("focusable"); !ok || focusable != true {
			t.Errorf("Expected default focusable true with ok=true, got %v, %v", focusable, ok)
		}
		if tabIndex, ok := props.GetOK("tabIndex"); !ok || tabIndex != 0 {
			t.Errorf("Expected default tabIndex 0 with ok=true, got %v, %v", tabIndex, ok)
		}
	})
}
