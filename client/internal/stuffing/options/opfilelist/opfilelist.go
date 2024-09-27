package opfileslist

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/checkdir"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

const listHeight = 14

var (
	getFilesDownload = false
	subtleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dotStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)

	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	quitViewStyle = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("170"))
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

type OpFilesList struct {
	data          *model.Data
	list          list.Model
	choice        string
	quitting      bool
	showFile      bool
	listFileByKey map[string]string
}

type Page interface {
	tea.Model
}

func NewOpFilesList(model *model.Data) Page {
	lp := initialModel(model)
	lp.data = model
	return lp
}

func (m OpFilesList) NextStepTo() {
	m.data.NextStep.NextStepByName = "files"
}

func (m OpFilesList) NextStepToGetFiles() {
	m.data.NextStep.NextStepByName = "requesthttp"
	m.data.NextStep.RequestByName = "getfiles"
}

func (m OpFilesList) NextStepToDialogCopyOrRewrite() {
	m.data.NextStep.NextStepByName = "showresponse"
	m.data.NextStep.RequestByName = "dialogfilescopyorrew"
}

func (m OpFilesList) NextStepToDialogRemoveOrNot() {
	m.data.NextStep.NextStepByName = "showresponse"
	m.data.NextStep.RequestByName = "dialogfilesdeleteornot"
}

func initialModel(model *model.Data) *OpFilesList {
	var listItem []list.Item
	listFileByKey := make(map[string]string)

	request := model.RequestHTTP["fileslist"]

	for key, val := range *request.FilesList {
		keyForItem := fmt.Sprintf("%d:%s", key, val.Filename)
		listItem = append(listItem, item(keyForItem))
		listFileByKey[keyForItem] = fmt.Sprintf("ID: %v\nFileName: %s\nSize: %d\nChunckSize: %d\nUploadDate: %s\n", val.ID, val.Filename, val.Length, val.ChunkSize, val.UploadDate.Format("02/01/2006 15:04:05"))

	}

	const defaultWidth = 20

	l := list.New(listItem, itemDelegate{}, defaultWidth, listHeight)
	l.Select(0)
	l.Title = "List files..."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &OpFilesList{list: l, listFileByKey: listFileByKey}
}

func (m OpFilesList) Init() tea.Cmd {
	return nil
}

func (m OpFilesList) showPasswordUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.showFile = false
			return m, tea.Quit
		case "ctrl+d":
			err := m.CheckFileAndSetOutputFileName()
			if err != nil {
				fmt.Println("CheckFileAndSetOutputFileName", err)
				time.Sleep(10 * time.Second)
			}
			return m, tea.Quit
		case "ctrl+r":
			request := m.data.RequestHTTP["fileslist"]
			chosedSplit := strings.Split(request.Choced, ":")
			request.Choced = chosedSplit[0]
			m.data.RequestHTTP["fileslist"] = request
			m.NextStepToDialogRemoveOrNot()

			return m, tea.Quit
		}
	}
	return m, nil
}

func (m OpFilesList) CheckFileAndSetOutputFileName() error {
	m.NextStepToGetFiles()
	getFilesDownload = true
	request := m.data.RequestHTTP["fileslist"]
	chosedSplit := strings.Split(request.Choced, ":")
	request.Choced = chosedSplit[0]

	filesList := *request.FilesList
	chocedKey, err := strconv.Atoi(request.Choced)
	if err != nil {
		fmt.Printf("strconv.Atoi(request.Choced): %v\n", err)
		time.Sleep(10 * time.Second)
		return err
	}
	selectedFile := filesList[chocedKey]
	// Открываем файл для записи
	_, outputFileName, isNotExistFile, err := checkdir.EnsureDirectoryExists("gokeeperspace", selectedFile.Filename)

	if err != nil {
		fmt.Printf("checkdir.EnsureDirectoryExists: %v\n", err)
		time.Sleep(10 * time.Second)
		return err
	}

	request.OutputFileName = outputFileName

	m.data.RequestHTTP["fileslist"] = request

	if !isNotExistFile {
		m.NextStepToDialogCopyOrRewrite()
	}

	return nil
}

func (m OpFilesList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.showFile {
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
				request := m.data.RequestHTTP["fileslist"]
				request.Choced = string(i)
				m.data.RequestHTTP["fileslist"] = request
				m.choice = string(i)
			}
			m.showFile = true
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m OpFilesList) View() string {
	var b strings.Builder

	tpl := "Hello! I'm your keeper!\n\n"

	if m.showFile {
		tpl += "%s\n\n"
		tpl += "ctrl+d: " + "dowload;"
		tpl += "\n"
		tpl += "ctrl+r: " + "remove;"
		tpl += "\n"
		tpl += "ctrl+c, q: " + "quit"

		request := m.data.RequestHTTP["fileslist"]
		fmt.Fprintf(&b, tpl, m.listFileByKey[request.Choced])
		return quitViewStyle.Render(b.String())
	}

	fmt.Fprintf(&b, "\n"+m.list.View())

	return style.SetStyleBeforeShowMenu(b.String())
}
