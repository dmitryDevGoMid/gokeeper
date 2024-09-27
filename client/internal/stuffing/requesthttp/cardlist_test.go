package requesthttp

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/config"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/opcardslist"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/keeperlog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/go-playground/assert/v2"

	"github.com/go-resty/resty/v2"
)

var PublicClientKey string = `-----BEGIN RSA PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAsDnFluV73ef/gz8jIUZO
/P3bZ6gyBn8KLFXaPnNh5ajuju4HnXZTJd42JWMWRyYFsA16And3Hf4w6AcHwS9E
TGxClV69LE/OPPbug0bITUWEMQgOEafW29BKBgmAtBEFL5M4OxbomCK7DQTnaBOD
VZqFmS87DZe5l8l31HZooowVIF+fhHjRsRqS3L7t5PLG1WWo606rtj2cuJtdz6df
z8/gSl/IQwMbeu6q0Gii8vna9Yxw1ENOhYWdF4bhTXjposefj4ICV8eLO5iqudV8
Unb3DTQNTJ6aRFfEaPbuKb3xS7wS63i1qFxh9cdyZGpxRzHoaRmVD4abH1kmtTVQ
CnOX95syTgJ4IwHBcgIlt1lpZoVWqQdSJQsCjM36Ax4I0H17vZyKZ7EXiXnpnk/X
HK4+sCQ9BqqfUgSxz7DpMJVGwgUCIYFO7Xgh0iDEjb0M1nBWZy6yc+SOSrOBmQuM
5dgytRhPaUnQcwPqERJClSwavM2C6lsBINGute+HPG2aS4mROedCu/iREGj+zGVB
BWNtim0/4iiqJtyGXyLHfS4vTeThHNjZnb+4FLq5UWdQ7gviL3m/YpcbVTuIYPjO
h6V98oFVtOsJhVoD04nu7X5X+nYu472qUeJ+cBVxPeC+xuPaOqj0CKdVntEh0KeR
X5l46qclblQnvpY7q9cpsyMCAwEAAQ==
-----END RSA PUBLIC KEY-----`

var PrivateClientKey string = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

var PublicServerKey string = `-----BEGIN RSA PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAz+qMWhe2MAXkotvQa3Yd
+euCV4i+mtMw8Ni9vN4SmG416+khsf2NdvLi3U1THvDmKkP1IYAO2zhxe7hTzHm+
M/Z6XF0pe85ax+RS1f4QzgqoaMDEZXhgTO2oBFtgdQrYl6Z/meIAC1MYj/zyF44M
+euSnuemjUiXO4keowfA1dUx9wu1x/703Ma1vd78MwZQzpgl0p5c1pGiaYoIDC/N
GY6VqsovI9DltVxpI2CjcrbZ18X4N+G9rSau9LtfSctYh71jJvxEW7zQS2/zAZMT
A5I31TUHSvYDbtc/HILMGz+5zw/wVgvQ5il3rWYQ3Hm1cBS4z6NVYJhZaY/P9gHd
IX/UYxHmQvMmASJTR0FO13bkJR3pyhNUkmAhWixa0kYB4in6SWJ1HYrR2ZqEh7CC
ImfBbLctpOrcls6GMWKklCRmp2a65DwjM/gtnjV+YpNaWwxFYCC9doWEtohWoSoX
P/lQBeCN56Hey2BjqAdvSlkXu7Wpqg8iY2gdkiKQe4so59MZPljnbtB0XF0UJNF7
MmHnQVmHZKCLAcIOLFQ711qlmP4R+uECLt2vnbJQkNrk8iZM0/K83PfLl9yZaMYI
80bDTCuWrbdthOxzU6TZ0ywYgaOOO+7DVOehIj5CxJfX/Kib2kE9wfXKaWTenDHg
t51NW2aJIuxlHXS3kvxFjh0CAwEAAQ==
-----END RSA PUBLIC KEY-----`

var PrivateServerKey string = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

// Возвращаем приватный rsa ключ
func GetRSAPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	// Decode the file from PEM format
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse RSA")
	}

	// Parse the private key in PKCS#1 format
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil

}

// Возвращаем публичный rsa ключ
func GetRSAPublicKey(data []byte) (*rsa.PublicKey, error) { // Decode the file from PEM format
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse RSA")
	}

	// Parse the public key in PKIX format
	publick, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := publick.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

// Устанавливаем ключи для тестирования в *asimencrypt.AsimEncryptStruct
func SetKey(encrypt *asimencrypt.AsimEncryptStruct) {
	encrypt.PublicServerKey = []byte(PublicServerKey)

	key, err := GetRSAPrivateKey([]byte(PrivateServerKey))
	if err != nil {
		fmt.Println("get RSA key failed GetRSAPrivateKey:", err)
		return
	}
	encrypt.PrivateServerKey = key

	keyPrivateClientKey, err := GetRSAPrivateKey([]byte(PrivateClientKey))
	if err != nil {
		fmt.Println("get RSA key failed GetRSAPrivateKey:", err)
		return
	}
	encrypt.PrivateKey = keyPrivateClientKey

	keyPublicClientKey, err := GetRSAPublicKey([]byte(PublicClientKey))
	if err != nil {
		fmt.Println("get RSA key failed GetRSAPublicKey:", err)
		return
	}
	encrypt.PublicKey = keyPublicClientKey
}

