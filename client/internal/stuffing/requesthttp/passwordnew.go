package requesthttp

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type PasswordRequest struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

func (h RequestHTTP) PasswordNew() tea.Msg {

	loginReq := &PasswordRequest{
		ID:          h.data.OptionPasswords.ID,
		Description: h.data.OptionPasswords.Description,
		Username:    h.data.OptionPasswords.Username,
		Password:    h.data.OptionPasswords.Password,
	}

	request := h.data.RequestHTTP["newpassword"]

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

	h.parseResponseNewPassword(resp)

	if err != nil {
		fmt.Println(err)
		return errMsg{err}
	}

	return statusMsg(resp.StatusCode())
}

func (h RequestHTTP) parseResponseNewPassword(resp *resty.Response) {

	request := h.data.RequestHTTP["newpassword"]

	if len(resp.Body()) > 0 {
		request.Response.Body = string(resp.Body())
	} else {
		request.Response.Body = "No Body"
	}

	request.Response.StatusCode = resp.StatusCode()

	h.data.RequestHTTP["newpassword"] = request

}
