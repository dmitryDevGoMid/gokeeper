package requesthttp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type PasswordList struct {
	Description string `json:"description"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

func (h RequestHTTP) PasswordsList() tea.Msg {

	type List struct {
		Type string `json:"type"`
	}

	list := &List{
		Type: "get_list_passwords",
	}

	request := h.data.RequestHTTP["passwordslist"]

	jsonRequest, err := json.Marshal(list)

	if err != nil {
		return errMsg{err}
	}

	request.Request = string(jsonRequest)

	encriptSend, err := h.asimencrypt.EncryptByServerKeyParts(string(jsonRequest))

	if err != nil {
		return errMsg{err}
	}

	resp, err := h.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Token", h.data.User.Token).
		SetBody(string(encriptSend)).
		Post(request.URL)

	if err != nil {
		return errMsg{err}
	}

	h.parseResponsePasswordsList(resp)

	return statusMsg(resp.StatusCode())

}

func (h RequestHTTP) parseResponsePasswordsList(resp *resty.Response) {

	passwordslist := h.data.RequestHTTP["passwordslist"]

	if len(resp.Body()) > 0 {

		passwordslist.Response.Body = string(resp.Body())
		err := json.Unmarshal(resp.Body(), passwordslist.ResponseLists)

		if err != nil {
			fmt.Println(" json.Unmarshal(resp.Body(), passwordslist.ResponseLists)===>", err)
		}

		modelOptionPasswords := []model.OptionPasswords{}

		for _, item := range *passwordslist.ResponseLists {
			baseDecode, err := base64.StdEncoding.DecodeString(item.Data)
			if err != nil {
				fmt.Println("Error decoding body of response", err)
			}
			decryptBody, err := h.asimencrypt.DecryptOAEP(baseDecode)
			if err != nil {
				passwordslist.Response.Error = fmt.Sprintf("Error decoding %s", err.Error())
			} else {
				passwordslist.Response.Body = string(decryptBody)

				passwordsListUnmarshal := model.OptionPasswords{}
				passwordsListUnmarshal.ID = item.ID
				err := json.Unmarshal(decryptBody, &passwordsListUnmarshal)
				if err != nil {
					fmt.Println("Error decoding password list", err)
				}
				modelOptionPasswords = append(modelOptionPasswords, passwordsListUnmarshal)
			}

		}

		passwordslist.PasswordsList = &modelOptionPasswords

	} else {
		passwordslist.Response.Body = "No Body"
	}

	passwordslist.Response.StatusCode = resp.StatusCode()

	h.data.RequestHTTP["passwordslist"] = passwordslist

}
