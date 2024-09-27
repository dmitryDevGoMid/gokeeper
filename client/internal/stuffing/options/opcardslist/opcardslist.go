package opcardslist

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

/*const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)*/

const listHeight = 14

var (
	/*subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dotStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)

	textGreenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#008000"))
	textRedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	textBlackStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4B0082"))*/

	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	//quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	//choiceStyle   = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("241"))
	quitViewStyle = lipgloss.NewStyle().Padding(2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("170"))
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type OpCardsList struct {
	data          *model.Data
	list          list.Model
	choice        string
	quitting      bool
	showCard      bool
	listCardByKey map[string]string
}

type Page interface {
	tea.Model
}

func NewOpCardsList(model *model.Data) Page {
	lp := initialModel(model)
	lp.data = model
	return lp
}

func (m OpCardsList) NextStepTo() {
	m.data.NextStep.NextStepByName = "cards"
}

func (m OpCardsList) NextStepToEditCards() {
	m.data.NextStep.NextStepByName = "description"
	m.data.DescriptionStep = "creditcards"
}

func (m OpCardsList) NextStepToDeleteCards() {
	m.data.NextStep.NextStepByName = "showresponse"
	m.data.NextStep.RequestByName = "dialogcarddeleteornot"
}

func initialModel(model *model.Data) *OpCardsList {
	if model.ModeTest {
		lipgloss.SetColorProfile(termenv.Ascii)
	}
	var listItem []list.Item
	listCardByKey := make(map[string]string)

	request := model.RequestHTTP["cardslist"]

	for key, val := range *request.CardsList {
		keyForItem := fmt.Sprintf("%d:%s", key, val.Description)
		listItem = append(listItem, item(keyForItem))
		listCardByKey[keyForItem] = fmt.Sprintf("ID: %v\n\nDescription: %s\nNumber: %s\nExp: %s\nCVC: %s", val.ID, val.Description, val.Number, val.Exp, val.Cvc)
	}

	const defaultWidth = 20

	l := list.New(listItem, itemDelegate{}, defaultWidth, listHeight)
	l.Select(0)
	l.Title = "List cards..."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &OpCardsList{list: l, listCardByKey: listCardByKey}
}

func (m OpCardsList) Init() tea.Cmd {
	return nil
}

func (m OpCardsList) showCardsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.showCard = false
			return m, tea.Quit
		case "ctrl+e":
			err := m.SetDataForEditeCards()
			if err != nil {
				fmt.Println("SetDataForEditeFile", err)
				time.Sleep(10 * time.Second)
			}
			m.showCard = false

			//Уходим на редактирование файла
			m.NextStepToEditCards()
			return m, tea.Quit
		case "ctrl+d":
			err := m.SetDataForDeleteCards()
			if err != nil {
				fmt.Println("SetDataForEditeFile", err)
				time.Sleep(10 * time.Second)
			}
			m.showCard = false

			//Уходим на редактирование файла
			m.NextStepToDeleteCards()
			return m, tea.Quit

		}
	}
	return m, nil
}

func (m OpCardsList) SetDataForDeleteCards() error {
	request := m.data.RequestHTTP["cardslist"]
	chosedSplit := strings.Split(request.Choced, ":")
	request.Choced = chosedSplit[0]

	cardsList := *request.CardsList
	chocedKey, err := strconv.Atoi(request.Choced)
	if err != nil {
		fmt.Printf("strconv.Atoi(request.Choced): %v\n", err)
		time.Sleep(10 * time.Second)
		return err
	}
	m.data.OptionCards = cardsList[chocedKey]
	m.data.RequestHTTP["cardslist"] = request

	return nil
}

func (m OpCardsList) SetDataForEditeCards() error {
	request := m.data.RequestHTTP["cardslist"]
	chosedSplit := strings.Split(request.Choced, ":")
	request.Choced = chosedSplit[0]

	cardsList := *request.CardsList
	chocedKey, err := strconv.Atoi(request.Choced)
	if err != nil {
		fmt.Printf("strconv.Atoi(request.Choced): %v\n", err)
		time.Sleep(10 * time.Second)
		return err
	}
	m.data.OptionCards = cardsList[chocedKey]
	m.data.RequestHTTP["cardslist"] = request
	m.data.OptionCards.Edite = true

	return nil
}

func (m OpCardsList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.showCard {
		return m.showCardsUpdate(msg)
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			m.NextStepTo()
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				request := m.data.RequestHTTP["cardslist"]
				request.Choced = string(i)
				m.data.RequestHTTP["cardslist"] = request
				m.choice = string(i)
			}
			m.showCard = true
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m OpCardsList) View() string {
	var b strings.Builder

	tpl := "Hello! I'm your keeper! You are viewing the credit card.\n\n"

	if m.showCard {
		/*tpl += "%s\n\n"
		tpl += dotStyle + subtleStyle.Render(textBlackStyle.Render("ctrl+e: ")+textGreenStyle.Render("edite;"))
		tpl += "\n"
		tpl += dotStyle + subtleStyle.Render(textBlackStyle.Render("ctrl+d: ")+textRedStyle.Render("delete;"))
		tpl += "\n"
		tpl += dotStyle + subtleStyle.Render(textBlackStyle.Render("ctrl+c, q: ")+"quit")*/

		tpl += "%s\n\n"
		tpl += "ctrl+e: " + "edite;"
		tpl += "\n"
		tpl += "ctrl+d: " + "delettttttttte;"
		tpl += "\n"
		tpl += "ctrl+c, q: " + "quit"

		request := m.data.RequestHTTP["cardslist"]
		fmt.Fprintf(&b, tpl, m.listCardByKey[request.Choced])
		return quitViewStyle.Render(b.String())

	}

	fmt.Fprint(&b, tpl, m.list.View())

	return style.SetStyleBeforeShowMenu(b.String())
}
