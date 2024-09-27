package opcreditcard

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

const (
	ccn = iota
	exp
	cvv
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

/*type model struct {
	inputs  []textinput.Model
	focused int
	err     error
}*/

type CreditCard struct {
	data    *model.Data
	inputs  []textinput.Model
	focused int
	err     error
}

type Page interface {
	tea.Model
}

func NewCreditCard(model *model.Data) Page {
	creditCard := initialModel()
	creditCard.data = model
	//Если режим редактирования инициализируем поля данными
	if creditCard.data.OptionCards.Edite {
		creditCard.inputs[0].SetValue(creditCard.data.OptionCards.Number)
		creditCard.inputs[1].SetValue(creditCard.data.OptionCards.Exp)
		creditCard.inputs[2].SetValue(creditCard.data.OptionCards.Cvc)
	}

	return creditCard
}

func (m CreditCard) NextStepBack() {
	m.data.NextStep.NextStepByName = "cards"
}

func (m CreditCard) NextStepSendToServer() {
	m.data.NextStep.NextStepByName = "requesthttp"
	m.data.NextStep.RequestByName = "newcreditcard"
}

func (m CreditCard) SaveCardToModel() {
	m.data.OptionCards.Number = m.inputs[0].Value()
	m.data.OptionCards.Exp = m.inputs[1].Value()
	m.data.OptionCards.Cvc = m.inputs[2].Value()

}

// Validator functions to ensure valid input
func ccnValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	// The last digit should be a number unless it is a multiple of 4 in which
	// case it should be a space
	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func expValidator(s string) error {
	// The 3 character should be a slash (/)
	// The rest should be numbers
	e := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return fmt.Errorf("EXP is invalid")
	}

	// There should be only one slash and it should be in the 2nd index (3rd character)
	if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("EXP is invalid")
	}

	return nil
}

func cvvValidator(s string) error {
	// The CVV should be a number of 3 digits
	// Since the input will already ensure that the CVV is a string of length 3,
	// All we need to do is check that it is a number
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func initialModel() CreditCard {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[ccn] = textinput.New()
	inputs[ccn].Placeholder = "4505 **** **** 1234"
	inputs[ccn].Focus()
	inputs[ccn].CharLimit = 20
	inputs[ccn].Width = 30
	inputs[ccn].Prompt = ""
	inputs[ccn].Validate = ccnValidator

	inputs[exp] = textinput.New()
	inputs[exp].Placeholder = "MM/YY "
	inputs[exp].CharLimit = 5
	inputs[exp].Width = 5
	inputs[exp].Prompt = ""
	inputs[exp].Validate = expValidator

	inputs[cvv] = textinput.New()
	inputs[cvv].Placeholder = "XXX"
	inputs[cvv].CharLimit = 3
	inputs[cvv].Width = 5
	inputs[cvv].Prompt = ""
	inputs[cvv].Validate = cvvValidator

	return CreditCard{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m CreditCard) Init() tea.Cmd {
	return textinput.Blink
}

func (m CreditCard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == 3 {
				m.NextStepSendToServer()
				m.SaveCardToModel()
				return m, tea.Quit
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			m.NextStepBack()
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		if len(m.inputs) != m.focused {
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			m.inputs[m.focused].Focus()
		} else {
			m.inputs[m.focused-1].Blur()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m CreditCard) View() string {
	var b strings.Builder
	fmt.Fprintf(&b,
		`

 %s
 %s

 %s  %s
 %s  %s
`, inputStyle.Width(30).Render("Card Number"), m.inputs[ccn].View(), inputStyle.Width(6).Render("EXP"),
		inputStyle.Width(6).Render("CVV"), m.inputs[exp].View(), m.inputs[cvv].View())

	button := &blurredButton
	if m.focused == len(m.inputs) {
		button = &focusedButton
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	fmt.Fprintf(&b, "%d", m.focused)

	return style.SetStyleBeforeShowMenu(b.String())
}

// nextInput focuses the next input field
func (m *CreditCard) nextInput() {
	m.focused = (m.focused + 1) % (len(m.inputs) + 1)
}

// prevInput focuses the previous input field
func (m *CreditCard) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs)
	}
}
