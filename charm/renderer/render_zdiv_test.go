package renderer

import (
	"testing"

	"github.com/xhd2015/go-dom-tui/dom"
)

// TestZDivBasicOverlay tests basic zdiv overlay functionality
func TestZDivBasicOverlay(t *testing.T) {
	t.Run("TwoSingleLineTexts", func(t *testing.T) {
		// Create first text node (background layer)
		text1 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "AAAAA",
		}

		// Create second text node (foreground layer)
		text2 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "B",
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

		// Create zdiv with two divs
		zdiv := &dom.Node{
			Type:     dom.ElementTypeZDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{div1, div2},
		}

		// Render the zdiv
		output := RenderToStringStripColor(zdiv)

		// Expected: "B" should overlay the first character of "AAAAA"
		expected := "BAAAA\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})

	t.Run("ThreeSingleLineTexts", func(t *testing.T) {
		// Create three text nodes
		text1 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "XXXXX",
		}

		text2 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "YY",
		}

		text3 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "Z",
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

		// Create zdiv with three divs
		zdiv := &dom.Node{
			Type:     dom.ElementTypeZDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{div1, div2, div3},
		}

		// Render the zdiv
		output := RenderToStringStripColor(zdiv)

		// Expected: "Z" overlays first char, "YY" overlays first two chars, "XXXXX" is base
		// Result: "Z" at pos 0, "Y" at pos 1 (from YY), "XXX" at pos 2-4 (from XXXXX)
		expected := "ZYXXX\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})
}

// TestZDivMultiLineOverlay tests zdiv with multi-line children
func TestZDivMultiLineOverlay(t *testing.T) {
	t.Run("TwoMultiLineDivs", func(t *testing.T) {
		// Create first div with 3 lines (background)
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

		// Create second div with 2 lines (foreground)
		text4 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "AAA",
		}
		text5 := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "BBB",
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

		// Create zdiv with two divs
		zdiv := &dom.Node{
			Type:     dom.ElementTypeZDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{div1, div2},
		}

		// Render the zdiv
		output := RenderToStringStripColor(zdiv)

		// Expected: Second div overlays first two lines of first div
		// Line 1: "AAA" overlays "Line1" -> "AAAe1"
		// Line 2: "BBB" overlays "Line2" -> "BBBe2"
		// Line 3: "Line3" (no overlay)
		expected := "AAAe1\nBBBe2\nLine3\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})

	t.Run("ThreeMultiLineDivs", func(t *testing.T) {
		// Create first div with 3 lines (bottom layer)
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
				createTextDiv("XXXXX"),
				createTextDiv("XXXXX"),
				createTextDiv("XXXXX"),
			},
		}

		// Create second div with 2 lines (middle layer)
		div2 := &dom.Node{
			Type:  dom.ElementTypeDiv,
			Props: dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{
				createTextDiv("YYY"),
				createTextDiv("YYY"),
			},
		}

		// Create third div with 1 line (top layer)
		div3 := &dom.Node{
			Type:  dom.ElementTypeDiv,
			Props: dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{
				createTextDiv("Z"),
			},
		}

		// Create zdiv with three divs
		zdiv := &dom.Node{
			Type:     dom.ElementTypeZDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{div1, div2, div3},
		}

		// Render the zdiv
		output := RenderToStringStripColor(zdiv)

		// Expected: Layers stack on top of each other
		// Line 1: "Z" at pos 0, "YY" at pos 1-2 (from YYY), "XX" at pos 3-4 (from XXXXX) -> "ZYYXX"
		// Line 2: "YYY" overlays "XXXXX" -> "YYYXX"
		// Line 3: "XXXXX" (no overlay)
		expected := "ZYYXX\nYYYXX\nXXXXX\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})
}

// TestZDivComplexOverlay tests complex overlay scenarios
func TestZDivComplexOverlay(t *testing.T) {
	t.Run("DialogOverBackground", func(t *testing.T) {
		// Simulate a background with content
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

		// Background: 5 lines of content
		background := &dom.Node{
			Type:  dom.ElementTypeDiv,
			Props: dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{
				createTextDiv("Background Line 1"),
				createTextDiv("Background Line 2"),
				createTextDiv("Background Line 3"),
				createTextDiv("Background Line 4"),
				createTextDiv("Background Line 5"),
			},
		}

		// Dialog: 3 lines centered-ish (starts at line 2)
		// We'll use a div with empty first line to position it
		emptyText := &dom.Node{
			Type:  dom.ElementTypeText,
			Props: dom.NewStructProps(dom.TextNodeProps{}),
			Text:  "",
		}
		emptyDiv := &dom.Node{
			Type:     dom.ElementTypeDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{emptyText},
		}

		dialog := &dom.Node{
			Type:  dom.ElementTypeDiv,
			Props: dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{
				emptyDiv, // Empty line to push dialog down
				createTextDiv("+----------+"),
				createTextDiv("| Dialog!! |"),
				createTextDiv("+----------+"),
			},
		}

		// Create zdiv with background and dialog
		zdiv := &dom.Node{
			Type:     dom.ElementTypeZDiv,
			Props:    dom.NewStructProps(dom.DivProps{}),
			Children: []*dom.Node{background, dialog},
		}

		// Render the zdiv
		output := RenderToStringStripColor(zdiv)

		// Expected: Dialog overlays lines 2-4 of background
		// Note: The space in "| Dialog!! |" is preserved, but "Background Line 3" starts with "B"
		// So we get "|" + " Dialog!! " + "|" overlaying "Background Line 3"
		// Position 0: "|" overlays "B"
		// Position 1: " " (space) is skipped, so "a" from "Background" shows through
		// Position 2-11: "Dialog!! " overlays "ckground L"
		// Position 12: "|" overlays "i"
		// Position 13+: "ne 3" from background
		expected := "" +
			"Background Line 1\n" +
			"+----------+ine 2\n" +
			"|aDialog!! |ine 3\n" +
			"+----------+ine 4\n" +
			"Background Line 5\n"

		// Check exact string match
		if output != expected {
			t.Errorf("Expected exact output:\n%q\nGot:\n%q", expected, output)
		}
	})
}
