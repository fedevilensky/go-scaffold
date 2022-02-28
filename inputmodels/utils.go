package inputmodels

import "github.com/charmbracelet/lipgloss"

func HelpStyle(str string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render(str)
}
