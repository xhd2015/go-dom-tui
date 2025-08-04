package dom

// handleDeleteBackWord deletes back a word from the current position
// returns new string and new position
func handleDeleteBackWord(currentValue string, pos int) (string, int) {
	if pos <= 0 {
		return currentValue, pos
	}

	p := pos
	// Skip trailing spaces
	for p > 0 && currentValue[p-1] == ' ' {
		p--
	}

	// Find start of word
	l := p
	for l > 0 && currentValue[l-1] != ' ' {
		l--
	}

	newStr := currentValue[:l] + currentValue[pos:]
	newPos := l

	return newStr, newPos
}

// UpdateInputValue updates input value based on key press
// return new pos
func UpdateInputValue(currentValue string, pos int, key string) (string, int) {
	switch key {
	case "backspace":
		if len(currentValue) > 0 {
			if pos > len(currentValue) {
				return currentValue, pos - 1
			}
			if pos > 0 {
				return currentValue[:pos-1] + currentValue[pos:], pos - 1
			} else {
				return currentValue[1:], 0
			}
		}
		return currentValue, pos
	case "delete":
		if len(currentValue) > 0 {
			if pos < len(currentValue) {
				return currentValue[:pos] + currentValue[pos+1:], pos
			} else {
				return currentValue[:len(currentValue)-1], pos
			}
		}
		return currentValue, pos
	case "ctrl+w":
		return handleDeleteBackWord(currentValue, pos)
	case "ctrl+a":
		return currentValue, 0
	case "ctrl+e":
		return currentValue, len(currentValue)
	case "ctrl+k":
		if pos < len(currentValue) {
			return currentValue[:pos], pos
		}
		return currentValue, pos
	case "enter", "tab", "escape", "up", "down", "left", "right":
		// Navigation keys don't change the value
		return currentValue, pos
	default:
		// For regular characters, append to the value
		// Handle special keys that shouldn't be added
		if len(key) == 1 {
			// Single character keys (letters, numbers, symbols)
			if pos > len(currentValue) {
				return currentValue + key, pos + 1
			} else {
				return currentValue[:pos] + key + currentValue[pos:], pos + 1
			}
		}
		// Multi-character keys like "ctrl+c", "shift+tab" etc. don't change value
		return currentValue, pos
	}
}
