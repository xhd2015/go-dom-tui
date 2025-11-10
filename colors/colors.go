package colors

const (
	GREEN_SUCCESS = "#00FF00"
	RED_ERROR     = "#FF0000"

	PURPLE_PRIMARY = "#874BFD"

	GREY_TEXT = "#808080"

	DARK_RED_1 = "#a6025f"
	DARK_RED_2 = "#c20470"
	DARK_RED_3 = "#db007d"
	DARK_RED_4 = "#f00089"
	DARK_RED_5 = "#fa028f"
)

// ANSI 256 color codes for terminal
const (
	// Basic colors
	Red         = "1"  // Red - for errors
	Yellow      = "3"  // Yellow - for focused/highlighted items
	Cyan        = "6"  // Cyan - for metadata
	Grey        = "8"  // Grey - for dimmed/secondary text
	BrightWhite = "15" // Bright white - for primary text

	// Semantic names
	TextPrimary   = BrightWhite // Primary readable text
	TextSecondary = Grey        // Secondary/dimmed text
	TextError     = Red         // Error messages
	TextHighlight = Yellow      // Highlighted/focused items
	TextMetadata  = Cyan        // Metadata information
)
