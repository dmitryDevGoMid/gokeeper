package requesthttp

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

func (h RequestHTTP) ExchangeSet() tea.Msg {

	request := h.data.RequestHTTP["exchangeset"]

	encriptSend, err := h.asimencrypt.EncryptByServerKeyParts(request.Request)
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

	if err != nil {
		return errMsg{err}
	}

	h.parseResponseExchangeSet(resp)

	return statusMsg(resp.StatusCode())

}

func (h RequestHTTP) parseResponseExchangeSet(resp *resty.Response) {

	exchange := h.data.RequestHTTP["exchangeset"]

	if len(resp.Body()) > 0 {
		exchange.Response.Body = string(resp.Body())
	} else {
		exchange.Response.Body = "No Body"
	}

	exchange.Response.Body = "No Body"

	exchange.Response.StatusCode = resp.StatusCode()

	h.data.RequestHTTP["exchangeset"] = exchange

}
