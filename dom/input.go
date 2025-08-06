package dom

import "github.com/xhd2015/go-dom-tui/log"

// UpdateInputValue updates input value based on key press
// return new pos
func UpdateInputValue(currentValue string, pos int, e *KeydownEvent) (string, int) {
	switch e.KeyType {
	case KeyTypeBackspace:
		if len(currentValue) > 0 {
			runes := []rune(currentValue)
			if pos > len(runes) {
				return currentValue, pos - 1
			}
			if pos > 0 {
				return string(runes[:pos-1]) + string(runes[pos:]), pos - 1
			} else {
				return string(runes[1:]), 0
			}
		}
		return currentValue, pos
	case KeyTypeDelete:
		if pos >= 0 && len(currentValue) > 0 {
			runes := []rune(currentValue)
			if pos < len(runes) {
				return string(runes[:pos]) + string(runes[pos+1:]), pos
			} else {
				return string(runes[:len(runes)-1]), pos
			}
		}
		return currentValue, pos
	case KeyTypeCtrlW:
		return handleDeleteBackWord(currentValue, pos)
	case KeyTypeCtrlA:
		return currentValue, 0
	case KeyTypeCtrlE:
		return currentValue, len([]rune(currentValue))
	case KeyTypeCtrlK:
		runes := []rune(currentValue)
		if pos < len(runes) {
			return string(runes[:pos]), pos
		}
		return currentValue, pos
	case KeyTypeEnter, KeyTypeTab, KeyTypeEsc, KeyTypeUp, KeyTypeDown, KeyTypeLeft, KeyTypeRight:
		// Navigation keys don't change the value
		return currentValue, pos
	default:
		// For regular characters, append to the value
		// Handle special keys that shouldn't be added
		if !e.Alt && len(e.Runes) > 0 {
			s := string(e.Runes)
			n := len(e.Runes)
			log.Logf("UpdateInputValue len(e.Runes): %v, pos: %v, n: %v", len(e.Runes), pos, n)

			// Single character keys (letters, numbers, symbols)

			runes := []rune(currentValue)
			if pos > len(runes) {
				return currentValue + s, pos + n
			} else {
				return string(runes[:pos]) + s + string(runes[pos:]), pos + n
			}
		}
		// Multi-character keys like "ctrl+c", "shift+tab" etc. don't change value
		return currentValue, pos
	}
}

// handleDeleteBackWord deletes back a word from the current position
// returns new string and new position
func handleDeleteBackWord(currentValue string, pos int) (string, int) {
	if pos <= 0 {
		return currentValue, pos
	}
	runes := []rune(currentValue)

	// Skip trailing spaces
	if pos > len(runes) {
		pos = len(runes)
	}
	p := pos

	for p > 0 && runes[p-1] == ' ' {
		p--
	}

	// Find start of word
	l := p
	for l > 0 && runes[l-1] != ' ' {
		l--
	}

	newStr := string(runes[:l]) + string(runes[pos:])
	newPos := l

	return newStr, newPos
}
