package requesthttp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type FilesList struct {
	Description string `json:"description"`
	Filename    string `json:"filename"`
	Size        string `json:"size"`
}

func (h RequestHTTP) FilesList() tea.Msg {

	type List struct {
		Type string `json:"type"`
	}

	list := &List{
		Type: "get_list_files",
	}

	request := h.data.RequestHTTP["fileslist"]

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

	h.parseResponseFilesList(resp)

	return statusMsg(resp.StatusCode())

}

func (h RequestHTTP) parseResponseFilesList(resp *resty.Response) {

	fileslist := h.data.RequestHTTP["fileslist"]

	if len(resp.Body()) > 0 {

		fileslist.Response.Body = string(resp.Body())
		err := json.Unmarshal(resp.Body(), fileslist.ResponseLists)
		if err != nil {
			fmt.Println(" json.Unmarshal(resp.Body(), passwordslist.ResponseLists)===>", err)
		}

		modelOptionFiles := []model.OptionFiles{}

		for _, item := range *fileslist.ResponseLists {
			baseDecode, err := base64.StdEncoding.DecodeString(item.Data)
			if err != nil {
				fmt.Println("Error decoding body of response", err)
			}

			decryptBody, err := h.asimencrypt.DecryptOAEP(baseDecode)
			//fmt.Println(string(decryptBody))
			if err != nil {
				fileslist.Response.Error = fmt.Sprintf("Error decoding %s", err.Error())
			} else {
				fileslist.Response.Body = string(decryptBody)

				filesListUnmarshal := model.OptionFiles{}
				filesListUnmarshal.ID = item.ID
				err := json.Unmarshal(decryptBody, &filesListUnmarshal)
				if err != nil {
					fmt.Println("Error decoding password list", err)
				}
				modelOptionFiles = append(modelOptionFiles, filesListUnmarshal)
			}

		}

		fileslist.FilesList = &modelOptionFiles

	} else {
		fileslist.Response.Body = "No Body"
	}

	fileslist.Response.StatusCode = resp.StatusCode()

	h.data.RequestHTTP["fileslist"] = fileslist

	//time.Sleep(10 * time.Second)

}
