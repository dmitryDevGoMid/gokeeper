package requesthttp

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

func (h RequestHTTP) PasswordDelete() tea.Msg {

	passwordReq := &PasswordRequest{
		ID:          h.data.OptionPasswords.ID,
		Description: h.data.OptionPasswords.Description,
		Username:    h.data.OptionPasswords.Username,
		Password:    h.data.OptionPasswords.Password,
	}

	request := h.data.RequestHTTP["deletepassword"]

	jsonRequest, err := json.Marshal(passwordReq)

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

	h.parseResponseDeletePassword(resp)

	return statusMsg(resp.StatusCode())
}

func (h RequestHTTP) parseResponseDeletePassword(resp *resty.Response) {

	request := h.data.RequestHTTP["deletepassword"]

	if len(resp.Body()) > 0 {
		request.Response.Body = string(resp.Body())
	} else {
		request.Response.Body = "No Body"
	}

	request.Response.StatusCode = resp.StatusCode()

	h.data.RequestHTTP["deletepassword"] = request

}
