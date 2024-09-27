package oppaswordslist

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

type OpPasswordsList struct {
	data              *model.Data
	list              list.Model
	choice            string
	quitting          bool
	showPassword      bool
	listPasswordByKey map[string]string
}

type Page interface {
	tea.Model
}

func NewOpPasswordsList(model *model.Data) Page {
	lp := initialModel(model)
	lp.data = model
	return lp
}

func (m OpPasswordsList) NextStepTo() {
	m.data.NextStep.NextStepByName = "passwords"
}

func (m OpPasswordsList) NextStepToEditPassword() {
	m.data.NextStep.NextStepByName = "description"
	m.data.DescriptionStep = "passwords"
}

func (m OpPasswordsList) NextStepToDeletePassword() {
	m.data.NextStep.NextStepByName = "showresponse"
	m.data.NextStep.RequestByName = "dialogpassworddeleteornot"
}

func initialModel(model *model.Data) *OpPasswordsList {
	var listItem []list.Item
	listPasswordByKey := make(map[string]string)

	request := model.RequestHTTP["passwordslist"]

	for key, val := range *request.PasswordsList {
		keyForItem := fmt.Sprintf("%d:%s", key, val.Description)
		listItem = append(listItem, item(keyForItem))
		listPasswordByKey[keyForItem] = fmt.Sprintf("ID: %v\n\nDescription: %s\nLogin: %s\nPassword: %s", val.ID, val.Description, val.Username, val.Password)
	}

	const defaultWidth = 20

	l := list.New(listItem, itemDelegate{}, defaultWidth, listHeight)
	l.Select(2)
	l.Title = "List passwords..."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &OpPasswordsList{list: l, listPasswordByKey: listPasswordByKey}
}

func (m OpPasswordsList) Init() tea.Cmd {
	return nil
}

func (m OpPasswordsList) showPasswordUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.showPassword = false
			return m, tea.Quit
		case "ctrl+e":
			err := m.SetDataForEditePassword()
			if err != nil {
				fmt.Println("SetDataForEditeFile", err)
				time.Sleep(10 * time.Second)
			}
			m.showPassword = false

			//Уходим на редактирование файла
			m.NextStepToEditPassword()
			return m, tea.Quit
		case "ctrl+d":
			err := m.SetDataForDeletePassword()
			if err != nil {
				fmt.Println("SetDataForEditePassword", err)
				time.Sleep(10 * time.Second)
			}
			m.showPassword = false

			//Уходим на редактирование файла
			m.NextStepToDeletePassword()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m OpPasswordsList) SetDataForDeletePassword() error {
	request := m.data.RequestHTTP["passwordslist"]
	chosedSplit := strings.Split(request.Choced, ":")
	request.Choced = chosedSplit[0]

	passwordsList := *request.PasswordsList
	chocedKey, err := strconv.Atoi(request.Choced)
	if err != nil {
		fmt.Printf("strconv.Atoi(request.Choced): %v\n", err)
		time.Sleep(10 * time.Second)
		return err
	}
	m.data.OptionPasswords = passwordsList[chocedKey]
	m.data.RequestHTTP["passwordslist"] = request

	return nil
}

func (m OpPasswordsList) SetDataForEditePassword() error {
	request := m.data.RequestHTTP["passwordslist"]
	chosedSplit := strings.Split(request.Choced, ":")
	request.Choced = chosedSplit[0]

	passwordsList := *request.PasswordsList
	chocedKey, err := strconv.Atoi(request.Choced)
	if err != nil {
		fmt.Printf("strconv.Atoi(request.Choced): %v\n", err)
		time.Sleep(10 * time.Second)
		return err
	}
	m.data.OptionPasswords = passwordsList[chocedKey]
	m.data.RequestHTTP["passwordslist"] = request
	m.data.OptionPasswords.Edite = true

	return nil
}

func (m OpPasswordsList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.showPassword {
		return m.showPasswordUpdate(msg)
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
				request := m.data.RequestHTTP["passwordslist"]
				request.Choced = string(i)
				m.data.RequestHTTP["passwordslist"] = request
				m.choice = string(i)
			}
			m.showPassword = true
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m OpPasswordsList) View() string {
	var b strings.Builder

	tpl := "Hello! I'm your keeper! You are viewing the password.\n\n"

	if m.showPassword {
		/*tpl += "%s\n\n"
		tpl += dotStyle + subtleStyle.Render(textBlackStyle.Render("ctrl+e: ")+textGreenStyle.Render("edite;"))
		tpl += "\n"
		tpl += dotStyle + subtleStyle.Render(textBlackStyle.Render("ctrl+d: ")+textRedStyle.Render("delete;"))
		tpl += "\n"
		tpl += dotStyle + subtleStyle.Render(textBlackStyle.Render("ctrl+c, q: ")+"quit")*/

		tpl += "%s\n\n"
		tpl += "ctrl+e: " + "edite;"
		tpl += "\n"
		tpl += "ctrl+d: " + "delete;"
		tpl += "\n"
		tpl += "ctrl+c, q: " + "quit"

		request := m.data.RequestHTTP["passwordslist"]
		fmt.Fprintf(&b, tpl, m.listPasswordByKey[request.Choced])
		return quitViewStyle.Render(b.String())

	}

	fmt.Fprint(&b, tpl, m.list.View())

	return style.SetStyleBeforeShowMenu(b.String())
}
