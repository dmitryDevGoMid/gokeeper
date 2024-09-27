package opdescritption

// A simple program demonstrating the textarea component from the Bubbles
// component library.

import (
	"fmt"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type errMsg error

type Page interface {
	tea.Model
}

type Description struct {
	data       *model.Data
	focusIndex int
	textarea   textarea.Model
	err        error
}

func NewOpDescription(model *model.Data, defaultDescription string) *Description {
	description := initialModel()
	description.data = model

	if defaultDescription != "" {
		description.textarea.SetValue(defaultDescription)
	}

	//Заполняем поле для редактирования пароля
	if description.data.OptionPasswords.Edite {
		description.textarea.SetValue(description.data.OptionPasswords.Description)
	}
	//Заполняем поля для редактирования карты
	if description.data.OptionCards.Edite {
		description.textarea.SetValue(description.data.OptionCards.Description)
	}

	return description
}

func initialModel() *Description {
	lipgloss.SetColorProfile(termenv.Ascii)

	ti := textarea.New()
	ti.Placeholder = "Description..."
	ti.Focus()
	ti.ShowLineNumbers = false

	return &Description{
		textarea: ti,
		err:      nil,
	}
}

func (m Description) NextStepTo() {
	switch m.data.DescriptionStep {
	case "passwords":
		m.data.NextStep.NextStepByName = "newpasswords"
	case "files":
		m.data.NextStep.NextStepByName = "selectfiles"
	case "creditcards":
		m.data.NextStep.NextStepByName = "newcreditcards"
	}
}

func (m Description) BackStepTo() {
	switch m.data.DescriptionStep {
	case "passwords":
		m.data.NextStep.NextStepByName = "passwords"
	case "files":
		m.data.NextStep.NextStepByName = "files"
	case "creditcards":
		m.data.NextStep.NextStepByName = "cards"
	}
}

// Сохраняем обписание в зависимости от модели которая была выбрана
func (m Description) SaveDescritpionToModel() {
	switch m.data.DescriptionStep {
	case "passwords":
		m.data.OptionPasswords.Description = m.textarea.Value()
	case "files":
		m.data.OptionFiles.Description = m.textarea.Value()
	case "creditcards":
		m.data.OptionCards.Description = m.textarea.Value()
	}
}

func (m Description) Init() tea.Cmd {
	return textarea.Blink
}

func (m Description) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			m.BackStepTo()
			return m, tea.Quit
		case tea.KeyEnter:
			if m.focusIndex == 1 {
				m.SaveDescritpionToModel()
				m.NextStepTo()
				return m, tea.Quit
			}
		case tea.KeyTab:
			m.focusIndex++
			if m.textarea.Focused() {
				m.textarea.Blur()
			} else {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

		//Сбрасываем счетчик для перевода фокуса на элементы
		if m.focusIndex > 1 {
			m.focusIndex = 0
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Description) View() string {
	var b strings.Builder

	button := &blurredButton
	if m.focusIndex == 1 {
		button = &focusedButton
	}
	b.WriteString("Description for login and password.\n\n")
	b.WriteString(m.textarea.View())
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	b.WriteString("(ctrl+c to quit)")
	b.WriteString("\n\n")

	return style.SetStyleBeforeShowMenu(b.String())
}
