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
	asimencrypt *asimencrypt.AsimEncryptStruct
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

func NewRequestHTTP(client *resty.Client, model *model.Data, asimencrypt *asimencrypt.AsimEncryptStruct) *RequestHTTP {
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

func (h RequestHTTP) NextStepToSendAuthDataToServer() {
	h.logField.Method = "NextStepToSendAuthDataToServer"
	h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info("Set to Next Step - topmenu")

	h.data.NextStep.NextStepByName = "topmenu"
}

func (h RequestHTTP) NextStepToShowShowDialog() {
	h.logField.Method = "NextStepToShowShowDialog"
	h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info("Set to Next Step - showresponse")

	h.data.NextStep.NextStepByName = "showresponse"

}

func (h RequestHTTP) NextStepToShowShowREsponse() {

	switch h.data.NextStep.RequestByName {
	case "getfiles":
		h.data.NextStep.NextStepByName = "files"
	case "getfileslist":
		h.data.NextStep.NextStepByName = "showfileslist"
	case "sendfiles":
		h.data.NextStep.NextStepByName = "files"
	case "passwordslist":
		h.data.NextStep.NextStepByName = "showpasswordslist"
	case "cardslist":
		h.data.NextStep.NextStepByName = "showcardslist"
	default:
		h.data.NextStep.NextStepByName = "showresponse"
	}
}

func (h RequestHTTP) Init() tea.Cmd {
	// Создание канала для отмены операции
	h.data.Cancel = make(chan struct{})

	h.logField.Method = "init"
	// Логирование информации о запросе
	h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info("Init request")

	// Обработка различных запросов
	switch h.data.NextStep.RequestByName {
	case "register":
		return h.Register
	case "login":
		return h.Login
	case "exchangeget":
		return h.ExchangeGet
	case "exchangeset":
		return h.ExchangeSet
	case "newpassword":
		return h.PasswordNew
	case "passwordslist":
		return h.PasswordsList
	case "sendfiles":
		return h.RunSendFiles
	case "getfiles":
		return h.RunGetFiles
	case "getfileslist":
		return h.FilesList
	case "deletefile":
		return h.DeleteFile
	case "newcreditcard":
		return h.Card
	case "cardslist":
		return h.CardsList
	case "deletecard":
		return h.CardDelete
	case "deletepassword":
		return h.PasswordDelete

	}
	return nil
}

func (h RequestHTTP) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	h.logField.Method = "update"
	switch msg := msg.(type) {
	case tea.KeyMsg:
		h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info("Key pressed")

		switch msg.String() {
		case "q", "ctrl+c":
			h.data.Cancel <- struct{}{}
			h.NextStepToSendAuthDataToServer()
			return h, tea.Quit
		default:
			return h, nil
		}

	case statusMsg:
		h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info(fmt.Sprintf("statusMsg - %d", int(msg)))

		h.status = int(msg)
		h.NextStepToShowShowREsponse()
		return h, tea.Quit

	case errMsg:
		h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Error(fmt.Sprintf("statusMsg - %s", msg.Error()))

		h.err = msg
		if msg.Error() == "ExistFile" {
			// Обработка ошибки ExistFile
		} else {
			h.NextStepToShowShowREsponse()
		}
		return h, tea.Quit

	default:
		h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info("Default")
		return h, nil
	}
}

func (h RequestHTTP) View() string {
	s := "Выполняем запрос, пожалуйста, ожидайте ..." + h.data.NextStep.RequestByName
	if h.err != nil {
		s += fmt.Sprintf("something went wrong: %s", h.err)
	} else if h.status != 0 {
		s += fmt.Sprintf("%d %s", h.status, http.StatusText(h.status))
	}
	return s + "\n"
}

// Обновляем токен прозрачно для клиента
func (h RequestHTTP) RefreshToken() error {
	h.logField.Method = "RefreshToken"
	//Получаем токен для обновления основного токена
	tokenRefresh := h.data.User.TokenRefresh
	// Send a request to refresh the token
	resp, err := h.client.R().
		SetHeader("Token-Refresh", tokenRefresh).
		Post(fmt.Sprintf("%s/files/download", h.data.Config.Server.AddressFileServer))

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to refresh token: %s", resp.String())
	}

	//Устанавливаем токены в модель
	h.data.User.Auth = true
	h.data.User.Token = resp.Header().Get("Token")
	h.data.User.TokenRefresh = resp.Header().Get("Token-Refresh")

	h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info("Successfully RefreshToken()")

	return nil
}
