package dom

// Window represents the global terminal window state
type Window struct {
	Width  int
	Height int
}

// Update sets the current window dimensions
func (w *Window) Update(width, height int) {
	w.Width = width
	w.Height = height
}

// Get returns the current window dimensions
func (w *Window) Get() (int, int) {
	return w.Width, w.Height
}

// GetWidth returns the current window width
func (w *Window) GetWidth() int {
	return w.Width
}

// GetHeight returns the current window height
func (w *Window) GetHeight() int {
	return w.Height
}
