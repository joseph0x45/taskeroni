package tui

import "github.com/charmbracelet/lipgloss"

var (
	redTextStyle = lipgloss.NewStyle().Foreground(
		lipgloss.Color("#ff8989"),
	)

	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
)
