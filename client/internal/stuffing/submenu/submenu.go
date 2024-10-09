package submenu

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
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
)

type Page interface {
	tea.Model
}

type SubMenu struct {
	data     *model.Data
	Choice   int
	Chosen   bool
	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool
}

func NewSubMenu(model *model.Data) Page {
	submenu := &SubMenu{}
	submenu.data = model
	return submenu
}

func (m SubMenu) Init() tea.Cmd {
	return nil
}

// Update loop for the first view where you're choosing a task.
func (m SubMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.SetNextStepByNameExit()
			return m, tea.Quit

		case "j", "down":
			m.Choice++
			if m.Choice > 2 {
				m.Choice = 2
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

func (m SubMenu) SetNextStepByNameExit() {
	m.data.NextStep.NextStepByName = "exit"
}

func (m SubMenu) SetNextStepByName() {
	if m.Choice == 0 {
		m.data.NextStep.NextStepByName = "passwords"
		return
	}

	if m.Choice == 1 {
		m.data.NextStep.NextStepByName = "cards"
		return
	}

	if m.Choice == 2 {
		m.data.NextStep.NextStepByName = "files"
		return
	}
}

// The first view, where you're choosing a task
func (m SubMenu) View() string {
	c := m.Choice

	tpl := "Hello! I'm your keeper! What to do?\n\n"
	tpl += fmt.Sprintf("id:%s\n\n", m.data.User.IDUser)
	tpl += "%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("ctrl+c, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n",
		checkbox("Passwords", c == 0),
		checkbox("Cards", c == 1),
		checkbox("Files", c == 2),
	)

	return style.SetStyleBeforeShowMenu(fmt.Sprintf(tpl, choices))
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
