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
		Type:  "text",
		Props: NewStructProps(TextNodeProps{Style: style}),
		Text:  text,
	}
}
func TextWithProps(text string, props TextNodeProps) *Node {
	return &Node{
		Type:  "text",
		Props: NewStructProps(props),
		Text:  text,
	}
}

// Div creates a div component
func Div(props DivProps, children ...*Node) *Node {
	return CreateNode("div", NewStructProps(props), children...)
}

// Span creates a span component
func Span(props DivProps, children ...*Node) *Node {
	return CreateNode("span", NewStructProps(props), children...)
}

// Br creates a line break component
func Br() *Node {
	return CreateNode("br", NewStructProps(EmptyProps{}))
}

func H1(props DivProps, children ...*Node) *Node {
	return CreateNode("h1", NewStructProps(props), children...)
}

func H2(props DivProps, children ...*Node) *Node {
	return CreateNode("h2", NewStructProps(props), children...)
}

func P(props DivProps, children ...*Node) *Node {
	return CreateNode("p", NewStructProps(props), children...)
}

func Input(props InputProps, children ...*Node) *Node {
	return CreateNode("input", NewStructProps(props), children...)
}

func Button(props ButtonProps, children ...*Node) *Node {
	return CreateNode("button", NewStructProps(props), children...)
}

func Ul(props DivProps, children ...*Node) *Node {
	return CreateNode("ul", NewStructProps(props), children...)
}

func Li(props ListItemProps, children ...*Node) *Node {
	return CreateNode("li", NewStructProps(props), children...)
}

func Fragment(children ...*Node) *Node {
	return CreateNode("fragment", NewStructProps(EmptyProps{}), children...)
}
