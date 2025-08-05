package styles

// Style represents styling options for div elements
type Style struct {
	Color         string // text color
	Bold          bool   // bold text
	Italic        bool   // italic text
	Underline     bool   // underline text
	Strikethrough bool   // strikethrough text

	BackgroundColor string // background color

	BorderColor  string // Color for border (empty = no border)
	BorderRouned bool

	PaddingLeft   *int
	PaddingRight  *int
	PaddingTop    *int
	PaddingBottom *int

	MarginLeft   *int
	MarginRight  *int
	MarginTop    *int
	MarginBottom *int

	FontSize int

	NoDefault bool
}

func Int(value int) *int {
	return &value
}

func String(value string) *string {
	return &value
}

func Bool(value bool) *bool {
	return &value
}
