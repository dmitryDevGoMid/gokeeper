package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/config"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/middlewares"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/login"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/opcardslist"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/opcreditcard"
	opdescritption "github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/opdescription"
	opfileslist "github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/opfilelist"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/opfiles"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/opfilessenddialog"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/oppasswords"
	oppaswordslist "github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/oppasswordslist"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/options/registration"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/checkrunapp"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/keeperlog"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/requesthttp"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/showresponse"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/submenu"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/submenu/cardsmenu"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/submenu/filesmenu"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/submenu/passwordsmenu"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/topmenu"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

var data *model.Data
var client *resty.Client
var encrypt asimencrypt.AsimEncrypt

func StartStep() {
	data.NextStep.NextStepByName = "requesthttp"
	data.NextStep.RequestByName = "exchangeget"
}

type LogField struct {
	Action  string
	Request string
	UserID  string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.ParseConfig() //config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	//Проверка, того что приложение уже запущено.
	next, err := checkrunapp.StartCheckRunApp(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if next {
		fmt.Println("Приложение уже запущено, если вы уверены что нет, то повторите попытку через 1 минуту...")
		os.Exit(0)
	}

	data = model.InitModel()
	data.Config = cfg
	model.InitRequestHTTP(data)
	StartStep()

	log, file, err := keeperlog.NewContextLogger("", true)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	data.Log = log

	defer file.Close()

	encrypt = asimencrypt.NewAsimEncrypt()
	errEncrypt := encrypt.GenerateKeyFile("keeper")
	if errEncrypt != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	encrypt.AllSet()

	client = resty.New()

	//Выполняем повоторный запрос в случае получаения 401 - не авторизован
	//Предварительно выполняет авторизацию с помощью refresh token
	client.SetRetryCount(1)
	client.AddRetryCondition(
		// RetryConditionFunc type is for retry condition function
		// input: non-nil Response OR request execution error
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == 401
		},
	)

	middlewares := middlewares.NewClientMiddleware(client, data)
	middlewares.SetTokenToRequest(data, encrypt)
	middlewares.CheckStatusCode401(data)

	//middlewares.DecryptResponse(encrypt)

	client.SetRetryWaitTime(1 * time.Second)

	var modelTea tea.Model
	for {

		switch data.NextStep.NextStepByName {

		case "topmenu":
			modelTea = topmenu.NewTopMenu(data)

		case "registration":
			modelTea = registration.NewRegistrate(data)

		case "login":
			modelTea = login.NewLogin(data)

		case "requesthttp":
			modelTea = requesthttp.NewRequestHTTP(client, data, encrypt)

		case "showresponse":
			modelTea = showresponse.NewShowResponse(data)

		case "submenu":
			subMain()

		default:
			os.Exit(1)
		}
		if _, err := tea.NewProgram(modelTea, tea.WithAltScreen()).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	}
}

func subMain() {

	var modelTea tea.Model

	for {
		//fmt.Println(data)

		switch data.NextStep.NextStepByName {
		case "submenu":
			modelTea = submenu.NewSubMenu(data)

		case "passwords":
			subMainPasswords()

		case "files":
			subMainFiles()

		case "cards":
			subMainCards()

		case "topmenu":
			return

		default:
			os.Exit(1)
		}

		if _, err := tea.NewProgram(modelTea, tea.WithAltScreen()).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	}
}

func subMainCards() {

	var modelTea tea.Model

	for {
		//fmt.Println(data)

		switch data.NextStep.NextStepByName {

		//Меню по работе с паролями
		case "cards":
			modelTea = cardsmenu.NewCardsMenu(data)

		case "newcreditcards":
			modelTea = opcreditcard.NewCreditCard(data)

		case "description":
			modelTea = opdescritption.NewOpDescription(data, "")

			//Общие для всех меню
		case "requesthttp":
			modelTea = requesthttp.NewRequestHTTP(client, data, encrypt)

		case "showcardslist":
			modelTea = opcardslist.NewOpCardsList(data)

		case "showresponse":
			modelTea = showresponse.NewShowResponse(data)

		case "submenu":
			return

		default:
			os.Exit(1)
		}

		if _, err := tea.NewProgram(modelTea, tea.WithAltScreen()).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	}
}

func subMainPasswords() {

	var modelTea tea.Model

	for {
		//fmt.Println(data)

		switch data.NextStep.NextStepByName {

		//Меню по работе с паролями
		case "passwords":
			modelTea = passwordsmenu.NewPasswordsMenu(data)

		case "passwordslist":
			modelTea = passwordsmenu.NewPasswordsMenu(data)

		case "newpasswords":
			modelTea = oppasswords.NewOpPasswords(data)

		case "showpasswordslist":
			modelTea = oppaswordslist.NewOpPasswordsList(data)

		//Общие для всех меню
		case "requesthttp":
			modelTea = requesthttp.NewRequestHTTP(client, data, encrypt)

		case "description":
			modelTea = opdescritption.NewOpDescription(data, "")

		case "showresponse":
			modelTea = showresponse.NewShowResponse(data)

		case "submenu":
			return

		default:
			os.Exit(1)
		}

		if _, err := tea.NewProgram(modelTea, tea.WithAltScreen()).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	}
}

func subMainFiles() {
	var modelTea tea.Model

	for {
		//fmt.Println(data)

		switch data.NextStep.NextStepByName {

		case "showfileslist":
			modelTea = opfileslist.NewOpFilesList(data)

		//Меню по работе с файлами
		case "opfilessenddialog":
			modelTea = opfilessenddialog.NewOpFilesDialog(data)

		case "files":
			modelTea = filesmenu.NewFilesMenu(data)

		case "selectfiles":
			data.OptionFiles.QuttingOfSelectFiles = false
			modelTea = opfiles.NewOpFiles(data)

		//Общие для всех меню
		case "requesthttp":
			modelTea = requesthttp.NewRequestHTTP(client, data, encrypt)

		case "description":
			modelTea = opdescritption.NewOpDescription(data, "")

		case "showresponse":
			modelTea = showresponse.NewShowResponse(data)

		case "submenu":
			return

		default:
			os.Exit(1)
		}

		if _, err := tea.NewProgram(modelTea, tea.WithAltScreen()).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	}
}
