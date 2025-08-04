package dom

import "testing"

func TestHandleDeleteBackWord(t *testing.T) {
	tests := []struct {
		name         string
		currentValue string
		pos          int
		expectedStr  string
		expectedPos  int
	}{
		{
			name:         "delete word at end",
			currentValue: "hello world",
			pos:          11,
			expectedStr:  "hello ",
			expectedPos:  6,
		},
		{
			name:         "delete word in middle",
			currentValue: "hello world test",
			pos:          11,
			expectedStr:  "hello  test",
			expectedPos:  6,
		},
		{
			name:         "delete word with trailing spaces",
			currentValue: "hello   world",
			pos:          8,
			expectedStr:  "   world",
			expectedPos:  0,
		},
		{
			name:         "delete word at beginning",
			currentValue: "hello world",
			pos:          5,
			expectedStr:  " world",
			expectedPos:  0,
		},
		{
			name:         "delete single word",
			currentValue: "hello",
			pos:          5,
			expectedStr:  "",
			expectedPos:  0,
		},
		{
			name:         "delete from middle of word",
			currentValue: "hello world",
			pos:          8,
			expectedStr:  "hello rld",
			expectedPos:  6,
		},
		{
			name:         "position at beginning",
			currentValue: "hello world",
			pos:          0,
			expectedStr:  "hello world",
			expectedPos:  0,
		},
		{
			name:         "empty string",
			currentValue: "",
			pos:          0,
			expectedStr:  "",
			expectedPos:  0,
		},
		{
			name:         "only spaces",
			currentValue: "   ",
			pos:          3,
			expectedStr:  "   ",
			expectedPos:  0,
		},
		{
			name:         "multiple spaces between words",
			currentValue: "hello    world",
			pos:          9,
			expectedStr:  "    world",
			expectedPos:  0,
		},
		{
			name:         "word at start with spaces after",
			currentValue: "hello   ",
			pos:          8,
			expectedStr:  "   ",
			expectedPos:  0,
		},
		{
			name:         "complex case with multiple words",
			currentValue: "one two three four",
			pos:          13,
			expectedStr:  "one two  four",
			expectedPos:  8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStr, gotPos := handleDeleteBackWord(tt.currentValue, tt.pos)
			if gotStr != tt.expectedStr {
				t.Errorf("handleDeleteBackWord() gotStr = %q, want %q", gotStr, tt.expectedStr)
			}
			if gotPos != tt.expectedPos {
				t.Errorf("handleDeleteBackWord() gotPos = %d, want %d", gotPos, tt.expectedPos)
			}
		})
	}
}

func TestUpdateInputValue(t *testing.T) {
	tests := []struct {
		name         string
		currentValue string
		pos          int
		key          string
		expectedStr  string
		expectedPos  int
	}{
		// Backspace tests
		{
			name:         "backspace at middle",
			currentValue: "hello",
			pos:          3,
			key:          "backspace",
			expectedStr:  "helo",
			expectedPos:  2,
		},
		{
			name:         "backspace at beginning",
			currentValue: "hello",
			pos:          0,
			key:          "backspace",
			expectedStr:  "ello",
			expectedPos:  0,
		},
		{
			name:         "backspace at end",
			currentValue: "hello",
			pos:          5,
			key:          "backspace",
			expectedStr:  "hell",
			expectedPos:  4,
		},
		{
			name:         "backspace empty string",
			currentValue: "",
			pos:          0,
			key:          "backspace",
			expectedStr:  "",
			expectedPos:  0,
		},

		// Delete tests
		{
			name:         "delete at middle",
			currentValue: "hello",
			pos:          2,
			key:          "delete",
			expectedStr:  "helo",
			expectedPos:  2,
		},
		{
			name:         "delete at beginning",
			currentValue: "hello",
			pos:          0,
			key:          "delete",
			expectedStr:  "ello",
			expectedPos:  0,
		},
		{
			name:         "delete at end",
			currentValue: "hello",
			pos:          5,
			key:          "delete",
			expectedStr:  "hell",
			expectedPos:  5,
		},
		{
			name:         "delete empty string",
			currentValue: "",
			pos:          0,
			key:          "delete",
			expectedStr:  "",
			expectedPos:  0,
		},

		// Regular character input tests
		{
			name:         "insert char at middle",
			currentValue: "hllo",
			pos:          2,
			key:          "e",
			expectedStr:  "hello",
			expectedPos:  3,
		},
		{
			name:         "insert char at beginning",
			currentValue: "ello",
			pos:          0,
			key:          "h",
			expectedStr:  "hello",
			expectedPos:  1,
		},
		{
			name:         "insert char at end",
			currentValue: "hell",
			pos:          4,
			key:          "o",
			expectedStr:  "hello",
			expectedPos:  5,
		},
		{
			name:         "insert char in empty string",
			currentValue: "",
			pos:          0,
			key:          "a",
			expectedStr:  "a",
			expectedPos:  1,
		},

		// Navigation keys tests
		{
			name:         "enter key",
			currentValue: "hello",
			pos:          3,
			key:          "enter",
			expectedStr:  "hello",
			expectedPos:  3,
		},
		{
			name:         "tab key",
			currentValue: "hello",
			pos:          3,
			key:          "tab",
			expectedStr:  "hello",
			expectedPos:  3,
		},
		{
			name:         "escape key",
			currentValue: "hello",
			pos:          3,
			key:          "escape",
			expectedStr:  "hello",
			expectedPos:  3,
		},
		{
			name:         "arrow keys",
			currentValue: "hello",
			pos:          3,
			key:          "left",
			expectedStr:  "hello",
			expectedPos:  3,
		},

		// Multi-character special keys
		{
			name:         "ctrl+c should not change value",
			currentValue: "hello",
			pos:          3,
			key:          "ctrl+c",
			expectedStr:  "hello",
			expectedPos:  3,
		},

		// Ctrl+w tests (word deletion)
		{
			name:         "ctrl+w simple word",
			currentValue: "hello world",
			pos:          5,
			key:          "ctrl+w",
			expectedStr:  " world",
			expectedPos:  1,
		},
		{
			name:         "ctrl+w with spaces",
			currentValue: "hello   world",
			pos:          8,
			key:          "ctrl+w",
			expectedStr:  " world",
			expectedPos:  1,
		},
		{
			name:         "ctrl+w at beginning",
			currentValue: "hello world",
			pos:          0,
			key:          "ctrl+w",
			expectedStr:  "hello world",
			expectedPos:  0,
		},
		{
			name:         "ctrl+w single word",
			currentValue: "hello",
			pos:          5,
			key:          "ctrl+w",
			expectedStr:  "",
			expectedPos:  0,
		},
		{
			name:         "ctrl+w multiple words",
			currentValue: "one two three",
			pos:          7,
			key:          "ctrl+w",
			expectedStr:  "one  three",
			expectedPos:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStr, gotPos := UpdateInputValue(tt.currentValue, tt.pos, tt.key)
			if gotStr != tt.expectedStr {
				t.Errorf("UpdateInputValue() gotStr = %q, want %q", gotStr, tt.expectedStr)
			}
			if gotPos != tt.expectedPos {
				t.Errorf("UpdateInputValue() gotPos = %d, want %d", gotPos, tt.expectedPos)
			}
		})
	}
}