package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/config"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/config/db"
	userHandlers "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/handlers"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/middlewares"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/security"
	filesRepository "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/files"
	userRepository "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	encrypt := asimencrypt.NewAsimEncrypt()
	err = encrypt.GenerateKeyFile("keeper")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	encrypt.AllSet()

	// Установка режима релиза

	// Создание маршрутизатора Gin без настроек по умолчанию
	r := gin.Default()

	r.Use(middlewares.DecryptMiddleware(encrypt))
	r.Use(middlewares.ChaeckAndGetUserByToken(encrypt))

	mongodb := db.NewConnectMongoDB(cfg)

	repoUser := userRepository.NewUserRepository(mongodb)
	repoFile := filesRepository.NewFilesRepository(mongodb)

	security := security.NewSecurity()

	userHandler := userHandlers.NewHandler(repoUser, repoFile, encrypt, security)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	//Authenticate
	r.POST("api/user/register", userHandler.Register)
	r.POST("api/user/login", userHandler.Login)
	r.POST("api/user/refresh/token", userHandler.RefreshToken)
	r.POST("api/user/check/token", userHandler.CheckToken)

	//Exchange key
	r.POST("api/user/exchange/get", userHandler.ExchangeGet)
	r.POST("api/user/exchange/set", userHandler.ExchangeSet)

	//Save Endpoints
	r.POST("api/user/save/password", userHandler.PasswordSave)
	r.POST("api/user/save/card", userHandler.CardSave)

	r.POST("api/user/delete/card", userHandler.CardDelete)
	r.POST("api/user/delete/file", userHandler.DeleteFile)
	r.POST("api/user/delete/password", userHandler.PasswordDelete)

	//List EndPoints
	r.POST("api/user/list/passwords", userHandler.PasswordList)
	r.POST("api/user/list/files", userHandler.FilesList)
	r.POST("api/user/list/cards", userHandler.CardsList)

	r.POST("api/user/savekeys", userHandler.Login)
	r.PUT("api/user/:username/password", userHandler.UpdatePassword)
	r.DELETE("api/user/:username", userHandler.DeleteByUsername)

	r.Run(":8000")

}
