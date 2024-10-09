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
	opdescritption "github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/opdescription"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/oppasswords"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/showresponse"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/go-playground/assert/v2"

	"github.com/go-resty/resty/v2"
	"github.com/muesli/termenv"
)

var serverPasswordNew *httptest.Server

// Создаем тестовый сервер
func SetServerPasswordNew(t *testing.T) {

	// Create a test server
	serverPasswordNew = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
			return
		}

		if r.URL.Path != "/save/password" {
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

func TestPasswordNew(t *testing.T) {

	lipgloss.SetColorProfile(termenv.Ascii)

	encrypt = asimencrypt.NewAsimEncrypt()

	SetKey(encrypt)
	SetServerPasswordNew(t)

	defer serverPasswordNew.Close()

	// Create a resty client with the test server URL
	client := resty.New()

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	data := model.InitModel()
	data.Config = cfg
	model.InitRequestHTTP(data)
	mPasswdNew := data.RequestHTTP["newpassword"]

	mPasswdNew.URL = serverPasswordNew.URL + "/save/password"
	data.User.Token = "test-token"
	data.RequestHTTP["newpassword"] = mPasswdNew

	data.OptionPasswords.Username = "test-username"
	data.OptionPasswords.Password = "test-password"
	data.OptionPasswords.Description = "my new password by site"

	handler := RequestHTTP{
		client:      client,
		data:        data,
		asimencrypt: encrypt,
	}

	SetLogger(data)

	msg := handler.PasswordNew()
	if _, ok := msg.(errMsg); ok {
		t.Errorf("CardsList returned an error")
	}

	if status, ok := msg.(statusMsg); ok {
		if status != http.StatusCreated {
			t.Errorf("Expected status OK, got %d", status)
		}
	}
	mPasswdNew = data.RequestHTTP["newpassword"]
	assert.Equal(t, 201, mPasswdNew.Response.StatusCode)

	modelTea := oppasswords.NewOpPasswords(data)

	modelTea.SetUsername(data.OptionPasswords.Username)
	modelTea.SetPassword(data.OptionPasswords.Password)

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

	pathFile := "testdata/newpassword/menu.out"
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

	//NewShowResponse
	modelTeaShowResponse := showresponse.NewShowResponse(data)

	data.NextStep.NextStepByName = "showresponse"
	data.NextStep.RequestByName = "newpassword"

	tmSh := teatest.NewTestModel(t, modelTeaShowResponse, teatest.WithInitialTermSize(300, 100))

	teatest.WaitFor(t, tmSh.Output(), func(bts []byte) bool {
		return len(bts) > 0
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	tmSh.Quit()

	tmSh.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	outResult, err := io.ReadAll(tmSh.FinalOutput(t, teatest.WithFinalTimeout(time.Second)))
	if err != nil {
		t.Error(err)
	}

	pathFile = "testdata/newpassword/response.out"
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

	//Description
	modelTeaDescription := opdescritption.NewOpDescription(data, data.OptionPasswords.Description)

	data.NextStep.RequestByName = "newpassword"

	tmDescr := teatest.NewTestModel(t, modelTeaDescription, teatest.WithInitialTermSize(300, 100))

	teatest.WaitFor(t, tmDescr.Output(), func(bts []byte) bool {
		return len(bts) > 0
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	tmDescr.Quit()

	tmDescr.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	outResultDescription, err := io.ReadAll(tmDescr.FinalOutput(t, teatest.WithFinalTimeout(time.Second)))
	if err != nil {
		t.Error(err)
	}

	pathFile = "testdata/newpassword/description.out"
	if fileExists(pathFile) {
		dataCompare, err = readFromFile(pathFile)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = writeToFile(pathFile, outResultDescription)
	if err != nil {
		fmt.Println("error write to file: ", err)
	}

	assert.Equal(t, outResultDescription, dataCompare)
}
