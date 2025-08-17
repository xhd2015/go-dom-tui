package react

import (
	"fmt"
	"reflect"

	"github.com/xhd2015/go-dom-tui/dom"
)

// DiffResult represents the result of diffing two VNodes
type DiffResult struct {
	Patches []Patch
}

// Diff compares two virtual DOM trees and returns patches (like React's reconciliation)
func Diff(oldTree, newTree *dom.Node) *DiffResult {
	result := &DiffResult{
		Patches: []Patch{},
	}

	diffNodes(oldTree, newTree, []int{}, result)
	return result
}

// diffNodes recursively diffs two nodes
func diffNodes(oldNode, newNode *dom.Node, path []int, result *DiffResult) {
	// Case 1: New node doesn't exist (deletion)
	if newNode == nil {
		if oldNode != nil {
			result.Patches = append(result.Patches, Patch{
				Type:    PatchDelete,
				Path:    copyPath(path),
				OldNode: oldNode,
			})
		}
		return
	}

	// Case 2: Old node doesn't exist (creation)
	if oldNode == nil {
		result.Patches = append(result.Patches, Patch{
			Type:    PatchCreate,
			Path:    copyPath(path),
			NewNode: newNode,
		})
		return
	}

	// Case 3: Different types (replacement)
	if oldNode.Type != newNode.Type {
		result.Patches = append(result.Patches, Patch{
			Type:    PatchReplace,
			Path:    copyPath(path),
			OldNode: oldNode,
			NewNode: newNode,
		})
		return
	}

	// Case 4: Same type, check for updates
	if oldNode.Type == dom.ElementTypeText {
		// Text node - compare text content
		if oldNode.Text != newNode.Text {
			result.Patches = append(result.Patches, Patch{
				Type:    PatchUpdate,
				Path:    copyPath(path),
				OldNode: oldNode,
				NewNode: newNode,
			})
		}
	} else {
		// Element node - compare props
		propChanges := diffProps(oldNode.Props, newNode.Props)
		if len(propChanges) > 0 {
			result.Patches = append(result.Patches, Patch{
				Type:    PatchUpdate,
				Path:    copyPath(path),
				OldNode: oldNode,
				NewNode: newNode,
				Props:   propChanges,
			})
		}

		// Diff children
		diffChildren(oldNode.Children, newNode.Children, path, result)
	}
}

// diffChildren diffs child nodes with key-based reconciliation
func diffChildren(oldChildren, newChildren []*dom.Node, path []int, result *DiffResult) {
	oldLen := len(oldChildren)
	newLen := len(newChildren)
	maxLen := oldLen
	if newLen > maxLen {
		maxLen = newLen
	}

	// Simple diffing without key optimization for now
	// TODO: Implement key-based reconciliation for better performance
	for i := 0; i < maxLen; i++ {
		childPath := append(path, i)

		var oldChild, newChild *dom.Node
		if i < oldLen {
			oldChild = oldChildren[i]
		}
		if i < newLen {
			newChild = newChildren[i]
		}

		diffNodes(oldChild, newChild, childPath, result)
	}
}

// diffProps compares two property maps
func diffProps(oldProps dom.Props, newProps dom.Props) map[string]interface{} {
	changes := make(map[string]interface{})

	rangeProps(oldProps, func(key string, value any) bool {
		if newValue, exists := newProps.Get(key); !exists || !reflect.DeepEqual(value, newValue) {
			changes[key] = newValue
		}
		return true
	})

	rangeProps(newProps, func(key string, value any) bool {
		if _, exists := oldProps.Get(key); !exists {
			changes[key] = value
		}
		return true
	})

	return changes
}

func rangeProps(props dom.Props, fn func(key string, value any) bool) {
	if rangeable, ok := props.(dom.RangeableProps); ok {
		rangeable.Range(fn)
	} else {
		panic(fmt.Errorf("rangeProps: props is not a RangeableProps"))
	}
}

// copyPath creates a copy of a path slice
func copyPath(path []int) []int {
	result := make([]int, len(path))
	copy(result, path)
	return result
}
