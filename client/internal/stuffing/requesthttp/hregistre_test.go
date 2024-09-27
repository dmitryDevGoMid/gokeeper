package requesthttp

import (
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
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/registration"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/showresponse"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/go-playground/assert/v2"

	"github.com/go-resty/resty/v2"
	"github.com/muesli/termenv"
)

var serverRegistre *httptest.Server

// Создаем тестовый сервер
func SetServerRegistre(t *testing.T) {

	// Create a test server
	serverRegistre = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
			return
		}

		if r.URL.Path != "/user/register" {
			t.Errorf("Expected /login path, got %s", r.URL.Path)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
			return
		}

		_, err = encrypt.DecryptOAEPServer(body)

		if err != nil {
			t.Errorf("Failed to DecryptOAEP: %v", err)
			return
		}

		// Return response
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		// Формирование ответа
		if err := json.NewEncoder(w).Encode(""); err != nil {
			log.Println("Failed to encode response", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}))

}

func TestRegister(t *testing.T) {

	lipgloss.SetColorProfile(termenv.Ascii)

	encrypt = asimencrypt.NewAsimEncrypt()

	SetKey(encrypt)
	SetServerRegistre(t)

	defer serverRegistre.Close()

	// Create a resty client with the test server URL
	client := resty.New() //.SetBaseURL(server.URL)

	cfg, err := config.ParseConfig() //config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	data := model.InitModel()
	data.Config = cfg
	model.InitRequestHTTP(data)
	mRegister := data.RequestHTTP["register"]

	mRegister.URL = serverRegistre.URL + "/user/register"
	data.User.Token = "test-token"
	data.RequestHTTP["register"] = mRegister

	data.User.Username = "test-username"
	data.User.Password = "test-password"

	handler := RequestHTTP{
		client:      client,
		data:        data,
		asimencrypt: encrypt,
	}

	msg := handler.Register()
	if _, ok := msg.(errMsg); ok {
		t.Errorf("CardsList returned an error")
	}

	if status, ok := msg.(statusMsg); ok {
		if status != http.StatusCreated {
			t.Errorf("Expected status OK, got %d", status)
		}
	}
	mRegister = data.RequestHTTP["register"]
	assert.Equal(t, 201, mRegister.Response.StatusCode)

	modelTea := registration.NewRegistrate(data)

	modelTea.SetUsername(data.User.Username)
	modelTea.SetPassword(data.User.Password)

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

	pathFile := "testdata/register/menu.out"
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

	SetLogger(data)
	modelTeaShowResponse := showresponse.NewShowResponse(data)

	data.NextStep.NextStepByName = "showresponse"
	data.NextStep.RequestByName = "register"

	tm = teatest.NewTestModel(t, modelTeaShowResponse, teatest.WithInitialTermSize(300, 100))

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return len(bts) > 0
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	tm.Quit()

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	outResult, err := io.ReadAll(tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second)))
	if err != nil {
		t.Error(err)
	}

	pathFile = "testdata/register/response.out"
	if fileExists(pathFile) {
		dataCompare, err = readFromFile(pathFile)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = writeToFile(pathFile, outResult)
	if err != nil {
		fmt.Println("error write to file: ", err)
	}

	assert.Equal(t, outResult, dataCompare)
}
