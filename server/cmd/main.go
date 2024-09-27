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
	//gin.SetMode(gin.ReleaseMode)

	// Создание маршрутизатора Gin без настроек по умолчанию
	r := gin.Default()
	//r := gin.New()

	//r.Use(middlewares.DecryptMiddleware(encrypt))

	r.Use(middlewares.DecryptMiddleware(encrypt))
	//r.Use(middlewares.CheckToken())
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

	//go HttpServer(mongodb, encrypt)

	r.Run(":8000")

	/*srv := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	// Запуск сервера в отдельной горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Ожидание сигналов для завершения работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Установка таймаута для завершения работы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exiting")*/

}

/*func HttpServer(mongodb db.MongoDBClient, encrypt asimencrypt.AsimEncrypt) {
	// Создаем HTTP сервер


	userHandler := userHandlers.NewIHandler(mongodb, encrypt)

	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	r.POST("/api/user/login/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	r.POST("/api/user/log/", userHandler.Login)

	server := &http.Server{
		Addr:    ":9000",
		Handler: r,
	}

	// Канал для получения сигналов завершения работы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :9000: %v\n", err)
		}
	}()

	// Ждем сигнала завершения работы
	<-stop
	log.Println("Shutting down server...")

	// Создаем контекст с таймаутом для плавного завершения работы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Завершаем работу сервера
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited properly")
}*/
