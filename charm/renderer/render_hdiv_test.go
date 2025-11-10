package renderer

import (
	"testing"

	"github.com/xhd2015/go-dom-tui/dom"
)

// TestHDivSingleLine tests rendering hdiv with single-line children
func TestHDivSingleLine(t *testing.T) {
	t.Run("TwoSingleLineSentences", func(t *testing.T) {
		// Create two text nodes
		text1 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Hello",
		}

		text2 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "World",
		}

		// Create two div elements
		div1 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text1},
		}

		div2 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text2},
		}

		// Create hdiv with two divs
		hdiv := &dom.Node{
			Type:     dom.ElementTypeHDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{div1, div2},
		}

		// Render the hdiv
		output := RenderToStringStripColor(hdiv)

		// Expected: "Hello" and "World" should be on the same line, side by side
		// Each div adds a newline, so the output should be "HelloWorld\n"
		expected := "HelloWorld\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})

	t.Run("ThreeSingleLineSentences", func(t *testing.T) {
		// Create three text nodes
		text1 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "A",
		}

		text2 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "B",
		}

		text3 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "C",
		}

		// Create three div elements
		div1 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text1},
		}

		div2 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text2},
		}

		div3 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text3},
		}

		// Create hdiv with three divs
		hdiv := &dom.Node{
			Type:     dom.ElementTypeHDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{div1, div2, div3},
		}

		// Render the hdiv
		output := RenderToStringStripColor(hdiv)

		// Expected: "A", "B", and "C" should be on the same line, side by side
		expected := "ABC\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})
}

// TestHDivMultiLine tests rendering hdiv with multi-line children
func TestHDivMultiLine(t *testing.T) {
	t.Run("TwoMultiLineDivs", func(t *testing.T) {
		// Create first div with 3 lines
		text1 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Line1",
		}
		text2 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Line2",
		}
		text3 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Line3",
		}

		subDiv1 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text1},
		}
		subDiv2 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text2},
		}
		subDiv3 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text3},
		}

		div1 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{subDiv1, subDiv2, subDiv3},
		}

		// Create second div with 2 lines
		text4 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "A",
		}
		text5 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "B",
		}

		subDiv4 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text4},
		}
		subDiv5 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text5},
		}

		div2 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{subDiv4, subDiv5},
		}

		// Create hdiv with two divs
		hdiv := &dom.Node{
			Type:     dom.ElementTypeHDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{div1, div2},
		}

		// Render the hdiv
		output := RenderToStringStripColor(hdiv)

		// Expected: First div has 3 lines, second has 2 lines
		// With top alignment (default), the output should be:
		// Line1A
		// Line2B
		// Line3
		expected := "Line1A\nLine2B\nLine3\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})

	t.Run("ThreeMultiLineDivs", func(t *testing.T) {
		// Create first div with 2 lines
		text1 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "X",
		}
		text2 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Y",
		}

		subDiv1 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text1},
		}
		subDiv2 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text2},
		}

		div1 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{subDiv1, subDiv2},
		}

		// Create second div with 3 lines
		text3 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "1",
		}
		text4 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "2",
		}
		text5 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "3",
		}

		subDiv3 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text3},
		}
		subDiv4 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text4},
		}
		subDiv5 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text5},
		}

		div2 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{subDiv3, subDiv4, subDiv5},
		}

		// Create third div with 1 line
		text6 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Z",
		}

		subDiv6 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{text6},
		}

		div3 := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{subDiv6},
		}

		// Create hdiv with three divs
		hdiv := &dom.Node{
			Type:     dom.ElementTypeHDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{div1, div2, div3},
		}

		// Render the hdiv
		output := RenderToStringStripColor(hdiv)

		// Expected: First div has 2 lines, second has 3 lines, third has 1 line
		// With top alignment (default), the output should be:
		// X1Z
		// Y2
		//  3
		expected := "X1Z\nY2\n 3\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})

	t.Run("ComplexMultiLineWithWords", func(t *testing.T) {
		// Create first div with 3 lines (1-3 words per line)
		// Line 1: "Hello" (1 word)
		// Line 2: "World Today" (2 words)
		// Line 3: "This Is Great" (3 words)
		createTextDiv := func(text string) *dom.Node {
			textNode := &dom.Node{
				Type:  dom.ElementTypeText,
				Props: dom.NewStructProps(dom.TextNodeProps{}),
				Text:  text,
			}
			return &dom.Node{
				Type:     dom.ElementTypeDiv,
				Props:    dom.NewStructProps(dom.DivProps{}),
				Children: []*dom.Node{textNode},
			}
		}

		div1 := &dom.Node{
			Type:  dom.ElementTypeDiv,
			Props: dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{
				createTextDiv("Hello"),
				createTextDiv("World Today"),
				createTextDiv("This Is Great"),
			},
		}

		// Create second div with 4 lines (1-4 words per line)
		// Line 1: "A" (1 word)
		// Line 2: "B C" (2 words)
		// Line 3: "D E F" (3 words)
		// Line 4: "G H I J" (4 words)
		div2 := &dom.Node{
			Type:  dom.ElementTypeDiv,
			Props: dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{
				createTextDiv("A"),
				createTextDiv("B C"),
				createTextDiv("D E F"),
				createTextDiv("G H I J"),
			},
		}

		// Create third div with 5 lines (1-4 words per line, cycling)
		// Line 1: "One" (1 word)
		// Line 2: "Two Three" (2 words)
		// Line 3: "Four Five Six" (3 words)
		// Line 4: "Seven Eight Nine Ten" (4 words)
		// Line 5: "End" (1 word)
		div3 := &dom.Node{
			Type:  dom.ElementTypeDiv,
			Props: dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{
				createTextDiv("One"),
				createTextDiv("Two Three"),
				createTextDiv("Four Five Six"),
				createTextDiv("Seven Eight Nine Ten"),
				createTextDiv("End"),
			},
		}

		// Create hdiv with three divs
		hdiv := &dom.Node{
			Type:     dom.ElementTypeHDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{div1, div2, div3},
		}

		// Render the hdiv
		output := RenderToStringStripColor(hdiv)

		// Expected: The result should be 5 lines high (max of 3, 4, 5)
		// Each line concatenates the corresponding line from each div
		// Divs with fewer lines will have empty space
		// Note: Each div's lines need to be padded to the max width of that div
		// Div1 max width: "This Is Great" = 13 chars
		// Div2 max width: "G H I J" = 7 chars
		// Div3 max width: "Seven Eight Nine Ten" = 20 chars (last child, no padding needed)
		expected := "" +
			"Hello        A      One\n" +
			"World Today  B C    Two Three\n" +
			"This Is GreatD E F  Four Five Six\n" +
			"             G H I JSeven Eight Nine Ten\n" +
			"                    End\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})
}
