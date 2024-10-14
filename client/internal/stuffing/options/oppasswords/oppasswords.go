package oppasswords

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type LoginAndPassword struct {
	data       *model.Data
	focusIndex int
	inputs     []textinput.Model
}

type Page interface {
	tea.Model
	SetUsername(username string)
	SetPassword(password string)
}

func (m LoginAndPassword) SetUsername(username string) {
	m.inputs[0].SetValue(username)
}

func (m LoginAndPassword) SetPassword(password string) {
	m.inputs[1].SetValue(password)
}

func NewOpPasswords(model *model.Data) Page {
	lipgloss.SetColorProfile(termenv.Ascii)

	lp := initialModel()
	lp.data = model

	//Если режим редактирования инициализируем поля данными
	if lp.data.OptionPasswords.Edite {
		lp.inputs[0].SetValue(lp.data.OptionPasswords.Username)
		lp.inputs[1].SetValue(lp.data.OptionPasswords.Password)
	}

	return lp
}

func (m LoginAndPassword) NextStepToBack() {
	m.data.NextStep.NextStepByName = "passwords"
}

func (m LoginAndPassword) NextStepToSendAuthDataToServer() {
	m.data.NextStep.NextStepByName = "requesthttp"
	m.data.NextStep.RequestByName = "newpassword"
}

func (m LoginAndPassword) SaveAuthorizedToModel() {
	for i := range m.inputs {
		if m.inputs[i].Placeholder == "Nickname" {
			m.data.OptionPasswords.Username = m.inputs[i].Value()
		}
		if m.inputs[i].Placeholder == "Password" {
			m.data.OptionPasswords.Password = m.inputs[i].Value()
		}
	}

}

func initialModel() LoginAndPassword {
	m := LoginAndPassword{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Nickname"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m LoginAndPassword) Init() tea.Cmd {
	return textinput.Blink
}

func (m LoginAndPassword) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.NextStepToBack()
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.SaveAuthorizedToModel()
				m.NextStepToSendAuthDataToServer()
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *LoginAndPassword) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m LoginAndPassword) View() string {
	var b strings.Builder

	tpl := "Hello! I'm your keeper?\n\n"
	fmt.Fprint(&b, tpl, fmt.Sprint("Description: ", m.data.OptionPasswords.Description), "\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	tplControl := "ctrl+c, q: quit"

	fmt.Fprintf(&b, "%s\n\n", tplControl)

	return style.SetStyleBeforeShowMenu(b.String())

}
