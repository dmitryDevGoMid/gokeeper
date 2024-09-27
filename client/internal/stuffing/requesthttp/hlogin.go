package requesthttp

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h RequestHTTP) Login() tea.Msg {

	// Логин
	loginReq := &LoginRequest{
		Username: h.data.User.Username,
		Password: h.data.User.Password,
	}

	request := h.data.RequestHTTP["login"]

	jsonRequest, err := json.Marshal(loginReq)

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
		SetBody(string(encriptSend)).
		Post(request.URL)

	if err != nil {
		return errMsg{err}
	}

	err = h.parseResponseLogin(resp)

	if err != nil {
		return errMsg{err}
	}

	return statusMsg(resp.StatusCode())
}

func (h RequestHTTP) parseResponseLogin(resp *resty.Response) error {

	err := json.Unmarshal(resp.Body(), h.data.User)
	if err != nil {
		return err
	}

	request := h.data.RequestHTTP["login"]

	if len(resp.Body()) > 0 {
		request.Response.Body = string(resp.Body())
	} else {
		request.Response.Body = "No Body"
	}

	request.Response.Body = "No Body"

	request.Response.StatusCode = resp.StatusCode()

	if request.Response.StatusCode == 200 {
		h.data.User.Auth = true
	} else {
		h.data.User.Auth = false
	}

	h.data.RequestHTTP["login"] = request

	return nil

}
