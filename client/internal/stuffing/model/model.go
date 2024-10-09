package model

import (
	"fmt"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/config"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/keeperlog"
)

type Data struct {
	NextStep        NextStep
	DescriptionStep string
	User            *User
	RequestHTTP     map[string]RequestHTTP
	Keys            *ServerExchangeKeys
	OptionPasswords OptionPasswords
	OptionFiles     OptionFiles
	OptionCards     OptionCards
	OptionTexts     OptionTexts
	Cancel          chan struct{}
	Log             *keeperlog.ContextLogger
	ModeTest        bool
	Config          *config.Config
}

//  Данные по текстам
type OptionTexts struct {
	Description string `json:"description"`
	Text        string `json:"text"`
}

// Данные по картам
type OptionCards struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Number      string `json:"number"`
	Exp         string `json:"exp"`
	Cvc         string `json:"cvc"`
	Edite       bool   `json:"edite"`
	Delete      bool   `json:"delete"`
}

// Данные по выбранным файлам
type OptionFiles struct {
	ID                   string `json:"_id"`
	ID_File              string `json:"id_file,omitempty"`
	Description          string `json:"description,omitempty"`
	SelectedFile         string `json:"selected_file"`
	QuttingOfSelectFiles bool   `json:"qutting_of_select_files"`
	Filename             string `json:"filename" bson:"filename"`
	ChunkSize            int64  `json:"chunk_size" bson:"chunkSize"`
	//Metadata   string             `json:"metadata" bson: "metadata"`
	Metadata struct {
		ClientID  string `bson:"client_id"`
		CountPart int    `bson:"count_part"`
		UID       string `bson:"uid"`
	} `json:"metadata" bson:"metadata"`
	UploadDate time.Time `json:"upload_date" bson:"uploadDate"`
	Length     int64     `json:"length" bson:"length"`
}

// Данные по пароля и логинам
type OptionPasswords struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Edite       bool   `json:"edite"`
}

type ServerExchangeKeys struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type ResponseLists struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type NextStep struct {
	NextStepByName  string
	RequestByName   string
	ShowErrorByName string
}

type User struct {
	ID           string `json:"id"`
	IDUser       string `json:"id_user"`
	Username     string `json:"username"`
	HashPassword string `json:"hash_password"`
	Password     string `json:"password"`
	Token        string `json:"token"`
	TokenRefresh string `json:"token_refresh"`
	Auth         bool   `json:"auth"`
}

type RequestHTTP struct {
	Name                    string `json:"name"`
	Type                    string `json:"type"`
	URL                     string `json:"url"`
	Request                 string `json:"request"`
	Response                Response
	ResponseLists           *[]ResponseLists
	PasswordsList           *[]OptionPasswords
	FilesList               *[]OptionFiles
	CardsList               *[]OptionCards
	OutputFileName          string `json:"output_file_name"`
	OutputFileCopyOrRewrite bool   `json:"output_file_copy_or_rewrite"` // false - copy, true - rewrite
	Choced                  string `json:"choced"`
}

type Response struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
	Body       string `json:"response"`
}

func InitModel() *Data {
	data := &Data{}
	data.NextStep.NextStepByName = "topmenu"
	data.Keys = &ServerExchangeKeys{}
	data.User = &User{}

	return data
}

func InitRequestHTTP(data *Data) {
	requests := make(map[string]RequestHTTP)
	domain := data.Config.Server.Address

	address := domain + "/api/user"
	requests["register"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "register", Type: "POST", URL: fmt.Sprintf("%s/%s/", address, "register")}
	requests["login"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "login", Type: "POST", URL: fmt.Sprintf("%s/%s", address, "login")}

	addressExchange := domain + "/api/user/exchange"
	requests["exchangeset"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "exchangeset", Type: "POST", URL: fmt.Sprintf("%s/%s", addressExchange, "set")}
	requests["exchangeget"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "exchangeget", Type: "POST", URL: fmt.Sprintf("%s/%s", addressExchange, "get")}

	addressSave := domain + "/api/user/save"
	requests["newcard"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "newcard", Type: "POST", URL: fmt.Sprintf("%s/%s", addressSave, "card")}
	requests["newpassword"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "newpassword", Type: "POST", URL: fmt.Sprintf("%s/%s", addressSave, "password")}
	requests["newtext"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "newtext", Type: "POST", URL: fmt.Sprintf("%s/%s", addressSave, "text")}

	addressList := domain + "/api/user/list"
	requests["cardslist"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "cardslist", Type: "POST", URL: fmt.Sprintf("%s/%s", addressList, "cards")}
	requests["passwordslist"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "passwordslist", Type: "POST", URL: fmt.Sprintf("%s/%s", addressList, "passwords")}
	requests["fileslist"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "fileslist", Type: "POST", URL: fmt.Sprintf("%s/%s", addressList, "files")}

	deleteFile := domain + "/api/user/delete"
	requests["deletefile"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "deletefile", Type: "POST", URL: fmt.Sprintf("%s/%s", deleteFile, "file")}
	requests["deletecard"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "deletecard", Type: "POST", URL: fmt.Sprintf("%s/%s", deleteFile, "card")}
	requests["deletepassword"] = RequestHTTP{ResponseLists: &[]ResponseLists{}, Name: "deletepassword", Type: "POST", URL: fmt.Sprintf("%s/%s", deleteFile, "password")}

	data.RequestHTTP = requests
}
