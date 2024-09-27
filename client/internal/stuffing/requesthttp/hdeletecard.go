package requesthttp

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

func (h RequestHTTP) CardDelete() tea.Msg {

	cardReq := &CardRequest{
		ID:          h.data.OptionCards.ID,
		Description: h.data.OptionCards.Description,
		Number:      h.data.OptionCards.Number,
		Exp:         h.data.OptionCards.Exp,
		Cvc:         h.data.OptionCards.Cvc,
	}

	request := h.data.RequestHTTP["deletecard"]

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

	h.parseResponseDeleteCard(resp)

	return statusMsg(resp.StatusCode())
}

func (h RequestHTTP) parseResponseDeleteCard(resp *resty.Response) {

	request := h.data.RequestHTTP["deletecard"]

	if len(resp.Body()) > 0 {
		request.Response.Body = string(resp.Body())
	} else {
		request.Response.Body = "No Body"
	}

	request.Response.StatusCode = resp.StatusCode()

	h.data.RequestHTTP["deletecard"] = request

}
