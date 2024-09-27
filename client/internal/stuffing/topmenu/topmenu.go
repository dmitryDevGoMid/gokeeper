package topmenu

import (
	"fmt"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"

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
	//keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	//ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	//progressEmpty = subtleStyle.Render(progressEmptyChar)
	dotStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	//mainStyle = lipgloss.NewStyle().MarginLeft(2)

	//quitViewStyle = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("170"))
	//choiceStyle   = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("241"))

	// Gradient colors we'll use for the progress bar
	//ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
)

type Page interface {
	tea.Model
}

type TopMenu struct {
	data     *model.Data
	Choice   int
	Chosen   bool
	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool
}

func NewTopMenu(model *model.Data) Page {
	topmenu := &TopMenu{}
	topmenu.data = model
	return topmenu
}

func (m TopMenu) Init() tea.Cmd {
	return nil
}

// Update loop for the first view where you're choosing a task.
func (m TopMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.SetNextStepByNameExit()
			return m, tea.Quit

		case "j", "down":
			m.Choice++
			if m.Choice > 1 {
				m.Choice = 1
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			m.SetNextStepByName()
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m TopMenu) SetNextStepByNameExit() {
	m.data.NextStep.NextStepByName = "exit"
}

func (m TopMenu) SetNextStepByName() {
	if m.Choice == 0 {
		m.data.NextStep.NextStepByName = "login"
		return
	}

	if m.Choice == 1 {
		m.data.NextStep.NextStepByName = "registration"
		return
	}
}

// The first view, where you're choosing a task
func (m TopMenu) View() string {
	c := m.Choice

	tpl := "Hello! I'm your keeper?\n\n"
	tpl += "%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("ctrl+c, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n",
		checkbox("Login", c == 0),
		checkbox("Registrate", c == 1),
	)

	return style.SetStyleBeforeShowMenu(fmt.Sprintf(tpl, choices))
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
