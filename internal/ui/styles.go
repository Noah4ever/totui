package ui

import "github.com/charmbracelet/lipgloss"

var (
	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("62"))

	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("250")).
		Underline(true).
		MarginBottom(1)

	Item = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	ItemSelected = Item.Copy().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("229")) // high contrast, colorblind-friendly

	Completed = lipgloss.NewStyle().
			Faint(true)

	Muted = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))
)

var (
	panelBase = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("245")).
			Padding(1, 2)

	panelActive = panelBase.Copy().
			BorderForeground(lipgloss.Color("212"))
)

func Panel(width, height int, active bool) lipgloss.Style {
	style := panelBase
	if active {
		style = panelActive
	}
	if width > 0 {
		style = style.Copy().Width(width)
	}
	if height > 0 {
		style = style.Copy().Height(height)
	}
	return style
}
