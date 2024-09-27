package requesthttp

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

func (h RequestHTTP) Card() tea.Msg {
	//Проверяем флаг редактирования иначе сбрасываем его
	if !h.data.OptionCards.Edite {
		h.data.OptionCards.ID = ""
	} else {
		h.data.OptionCards.Edite = false
	}

	cardReq := &CardRequest{
		ID:          h.data.OptionCards.ID,
		Description: h.data.OptionCards.Description,
		Number:      h.data.OptionCards.Number,
		Exp:         h.data.OptionCards.Exp,
		Cvc:         h.data.OptionCards.Cvc,
	}

	request := h.data.RequestHTTP["newcard"]

	jsonRequest, err := json.Marshal(cardReq)

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

	h.parseResponseNewCard(resp)

	return statusMsg(resp.StatusCode())
}

func (h RequestHTTP) parseResponseNewCard(resp *resty.Response) {

	request := h.data.RequestHTTP["newcard"]

	if len(resp.Body()) > 0 {
		request.Response.Body = string(resp.Body())
	} else {
		request.Response.Body = "No Body"
	}

	request.Response.StatusCode = resp.StatusCode()

	h.data.RequestHTTP["newcard"] = request

}
