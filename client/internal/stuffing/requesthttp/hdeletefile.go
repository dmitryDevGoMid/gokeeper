package requesthttp

import (
	"encoding/json"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

func (h RequestHTTP) DeleteFile() tea.Msg {

	/*passwordReq := &PasswordRequest{
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
	}*/

	file := h.data.RequestHTTP["fileslist"]

	chocedKey, err := strconv.Atoi(file.Choced)
	if err != nil {
		return err
	}
	filesList := *file.FilesList
	selectedFile := filesList[chocedKey]
	selectedFile.ID_File = selectedFile.ID

	request := h.data.RequestHTTP["deletefile"]

	jsonRequest, err := json.Marshal(selectedFile)

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
		SetHeader("Token", h.data.User.Token).
		//SetHeader("File-ID", selectedFile.ID).
		SetBody(string(encriptSend)).
		Post(request.URL)

	if err != nil {
		return errMsg{err}
	}

	h.parseResponseDelete(resp)

	return statusMsg(200)
}

func (h RequestHTTP) parseResponseDelete(resp *resty.Response) {

	request := h.data.RequestHTTP["deletefile"]

	request.Response.StatusCode = resp.StatusCode()

	h.data.RequestHTTP["deletefile"] = request
}
