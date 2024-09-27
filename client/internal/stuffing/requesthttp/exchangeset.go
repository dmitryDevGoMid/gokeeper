package requesthttp

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

func (h RequestHTTP) ExchangeSet() tea.Msg {

	fmt.Println("ExchangeSet!")

	//type Exchange struct {
	//	Key string `json:"key"`
	//Type string `json:"type"`
	//}

	//key := h.asimencrypt.GetBytePublic()

	//fmt.Println("KEY====>", key)

	//exchange := &Exchange{
	//	Key: string(key),
	//Type: "set_public_key",
	//}

	request := h.data.RequestHTTP["exchangeset"]

	//jsonRequest, err := json.Marshal(exchange)

	//fmt.Println("jsonRequest=====>", jsonRequest)

	//if err != nil {
	//	return errMsg{err}
	//}

	//fmt.Println("SIZE MESSAGE====>", len(jsonRequest))
	//request.Request = string("jsonRequest")

	//fmt.Println("jsonRequest=====>", request.Request)

	//encriptSend, err := h.asimencrypt.EncryptByServerKey(string(jsonRequest))
	//encriptSend := jsonRequest // h.asimencrypt.EncryptByServerKey(string(jsonRequest))
	//rsa, _ := fmt.Println(h.asimencrypt.ReadServerPublicKey())
	//fmt.Println("RSA======>", rsa)

	encriptSend, err := h.asimencrypt.EncryptByServerKeyParts(request.Request)
	//fmt.Println(err)

	//time.Sleep(2 * time.Second)
	if err != nil {
		return errMsg{err}
	}

	resp, err := h.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(encriptSend)).
		//SetResult(h.data.Keys).
		Post(request.URL)

	if err != nil {
		return errMsg{err}
	}

	//err = h.asimencrypt.SetPublicServerKey(h.data.Keys.Key)

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
