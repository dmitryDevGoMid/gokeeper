package showresponse

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/keeperlog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

var (
	maxFocusIndex = 0
)

// General stuff for styling the view
var (
	keywordStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	textGreenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#008000"))
	textRedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

	dotStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle = lipgloss.NewStyle().MarginLeft(2)

	quitViewStyle = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("170"))
	choiceStyle   = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("241"))

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

var (
	focusedButtonCancel = focusedStyle.Render("[ Cancel ]")
	blurredButtonCancel = fmt.Sprintf("[ %s ]", blurredStyle.Render("Cancel"))
)

var (
	focusedButtonCopy    = focusedStyle.Render("[ Copy ]")
	blurredButtonCopy    = fmt.Sprintf("[ %s ]", blurredStyle.Render("Copy"))
	focusedButtonRewrite = focusedStyle.Render("[ Rewrite ]")
	blurredButtonRewrite = fmt.Sprintf("[ %s ]", blurredStyle.Render("Rewrite"))
)

var (
	focusedButtonDeleteFileYes = focusedStyle.Render("[ Yes ]")
	blurredButtonDeleteFileYes = fmt.Sprintf("[ %s ]", blurredStyle.Render("Yes"))
	focusedButtonDeleteFileNo  = focusedStyle.Render("[ No ]")
	blurredButtonDeleteFileNo  = fmt.Sprintf("[ %s ]", blurredStyle.Render("No"))
)

type Page interface {
	tea.Model
}

type ShowResponse struct {
	data       *model.Data
	focusIndex int
	logField   keeperlog.LogField
}

func NewShowResponse(model *model.Data) Page {
	showResponse := &ShowResponse{}
	showResponse.data = model

	// Логируем запросы на сервер в файл keeper.log
	logField := keeperlog.LogField{NextStepByName: showResponse.data.NextStep.NextStepByName}
	logField.Action = "NewShowResponse"
	logField.RequestByName = showResponse.data.NextStep.RequestByName

	// Установка UserID в логгер
	showResponse.data.Log.SetUserID(showResponse.data.User.IDUser)
	showResponse.logField = logField
	return showResponse
}

func (m ShowResponse) Init() tea.Cmd {
	return nil
}

// Файлы
func (m ShowResponse) NextStepToBackListFile() {
	m.logField.Method = "NextStepToBackListFile"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set to Next Step: back to list file")

	m.data.NextStep.NextStepByName = "showfileslist"
}

func (m ShowResponse) NextStepToDeleteFile() {
	m.logField.Method = "NextStepToDeleteFile"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set to Next Step: requesthttp - remove file")

	m.data.NextStep.NextStepByName = "requesthttp"
	m.data.NextStep.RequestByName = "deletefile"
}

func (m ShowResponse) NextStepToSend() {
	m.logField.Method = "NextStepToSend"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set to Next Step: requesthttp - getfiles")

	m.data.NextStep.NextStepByName = "requesthttp"
	m.data.NextStep.RequestByName = "getfiles"
}

// Карты
func (m ShowResponse) NextStepToDeleteCard() {
	m.logField.Method = "NextStepToDeleteCard"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set to Next Step: requesthttp - remove file")
	m.data.NextStep.NextStepByName = "requesthttp"
	m.data.NextStep.RequestByName = "deletecard"
}

func (m ShowResponse) NextStepToBackListCard() {
	m.logField.Method = "NextStepToBackListCard"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set to Next Step: back to list file")

	m.data.NextStep.NextStepByName = "showcardslist"
}

// Пароли
func (m ShowResponse) NextStepToDeletePassword() {
	m.logField.Method = "NextStepToDeletePassword"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set to Next Step: requesthttp - remove file")
	m.data.NextStep.NextStepByName = "requesthttp"
	m.data.NextStep.RequestByName = "deletepassword"
}

