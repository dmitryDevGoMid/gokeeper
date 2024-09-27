package requesthttp

// A simple program that makes a GET request and prints the response status.

import (
	"fmt"
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/keeperlog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type RequestHTTP struct {
	data        *model.Data
	asimencrypt asimencrypt.AsimEncrypt
	client      *resty.Client
	status      int
	err         error
	logField    keeperlog.LogField
}

type CardRequest struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Number      string `json:"login"`
	Exp         string `json:"exp"`
	Cvc         string `json:"cvc"`
}

type Page interface {
	tea.Model
	NextStepToSendAuthDataToServer()
	NextStepToShowShowREsponse()
	EnsureDirectoryExists(dirName string) (string, error)
}

type statusMsg int

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func NewRequestHTTP(client *resty.Client, model *model.Data, asimencrypt asimencrypt.AsimEncrypt) *RequestHTTP {
	request := &RequestHTTP{data: model, client: client, asimencrypt: asimencrypt}
	request.data = model

	// Логируем запросы на сервер в файл keeper.log
	logField := keeperlog.LogField{NextStepByName: request.data.NextStep.NextStepByName}
	logField.Action = "request"
	logField.RequestByName = request.data.NextStep.RequestByName

	// Установка UserID в логгер
	request.data.Log.SetUserID(request.data.User.IDUser)
	request.logField = logField
	return request
}

func (m RequestHTTP) NextStepToSendAuthDataToServer() {
	m.logField.Method = "NextStepToSendAuthDataToServer"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set to Next Step - topmenu")

	m.data.NextStep.NextStepByName = "topmenu"
}

func (m RequestHTTP) NextStepToShowShowDialog() {
	m.logField.Method = "NextStepToShowShowDialog"
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Set to Next Step - showresponse")

	m.data.NextStep.NextStepByName = "showresponse"

}

func (m RequestHTTP) NextStepToShowShowREsponse() {

	switch m.data.NextStep.RequestByName {
	case "getfiles":
		m.data.NextStep.NextStepByName = "files"
	case "getfileslist":
		m.data.NextStep.NextStepByName = "showfileslist"
	case "sendfiles":
		m.data.NextStep.NextStepByName = "files"
	case "passwordslist":
		m.data.NextStep.NextStepByName = "showpasswordslist"
	case "cardslist":
		m.data.NextStep.NextStepByName = "showcardslist"
	default:
		m.data.NextStep.NextStepByName = "showresponse"
	}
}

func (m RequestHTTP) Init() tea.Cmd {
	// Создание канала для отмены операции
	m.data.Cancel = make(chan struct{})

	m.logField.Method = "init"
	// Логирование информации о запросе
	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Init request")

	// Обработка различных запросов
	switch m.data.NextStep.RequestByName {
	case "register":
		return m.Register
	case "login":
		return m.Login
	case "exchangeget":
		return m.ExchangeGet
	case "exchangeset":
		return m.ExchangeSet
	case "newpassword":
		return m.PasswordNew
	case "passwordslist":
		return m.PasswordsList
	case "sendfiles":
		return m.RunSendFiles
	case "getfiles":
		return m.RunGetFiles
	case "getfileslist":
		return m.FilesList
	case "deletefile":
		return m.DeleteFile
	case "newcreditcard":
		return m.Card
	case "cardslist":
		return m.CardsList
	case "deletecard":
		return m.CardDelete
	case "deletepassword":
		return m.PasswordDelete

	}
	return nil
}

func (m RequestHTTP) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logField.Method = "update"
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Key pressed")

		switch msg.String() {
		case "q":
			m.NextStepToSendAuthDataToServer()
			return m, tea.Quit
		case "ctrl+c":
			m.data.Cancel <- struct{}{}
			m.NextStepToShowShowREsponse()
			return m, tea.Quit
		default:
			return m, nil
		}

	case statusMsg:
		m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info(fmt.Sprintf("statusMsg - %d", int(msg)))

		m.status = int(msg)
		m.NextStepToShowShowREsponse()
		return m, tea.Quit

	case errMsg:
		m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Error(fmt.Sprintf("statusMsg - %s", msg.Error()))

		m.err = msg
		if msg.Error() == "ExistFile" {
			// Обработка ошибки ExistFile
		} else {
			m.NextStepToShowShowREsponse()
		}
		return m, tea.Quit

	default:
		//logField := keeperlog.LogField{Action: "Update", RequestByName: "default"}
		m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Default")
		return m, nil
	}
}

func (m RequestHTTP) View() string {
	s := "Выполняем запрос, пожалуйста, ожидайте ..." + m.data.NextStep.RequestByName
	if m.err != nil {
		s += fmt.Sprintf("something went wrong: %s", m.err)
	} else if m.status != 0 {
		s += fmt.Sprintf("%d %s", m.status, http.StatusText(m.status))
	}
	return s + "\n"
}

// Обновляем токен прозрачно для клиента
func (m RequestHTTP) RefreshToken() error {
	m.logField.Method = "RefreshToken"
	//Получаем токен для обновления основного токена
	tokenRefresh := m.data.User.TokenRefresh
	// Send a request to refresh the token
	resp, err := m.client.R().
		//SetQueryParam("refresh_token", "your_refresh_token").
		SetHeader("Token-Refresh", tokenRefresh).
		Post("http://localhost:8000/api/user/refresh/token")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to refresh token: %s", resp.String())
	}

	//Устанавливаем токены в модель
	m.data.User.Auth = true
	m.data.User.Token = resp.Header().Get("Token")
	m.data.User.TokenRefresh = resp.Header().Get("Token-Refresh")

	m.data.Log.WithFields(keeperlog.ToMap(m.logField)).Info("Successfully RefreshToken()")

	return nil
}
