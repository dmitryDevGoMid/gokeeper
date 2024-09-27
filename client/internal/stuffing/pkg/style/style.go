package style

import "github.com/charmbracelet/lipgloss"

var quitViewStyle = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("170"))

func SetStyleBeforeShowMenu(menu string) string {

	text := lipgloss.JoinHorizontal(lipgloss.Left, menu)

	return quitViewStyle.Render(text)
}
