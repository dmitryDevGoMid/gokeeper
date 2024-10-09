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
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/login"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/showresponse"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/go-playground/assert/v2"

	"github.com/go-resty/resty/v2"
	"github.com/muesli/termenv"
)

var serverLogin *httptest.Server

// Создаем тестовый сервер
func SetServerLogin(t *testing.T) {

	// Create a test server
	serverLogin = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
			return
		}

		if r.URL.Path != "/user/login" {
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

		userFind := &model.User{
			ID:           "1",
			IDUser:       "1",
			Username:     "egorka@gorka.com",
			HashPassword: "uwyfgweu%^54723siufygwio^&#VDFJFUWEFUUWY",
			Password:     "",
			Token:        "",
			TokenRefresh: "",
		}

		// Return response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		// Формирование ответа
		if err := json.NewEncoder(w).Encode(userFind); err != nil {
			log.Println("Failed to encode response", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}))

}

func TestLogin(t *testing.T) {

	lipgloss.SetColorProfile(termenv.Ascii)

	encrypt = asimencrypt.NewAsimEncrypt()

	SetKey(encrypt)
	SetServerLogin(t)

	defer serverLogin.Close()

	// Create a resty client with the test server URL
	client := resty.New()

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	data := model.InitModel()
	data.Config = cfg
	model.InitRequestHTTP(data)
	mLogin := data.RequestHTTP["login"]

	mLogin.URL = serverLogin.URL + "/user/login"
	data.User.Token = "test-token"
	data.RequestHTTP["login"] = mLogin

	handler := RequestHTTP{
		client:      client,
		data:        data,
		asimencrypt: encrypt,
	}

	msg := handler.Login()
	if _, ok := msg.(errMsg); ok {
		t.Errorf("CardsList returned an error")
	}

	if status, ok := msg.(statusMsg); ok {
		if status != http.StatusOK {
			t.Errorf("Expected status OK, got %d", status)
		}
	}
	mLogin = data.RequestHTTP["login"]
	modelTea := login.NewLogin(data)

	assert.Equal(t, "1", data.User.ID)
	assert.Equal(t, "egorka@gorka.com", data.User.Username)

	modelTea.SetUsername(data.User.Username)
	modelTea.SetPassword("password")

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

	pathFile := "testdata/login/menu.out"
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

	SetLogger(data)
	modelTeaShowResponse := showresponse.NewShowResponse(data)

	data.NextStep.NextStepByName = "showresponse"
	data.NextStep.RequestByName = "login"

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

	pathFile = "testdata/login/response.out"
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
