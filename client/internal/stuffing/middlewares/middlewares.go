package middlewares

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/keeperlog"

	"github.com/go-resty/resty/v2"
)

type CleintInterface interface {
	SetTokenToRequest(data *model.Data, asimencrypt *asimencrypt.AsimEncryptStruct)
	CheckStatusCode401(data *model.Data)
}

type Client struct {
	client   *resty.Client
	data     *model.Data
	logField keeperlog.LogField
}

func NewClientMiddleware(client *resty.Client, data *model.Data) CleintInterface {
	clientType := &Client{client: client, data: data}

	// Логируем запросы на сервер в файл keeper.log
	logField := keeperlog.LogField{NextStepByName: clientType.data.NextStep.NextStepByName}
	logField.Action = "NewClientMiddleware"
	logField.RequestByName = clientType.data.NextStep.RequestByName

	// Установка UserID в логгер
	clientType.data.Log.SetUserID(clientType.data.User.IDUser)
	clientType.logField = logField
	return clientType
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func (cl *Client) CheckStatusCode401(data *model.Data) {

	cl.client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		cl.logField.Method = "CheckStatusCode401 - OnAfterResponse"
		var url string
		splitUrl := strings.Split(resp.Request.URL, "/api/user/")
		if 1 < len(splitUrl) {
			url = splitUrl[1]
		}

		if url != "login" {
			cl.data.Log.WithFields(keeperlog.ToMap(cl.logField)).Info("Url " + url)

			if resp.StatusCode() == 401 {
				data.User.Auth = false
				if err := cl.RefreshToken(data.User.TokenRefresh); err != nil {
					cl.data.Log.WithFields(keeperlog.ToMap(cl.logField)).Error("Error refreshing token: " + err.Error())
					return err
				}
				c.SetRetryCount(1)
				cl.data.Log.WithFields(keeperlog.ToMap(cl.logField)).Info("Refreshing toke successfully")

			}
		}
		return nil
	})

}

// Обновляем токен прозрачно для клиента
func (cl *Client) RefreshToken(tokenRefresh string) error {
	cl.logField.Method = "RefreshToken"

	// Send a request to refresh the token
	resp, err := cl.client.R().
		SetHeader("Token-Refresh", tokenRefresh).
		Post("http://localhost:8000/api/user/refresh/token")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		cl.data.Log.WithFields(keeperlog.ToMap(cl.logField)).Error(fmt.Sprintf("Status code: %d", resp.StatusCode()))
		return fmt.Errorf("failed to refresh token: %s", resp.String())
	}

	//Устанавливаем токены в модель
	cl.data.User.Auth = true
	cl.data.User.Token = resp.Header().Get("Token")
	cl.data.User.TokenRefresh = resp.Header().Get("Token-Refresh")

	return nil
}

// Set Token to request
func (cl *Client) SetTokenToRequest(data *model.Data, asimencrypt *asimencrypt.AsimEncryptStruct) {
	cl.client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		cl.logField.Method = "SetTokenToRequest - OnBeforeRequest - Public-Key"
		var url string
		splitUrl := strings.Split(req.URL, "/api/user/")
		if 1 < len(splitUrl) {
			url = splitUrl[1]
		}
		cl.data.Log.WithFields(keeperlog.ToMap(cl.logField)).Info(fmt.Sprintf("Status code: %s", url))
		if url != "login" && url != "exchange/set" && url != "exchange/get" && data.User.Token != "" {
			req.Header.Set("Token", data.User.Token)
			req.Header.Set("Public-Key", base64.StdEncoding.EncodeToString(asimencrypt.GetBytePublic()))
		}
		return nil
	})
}
