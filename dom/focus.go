package dom

import (
	"fmt"

	"github.com/xhd2015/go-dom-tui/log"
)

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

func (d *DOM) MoveFocus(direction int) bool {
	curFocus := d.FocusedNode
	nextFocus := d.NextFocuseable
	if direction < 0 {
		nextFocus = d.PreviousFocuseable
		if nextFocus == nil {
			nextFocus = d.LastFocusable
		}
	} else if nextFocus == nil {
		nextFocus = d.FirstFocusable
	}

	if curFocus == nextFocus {
		return true
	}
	if curFocus != nil {
		curFocus.SetFocused(false)
	}
	if nextFocus != nil {
		nextFocus.SetFocused(true)
	}
	return true
}

// MoveFocus moves focus in the tab order
func (d *DOM) MoveFocusLegacy(direction int) bool {
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

func (c *Node) FindFocused() *Node {
	for _, child := range c.Children {
		if child == nil {
			continue
		}
		res := child.FindFocused()
		if res != nil {
			return res
		}
	}
	if c.IsFocused() {
		return c
	}
	return nil
}

func (c *Node) FindFocusable() *Node {
	for _, child := range c.Children {
		if child == nil {
			continue
		}
		res := child.FindFocusable()
		if res != nil {
			return res
		}
	}
	if c.IsFocusable() {
		return c
	}
	return nil
}

func (c *Node) FindAllFocusable() []*Node {
	var res []*Node
	for _, child := range c.Children {
		if child == nil {
			continue
		}
		res = append(res, child.FindAllFocusable()...)
	}
	if c.IsFocusable() {
		res = append(res, c)
	}
	return res
}

func (c *Node) SetFocused(focused bool) {
	if c.Props != nil {
		if focused {
			if onFocus, ok := c.Props.Get("onFocus"); ok {
				if handler, ok := onFocus.(func()); ok && handler != nil {
					handler()
				}
			}
		} else {
			if onBlur, ok := c.Props.Get("onBlur"); ok {
				if handler, ok := onBlur.(func()); ok && handler != nil {
					handler()
				}
			}
		}
	}
}

func (c *Node) IsFocused() bool {
	if c.Props != nil {
		if focused, ok := c.Props.Get("focused"); ok {
			if focusedBool, ok := focused.(bool); ok {
				return focusedBool
			}
		}
	}
	return false
}

func (c *Node) IsFocusable() bool {
	if c.Props != nil {
		if focusable, ok := c.Props.Get("focusable"); ok {
			if pbool, ok := focusable.(*bool); ok {
				if pbool == nil {
					// default value
					if c.Type == ElementTypeInput {
						return true
					}
					return false
				}
				return *pbool
			}
			if focusableBool, ok := focusable.(bool); ok {
				return focusableBool
			} else {
				panic(fmt.Errorf("focusable expect bool, got %T", focusable))
			}
		}
	}
	return false
}
