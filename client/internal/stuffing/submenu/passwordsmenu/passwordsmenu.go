package passwordsmenu

import (
	"fmt"
	"strings"

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

type PasswordsMenu struct {
	data     *model.Data
	Choice   int
	Chosen   bool
	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool
}

func NewPasswordsMenu(model *model.Data) Page {
	passwordsmenu := &PasswordsMenu{}
	passwordsmenu.data = model

	//Сбрасываем статус редактирования
	passwordsmenu.data.OptionPasswords.Edite = false

	return passwordsmenu
}

func (m PasswordsMenu) Init() tea.Cmd {
	return nil
}

// Update loop for the first view where you're choosing a task.
func (m PasswordsMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m PasswordsMenu) SetNextStepByNameExit() {
	m.data.NextStep.NextStepByName = "submenu"
}

func (m PasswordsMenu) SetNextStepByName() {
	if m.Choice == 0 {
		m.data.NextStep.NextStepByName = "requesthttp"
		m.data.NextStep.RequestByName = "passwordslist"
		return
	}

	if m.Choice == 1 {
		m.data.NextStep.NextStepByName = "description"
		m.data.DescriptionStep = "passwords"
		return
	}
}

// The first view, where you're choosing a task
func (m PasswordsMenu) View() string {
	var b strings.Builder

	c := m.Choice

	tpl := "Hello! I'm your keeper?\n\n"
	tpl += "%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("ctrl+c, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n",
		checkbox("List passwords", c == 0),
		checkbox("New password", c == 1),
	)

	fmt.Fprintf(&b, tpl, choices)

	return style.SetStyleBeforeShowMenu(b.String())
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
