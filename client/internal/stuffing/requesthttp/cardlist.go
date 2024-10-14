package requesthttp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type CardList struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Number      string `json:"login"`
	Exp         string `json:"exp"`
	Cvc         int    `json:"cvc"`
}

func (h RequestHTTP) CardsList() tea.Msg {

	type List struct {
		Type string `json:"type"`
	}

	list := &List{
		Type: "get_list_cards",
	}

	request := h.data.RequestHTTP["cardslist"]

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

	h.parseResponseCardsList(resp)

	return statusMsg(resp.StatusCode())

}

func (h RequestHTTP) parseResponseCardsList(resp *resty.Response) {

	cardslist := h.data.RequestHTTP["cardslist"]

	if len(resp.Body()) > 0 {

		cardslist.Response.Body = string(resp.Body())
		err := json.Unmarshal(resp.Body(), cardslist.ResponseLists)

		if err != nil {
			fmt.Println(" json.Unmarshal(resp.Body(), cardslist.ResponseLists)===>", err)
		}

		modelOptionCards := []model.OptionCards{}

		for _, item := range *cardslist.ResponseLists {
			baseDecode, err := base64.StdEncoding.DecodeString(item.Data)
			if err != nil {
				fmt.Println("Error decoding body of response", err)
			}
			decryptBody, err := h.asimencrypt.DecryptOAEP(baseDecode)

			if err != nil {
				cardslist.Response.Error = fmt.Sprintf("Error decoding %s", err.Error())
			} else {
				cardslist.Response.Body = string(decryptBody)

				cardsListUnmarshal := model.OptionCards{}
				err := json.Unmarshal(decryptBody, &cardsListUnmarshal)
				if err != nil {
					fmt.Println("Error decoding cards list", err)
				}
				cardsListUnmarshal.ID = item.ID
				modelOptionCards = append(modelOptionCards, cardsListUnmarshal)
			}

		}

		cardslist.CardsList = &modelOptionCards

	} else {
		cardslist.Response.Body = "No Body"
	}

	cardslist.Response.StatusCode = resp.StatusCode()

	h.data.RequestHTTP["cardslist"] = cardslist

}
