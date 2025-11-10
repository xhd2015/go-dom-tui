package dom

import (
	"github.com/xhd2015/go-dom-tui/styles"
)

// CreateNode creates a virtual element (like React.createElement)
func CreateNode(typ string, props Props, children ...*Node) *Node {
	return &Node{
		Type:     typ,
		Props:    props,
		Children: children,
		Key:      GetStringProp(props, "key"),
	}
}

func CreateComponent(component Component, props Props, children ...*Node) *Node {
	panic("TODO: CreateComponent")
}

// Text creates a text node
func Text(text string, styleList ...styles.Style) *Node {
	style := styles.Style{}
	if len(styleList) > 0 {
		if len(styleList) > 1 {
			panic("only one style is allowed")
		}
		style = styleList[0]
	}
	return &Node{
		Type:  ElementTypeText,
		Props: NewStructProps(TextNodeProps{Style: style}),
		Text:  text,
	}
}
func TextWithProps(text string, props TextNodeProps) *Node {
	return &Node{
		Type:  ElementTypeText,
		Props: NewStructProps(props),
		Text:  text,
	}
}

// Div creates a div component
func Div(props DivProps, children ...*Node) *Node {
	return CreateNode(ElementTypeDiv, NewStructProps(props), children...)
}

// HDiv creates a horizontal div component that places children left to right
func HDiv(props DivProps, children ...*Node) *Node {
	return CreateNode(ElementTypeHDiv, NewStructProps(props), children...)
}

// Span creates a span component
func Span(props DivProps, children ...*Node) *Node {
	return CreateNode(ElementTypeSpan, NewStructProps(props), children...)
}

// Br creates a line break component
func Br() *Node {
	return CreateNode(ElementTypeBr, NewStructProps(EmptyProps{}))
}

func H1(props DivProps, children ...*Node) *Node {
	return CreateNode(ElementTypeH1, NewStructProps(props), children...)
}

func H2(props DivProps, children ...*Node) *Node {
	return CreateNode(ElementTypeH2, NewStructProps(props), children...)
}

func P(props DivProps, children ...*Node) *Node {
	return CreateNode(ElementTypeP, NewStructProps(props), children...)
}

func Input(props InputProps, children ...*Node) *Node {
	return CreateNode(ElementTypeInput, NewStructProps(props), children...)
}

func Button(props ButtonProps, children ...*Node) *Node {
	return CreateNode(ElementTypeButton, NewStructProps(props), children...)
}

func Ul(props DivProps, children ...*Node) *Node {
	return CreateNode(ElementTypeUl, NewStructProps(props), children...)
}

func Li(props ListItemProps, children ...*Node) *Node {
	return CreateNode(ElementTypeLi, NewStructProps(props), children...)
}

func Fragment(children ...*Node) *Node {
	return CreateNode(ElementTypeFragment, NewStructProps(EmptyProps{}), children...)
}

// Spacer creates a spacer component that expands to fill available horizontal space
func Spacer() *Node {
	return CreateNode(ElementTypeSpacer, NewStructProps(SpacerProps{MinSize: 1}))
}
