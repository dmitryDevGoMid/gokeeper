package requesthttp

import (
	"encoding/json"
	"fmt"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/keeperlog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

func (h RequestHTTP) ExchangeGet() tea.Msg {

	h.logField.Method = "ExchangeGet"

	type Exchange struct {
		Type string `json:"type"`
	}

	exchange := &Exchange{
		Type: "get_public_key",
	}

	request := h.data.RequestHTTP["exchangeget"]

	jsonRequest, err := json.Marshal(exchange)

	if err != nil {
		return errMsg{err}
	}

	request.Request = string(jsonRequest)

	//encriptSend, err := h.asimencrypt.EncryptByServerKey(string(jsonRequest))
	encriptSend := jsonRequest

	if err != nil {
		return errMsg{err}
	}

	resp, err := h.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(encriptSend)).
		SetResult(h.data.Keys).
		Post(request.URL)

	if err != nil {
		h.parseResponseExchangeSetError(err)
		return errMsg{err}
	}

	err = h.asimencrypt.SetPublicServerKey(h.data.Keys.Key)

	//fmt.Println(h.data.Keys.Key)
	//time.Sleep(2 * time.Second)

	if err != nil {
		return errMsg{err}
	}

	h.parseResponseExchangeGet(resp)

	return statusMsg(resp.StatusCode())

}

func (h RequestHTTP) parseResponseExchangeSetError(err error) {
	h.logField.Method = "parseResponseExchangeSetError"
	h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info("Request")

	exchange := h.data.RequestHTTP["exchangeget"]
	exchange.Response.Error = fmt.Sprintf("Error: %v", err)
	h.data.RequestHTTP["exchangeget"] = exchange
}

func (h RequestHTTP) parseResponseExchangeGet(resp *resty.Response) {

	h.logField.Method = "parseResponseExchangeGet"
	h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info(fmt.Sprintf("Request: %d", resp.StatusCode()))

	exchange := h.data.RequestHTTP["exchangeget"]

	if len(resp.Body()) > 0 {
		exchange.Response.Body = string(resp.Body())
	} else {
		exchange.Response.Body = "No Body"
	}
	exchange.Response.Body = "No Body"

	exchange.Response.StatusCode = resp.StatusCode()

	h.data.RequestHTTP["exchangeget"] = exchange

}
