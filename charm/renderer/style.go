package renderer

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/xhd2015/go-dom-tui/colors"
	"github.com/xhd2015/go-dom-tui/dom"
)

// CharmStyles holds various lipgloss styles for rendering
type CharmStyles struct {
	Title          lipgloss.Style
	Subtitle       lipgloss.Style
	Text           lipgloss.Style
	InputText      lipgloss.Style
	Button         lipgloss.Style
	Container      lipgloss.Style
	CompactDiv     lipgloss.Style
	NoBorderDiv    lipgloss.Style
	CompactText    lipgloss.Style
	CompactSuccess lipgloss.Style
	Input          lipgloss.Style
	Prompt         lipgloss.Style
	Success        lipgloss.Style
	Error          lipgloss.Style
}

func defaultStyles() CharmStyles {
	return CharmStyles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2).
			Margin(1, 0),
		Subtitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#F25D94")).
			Padding(0, 1).
			Margin(0, 0),
		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Margin(0, 1),
		Button: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(lipgloss.Color("#04B575")).
			Padding(0, 3).
			Margin(0, 1).
			Bold(true),
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(2, 3).
			Margin(1, 0),
		CompactDiv: lipgloss.NewStyle().
			Padding(0, 0).
			Margin(0, 0),
		NoBorderDiv: lipgloss.NewStyle().
			Padding(1, 2).
			Margin(0, 0),
		CompactText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Margin(0, 0),
		CompactSuccess: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.GREEN_SUCCESS)).
			Bold(true).
			Margin(0, 0),
		Input: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1).
			Margin(0, 1).
			Width(60),
		InputText: lipgloss.NewStyle().
			Margin(0, 1),
		Prompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Bold(true).
			Margin(0, 1),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Margin(0, 1),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Margin(0, 1),
	}
}

// dom style to charm style
func domStyleToCharmStyle(base lipgloss.Style, style dom.Style) lipgloss.Style {
	if style.Color != "" {
		base = base.Foreground(lipgloss.Color(style.Color))
	}
	if style.BackgroundColor != "" {
		base = base.Background(lipgloss.Color(style.BackgroundColor))
	}
	if style.BorderColor != "" {
		base = base.BorderForeground(lipgloss.Color(style.BorderColor))
	}
	if style.Bold {
		base = base.Bold(true)
	}
	if style.Italic {
		base = base.Italic(true)
	}
	if style.Underline {
		base = base.Underline(true)
	}
	if style.Strikethrough {
		base = base.Strikethrough(true)
	}
	return base
}
