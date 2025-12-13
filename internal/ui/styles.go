package ui

import "github.com/charmbracelet/lipgloss"

var (
	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("62"))

	Title = lipgloss.NewStyle().
		Bold(true).
		MarginBottom(1)

	Selected = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	Item = lipgloss.NewStyle()

	ItemSelected = Item.Copy().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	Completed = lipgloss.NewStyle().
			Faint(true)
)

var (
	panelBase = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("247")).
			Padding(1, 2)

	panelActive = panelBase.Copy().
			BorderForeground(lipgloss.Color("212"))
)

func Panel(width int, active bool) lipgloss.Style {
	if width > 0 {
		if active {
			return panelActive.Copy().Width(width)
		}
		return panelBase.Copy().Width(width)
	}
	if active {
		return panelActive.Copy()
	}
	return panelBase.Copy()
}
