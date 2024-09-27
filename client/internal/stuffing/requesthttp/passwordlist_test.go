package requesthttp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/config"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	oppaswordslist "github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/oppasswordslist"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/asimencrypt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/go-playground/assert/v2"

	//"github.com/go-playground/assert/v2"
	"github.com/go-resty/resty/v2"
	"github.com/muesli/termenv"
)

var serverPasswordList *httptest.Server

// Создаем тестовый сервер
func SetServerPasswordList(t *testing.T) {

	// Create a test server
	serverPasswordList = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
			return
		}

		if r.URL.Path != "/list/passwords" {
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

		if payload.Type != "get_list_passwords" {
			t.Errorf("Unexpected payload type: got %s", payload.Type)
			return
		}

		type ResponseLists struct {
			ID   string `json:"id"`
			Data string `json:"data"`
		}

		//Берем ключ клиента публичны и шифруем данные
		data1, err := encrypt.EncryptByClientKeyParts(string(`{"ID": "1", "Description": "password by site confetki.com", "Username": "aftor@gml.com", "Password": "Tgdhf_465%"}`), PublicClientKey)
		if err != nil {
			log.Println("asimencrypt failed to encrypt", err)
		}

		//Берем ключ клиента публичны и шифруем данные
		data2, err := encrypt.EncryptByClientKeyParts(string(`{"ID": "2", "Description": "password by sitew rabota.ru", "Username": "aftor@gml.com", "Password": "Tgdhf_981%"}`), PublicClientKey)
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

func TestPasswordsList(t *testing.T) {

	lipgloss.SetColorProfile(termenv.Ascii)

	encrypt = asimencrypt.NewAsimEncrypt()

	SetKey(encrypt)
	SetServerPasswordList(t)

	defer serverPasswordList.Close()

	// Create a resty client with the test server URL
	client := resty.New() //.SetBaseURL(server.URL)

	cfg, err := config.ParseConfig() //config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	data := model.InitModel()
	data.Config = cfg
	model.InitRequestHTTP(data)
	passwordslist := data.RequestHTTP["passwordslist"]

	passwordslist.URL = serverPasswordList.URL + "/list/passwords"
	data.User.Token = "test-token"
	data.RequestHTTP["passwordslist"] = passwordslist

	handler := RequestHTTP{
		client:      client,
		data:        data,
		asimencrypt: encrypt,
	}

	msg := handler.PasswordsList()
	if _, ok := msg.(errMsg); ok {
		t.Errorf("CardsList returned an error")
	}

	if status, ok := msg.(statusMsg); ok {
		if status != http.StatusOK {
			t.Errorf("Expected status OK, got %d", status)
		}
	}
	passwordslist = data.RequestHTTP["passwordslist"]
	//fmt.Println("=>", cardslist.CardsList)
	modelTea := oppaswordslist.NewOpPasswordsList(data)

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
	var dataCompare []byte
	pathFile := "testdata/passwordlist/menu.out"
	if fileExists(pathFile) {
		dataCompare, err = readFromFile(pathFile)
		if err != nil {
			fmt.Println(err)
		}
	}

	//fmt.Println(string(dataCompare))
	err = writeToFile(pathFile, out)
	if err != nil {
		fmt.Println("error write to file: ", err)
	}

	assert.Equal(t, out, dataCompare)

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

	out, err = io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second*5)))
	if err != nil {
		t.Error(err)
	}

	pathFile = "testdata/passwordlist/select.id.2.out"
	if fileExists(pathFile) {
		dataCompare, err = readFromFile(pathFile)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = writeToFile(pathFile, out)
	if err != nil {
		fmt.Println("error write to file: ", err)
	}

	assert.Equal(t, out, dataCompare)

}