func SetLogger(data *model.Data) {
	log, _, err := keeperlog.NewContextLogger("", false)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	data.Log = log
}

var server *httptest.Server
var encrypt *asimencrypt.AsimEncryptStruct

// Создаем тестовый сервер
func SetServerCardList(t *testing.T) {

	// Create a test server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
			return
		}

		if r.URL.Path != "/list/cards" {
			t.Errorf("Expected /cardslist path, got %s", r.URL.Path)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
			return
		}

		if r.Header.Get("Token") != "test-token" {
			t.Errorf("Expected Token: test-token, got %s", r.Header.Get("Token"))
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
			return
		}

		body, err = encrypt.DecryptOAEPServer(body)

		if err != nil {
			t.Errorf("Failed to DecryptOAEP: %v", err)
			return
		}

		var payload struct {
			Type string `json:"type"`
		}
		err = json.Unmarshal(body, &payload)
		if err != nil {
			t.Errorf("Failed to unmarshal request body: %v", err)
			return
		}

		if payload.Type != "get_list_cards" {
			t.Errorf("Unexpected payload type: got %s", payload.Type)
			return
		}

		type ResponseLists struct {
			ID   string `json:"id"`
			Data string `json:"data"`
		}

		//Берем ключ клиента публичны и шифруем данные
		data1, err := encrypt.EncryptByClientKeyParts(string(`{"ID": "1", "Description": "Card 1", "Number": "1235", "Exp": "12/23", "Cvc": "123"}`), PublicClientKey)
		if err != nil {
			log.Println("asimencrypt failed to encrypt", err)
		}

		//Берем ключ клиента публичны и шифруем данные
		data2, err := encrypt.EncryptByClientKeyParts(string(`{"ID": "2", "Description": "Card 2", "Number": "3215", "Exp": "12/24", "Cvc": "456"}`), PublicClientKey)
		if err != nil {
			log.Println("asimencrypt failed to encrypt", err)
		}

		// Return response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		// Формирование ответа
		response := []ResponseLists{
			{ID: "1", Data: base64.StdEncoding.EncodeToString(data1)},
			{ID: "2", Data: base64.StdEncoding.EncodeToString(data2)},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Failed to encode response", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}))

}

func TestCardsList(t *testing.T) {

	//lipgloss.SetColorProfile(termenv.Ascii)

	encrypt = asimencrypt.NewAsimEncrypt()

	SetKey(encrypt)
	SetServerCardList(t)

	defer server.Close()

	// Create a resty client with the test server URL
	client := resty.New() //.SetBaseURL(server.URL)

	cfg, err := config.ParseConfig() //config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	data := model.InitModel()
	data.Config = cfg
	model.InitRequestHTTP(data)
	cardslist := data.RequestHTTP["cardslist"]

	cardslist.URL = server.URL + "/list/cards"
	data.User.Token = "test-token"
	data.RequestHTTP["cardslist"] = cardslist
	data.ModeTest = true

	handler := RequestHTTP{
		client:      client,
		data:        data,
		asimencrypt: encrypt,
	}

	msg := handler.CardsList()
	if _, ok := msg.(errMsg); ok {
		t.Errorf("CardsList returned an error")
	}

	if status, ok := msg.(statusMsg); ok {
		if status != http.StatusOK {
			t.Errorf("Expected status OK, got %d", status)
		}
	}
	cardslist = data.RequestHTTP["cardslist"]
	modelTea := opcardslist.NewOpCardsList(data)

	pathFile := "testdata/cardlist/menu.out"
	dataCompare, err := readFromFile(pathFile)
	if err != nil {
		fmt.Println(err)
	}

	tm := teatest.NewTestModel(t, modelTea, teatest.WithInitialTermSize(300, 100))

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {

		return len(bts) > 0
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("ctrl+c"),
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second*5)))
	if err != nil {
		t.Error(err)
	}

	err = writeToFile(pathFile, out)
	if err != nil {
		fmt.Println("error write to file: ", err)
	}

	assert.Equal(t, out, dataCompare)

	pathFile = "testdata/cardlist/select.id.2.out"
	if fileExists(pathFile) {
		dataCompare, err = readFromFile(pathFile)
		if err != nil {
			fmt.Println(err)
		}
	}

	tm = teatest.NewTestModel(t, modelTea, teatest.WithInitialTermSize(300, 100))

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("down"),
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {

		return len(bts) > 0
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	tm.Quit()

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	outCard, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second*5)))
	if err != nil {
		t.Error(err)
	}

	err = writeToFile(pathFile, outCard)
	if err != nil {
		fmt.Println("error write to file: ", err)
	}

	assert.Equal(t, outCard, dataCompare)
}

// Функция для проверки существования файла
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Функция для записи данных в файл
func writeToFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// Функция для чтения данных из файла
func readFromFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}
