package showerrors

import (
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

// General stuff for styling the view
var (
	keywordStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	textGreenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#008000"))
	textRedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	dotStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle      = lipgloss.NewStyle().MarginLeft(2)

	quitViewStyle = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("170"))
	choiceStyle   = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("241"))
)

type Page interface {
	tea.Model
}

type ShowErrors struct {
	data *model.Data
}

func NewShowErrors(model *model.Data) Page {
	showResponse := &ShowErrors{}
	showResponse.data = model
	return showResponse
}

func (m ShowErrors) Init() tea.Cmd {
	return nil
}

// Update loop for the first view where you're choosing a task.
func (m ShowErrors) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.SetNextStepByNameExit()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ShowErrors) SetNextStepByNameExit() {

	if m.data.User.Auth {
		m.data.NextStep.NextStepByName = "submenu"
	} else {
		m.data.NextStep.NextStepByName = "topmenu"
	}
}

// The first view, where you're choosing a task
func (m ShowErrors) View() string {
	switch m.data.NextStep.ShowErrorByName {
	case "errors":
		return m.ShowError()
	}

	return ""

}