func (m ShowResponse) NextStepToBackListPassword() {
	m.logField.Method = "NextStepToBackListPassword"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set to Next Step: back to list file")

	m.data.NextStep.NextStepByName = "showpasswordlist"
}

// Update loop for the first view where you're choosing a task.
func (m ShowResponse) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.SetNextStepByNameExit()
			return m, tea.Quit
			// Set focus to next input
		case "tab", "enter":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" {
				switch m.data.NextStep.RequestByName {
				case "dialogfilescopyorrew":
					m.FileSelectRewriteOrNot()
				case "dialogfilesdeleteornot":
					m.FileDeleteOrNot()
				case "dialogcarddeleteornot":
					m.CardDeleteOrNot()
				case "dialogpassworddeleteornot":
					m.PasswordDeleteOrNot()
				}
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "tab" {
				m.focusIndex++
			}

			if m.focusIndex > maxFocusIndex {
				m.focusIndex = 0
			}

			return m, nil
		}
	}
	return m, nil
}

// Карты
func (m ShowResponse) PasswordDeleteOrNot() {
	if m.focusIndex == 1 {
		m.NextStepToDeletePassword()
		return
	}

	m.NextStepToBackListPassword()
}

// Карты
func (m ShowResponse) CardDeleteOrNot() {
	if m.focusIndex == 1 {
		m.NextStepToDeleteCard()
		return
	}

	m.NextStepToBackListCard()
}

// Файлы
func (m ShowResponse) FileDeleteOrNot() {
	if m.focusIndex == 1 {
		m.NextStepToDeleteFile()
		return
	}

	m.NextStepToBackListFile()
}

func (m ShowResponse) FileSelectRewriteOrNot() {
	request := m.data.RequestHTTP["fileslist"]

	if m.focusIndex == 2 {
		m.NextStepToBackListFile()
		return
	}

	if m.focusIndex == 1 {

		request.OutputFileCopyOrRewrite = true

	} else {
		currentTime := time.Now()
		formattedTime := currentTime.Format("2006-01-02_15_04_05_")
		// Получаем путь к директории
		dirPath := filepath.Dir(request.OutputFileName)

		// Получаем имя файла
		fileName := filepath.Base(request.OutputFileName)
		request.OutputFileName = filepath.Join(dirPath, formattedTime+fileName)
	}

	m.data.RequestHTTP["fileslist"] = request
	m.NextStepToSend()
}

func (m ShowResponse) SetNextStepByNameExit() {
	m.logField.Method = "SetNextStepByNameExit"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set to Next Step: requesthttp - getfiles")

	if m.data.User.Auth {
		m.data.NextStep.NextStepByName = "submenu"
	} else {
		m.data.NextStep.NextStepByName = "topmenu"
	}

}

// The first view, where you're choosing a task
func (m ShowResponse) View() string {
	switch m.data.NextStep.RequestByName {
	case "register":
		return m.ResponseRegister()
	case "login":
		return m.ResponseLogin()
	case "exchangeget":
		return m.ResponseExchangeGet()
	case "exchangeset":
		return m.ResponseExchangeSet()
	case "passwordslist":
		return m.ResponsePasswordsList()
	case "newpassword":
		return m.ResponseNewPassword()
	case "dialogfilescopyorrew":
		maxFocusIndex = 2
		return m.DialogSelect()
	case "dialogfilesdeleteornot":
		maxFocusIndex = 1
		return m.DialogFilesDeleteOrNot()
	case "deletefile":
		return m.ResponseDeleteFile()
	case "newcreditcard":
		return m.ResponseNewCard()
	case "dialogcarddeleteornot":
		maxFocusIndex = 1
		return m.DialogCardDeleteOrNot()
	case "deletecard":
		return m.ResponseDeleteCard()
	case "dialogpassworddeleteornot":
		maxFocusIndex = 1
		return m.DialogPasswordDeleteOrNot()
	case "deletepassword":
		return m.ResponseDeletePassword()

	}
	return ""
}
