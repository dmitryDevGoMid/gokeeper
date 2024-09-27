package middlewares

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/keeperlog"

	"github.com/go-resty/resty/v2"
)

type CleintInterface interface {
	DecryptResponse(asimencrypt asimencrypt.AsimEncrypt)
	SetTokenToRequest(data *model.Data, asimencrypt asimencrypt.AsimEncrypt)
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
				//c.SetRetryWaitTime(1 * time.Microsecond)
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

func (cl *Client) DecryptResponse(asimencrypt asimencrypt.AsimEncrypt) {
	// Registering Response Middleware
	cl.client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		splitUrl := strings.Split(resp.Request.URL, "/api/user/")
		url := splitUrl[1]
		fmt.Println("DecryptResponse URL_SEND_ADRESS:", splitUrl[1])
		//time.Sleep(3 * time.Second)
		if url != "login" && url != "exchange/set" && url != "exchange/get" {
			//fmt.Println("Middleware DecryptResponse:")
			//fmt.Println(string(resp.Body()))
			//time.Sleep(3 * time.Second,
			body, err := io.ReadAll(resp.RawBody())
			if err != nil {
				fmt.Println("error io.ReadAll middlewwre", err)
			}

			fmt.Println("ОН САМЫЙ БОДИ:", string(body))

			os.Exit(1)

			fmt.Println("BODY Response:")
			fmt.Println("BODY Response:")
			fmt.Println("BODY Response:")
			fmt.Println("BODY Response:")
			fmt.Println("BODY Response:")
			fmt.Println("BODY Response:")
			fmt.Println("BODY Response:", string(resp.Body()))
			fmt.Println("BODY Response:")
			fmt.Println("BODY Response:")
			fmt.Println("BODY Response:")
			fmt.Println("BODY Response:")
			fmt.Println("BODY Response:")
			time.Sleep(3 * time.Second)
			base64Body, err := base64.StdEncoding.DecodeString(string(resp.Body()))
			if err != nil {
				fmt.Println("Error decoding base64:", err)
			}
			decodeData, err := asimencrypt.DecryptOAEP(base64Body)
			if err != nil {
				fmt.Println("error decode middlewwre", err)
			}
			//fmt.Println(" DecryptResponse RESULT:", string(decodeData))
			//time.Sleep(3 * time.Second)
			resp.SetBody(decodeData)
		}
		//body, _ := io.ReadAll(c.Request.Body)
		/*body, err := io.ReadAll(resp.RawResponse.Body)
		if err != nil {
			fmt.Println("error io.ReadAll middlewwre", err)
		}

		decodeData, err := asimencrypt.DecryptOAEP(body)
		if err != nil {
			fmt.Println("error decode middlewwre", err)
		}

		fmt.Println("ENCRYPTE=====>>>>")
		fmt.Println("ENCRYPTE=====>>>>", string(decodeData))

		//decompressBody, _ := asimencrypt.DecryptOAEP(body)

		fmt.Println("decompressBody=====>>>>")
		fmt.Println("decompressBody=====>>>>", string(decodeData))

		//c.Request.Body = io.NopCloser(bytes.NewReader(decodeData))
		resp. = io.NopCloser(bytes.NewReader(decodeData))*/
		// Now you have access to Client and current Response object
		// manipulate it as per your need
		//fmt.Println("RESPONSE:", resp.Header().Get("Content-Type"))
		//fmt.Println(string(resp.Body()))

		return nil // if its success otherwise return error
	})
}

// Set Token to request
func (cl *Client) SetTokenToRequest(data *model.Data, asimencrypt asimencrypt.AsimEncrypt) {
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
