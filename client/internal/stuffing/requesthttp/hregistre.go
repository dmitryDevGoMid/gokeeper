package requesthttp

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type RegisterRequest struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Key      string `json:"key"`
}

func (h RequestHTTP) Register() tea.Msg {

	// Публичный ключ клиента
	key := h.asimencrypt.GetBytePublic()

	// Регистрация
	registerReq := &RegisterRequest{
		ID:       h.data.User.ID,
		Username: h.data.User.Username,
		Password: h.data.User.Password,
		Key:      string(key),
	}

	request := h.data.RequestHTTP["register"]

	jsonRequest, err := json.Marshal(registerReq)

	if err != nil {
		return errMsg{err}
	}

	request.Request = string(jsonRequest)

	//Шифруем данные ключом сервера!
	encriptSend, err := h.asimencrypt.EncryptByServerKeyParts(string(jsonRequest))

	if err != nil {
		return errMsg{err}
	}

	resp, err := h.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Encrypt", "yes").
		SetBody(string(encriptSend)).
		Post(request.URL)

	if err != nil {
		return errMsg{err}
	}

	h.parseResponseRegister(resp)

	return statusMsg(resp.StatusCode())
}

func (h RequestHTTP) parseResponseRegister(resp *resty.Response) {

	request := h.data.RequestHTTP["register"]

	if len(resp.Body()) > 0 {
		request.Response.Body = string(resp.Body())
	} else {
		request.Response.Body = "No Body"
	}

	request.Response.StatusCode = resp.StatusCode()

	h.data.RequestHTTP["register"] = request

}
