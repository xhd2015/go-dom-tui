package dom

import "fmt"

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
					if c.Type == "input" {
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
