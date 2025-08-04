package dom

// Style represents styling options for div elements
type Style struct {
	Color         string // text color
	Bold          bool   // bold text
	Italic        bool   // italic text
	Underline     bool   // underline text
	Strikethrough bool   // strikethrough text

	BackgroundColor string // background color

	BorderColor string // Color for border (empty = no border)
}
