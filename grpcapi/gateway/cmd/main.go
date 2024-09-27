package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/config"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/files"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/files/pb"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/files/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	routes.StreamMap = make(map[string]pb.DataStreamer_SendFilesClient)

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	// Установка режима релиза
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	files.RegisterRoutes(r, cfg)

	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: r,
	}

	fmt.Println("cfg.Server.Address:", cfg.Server.Address)
	fmt.Println("cfg.GrpcServerAdress.AddressGrpc:", cfg.GrpcServerAdress.AddressGrpc)
	fmt.Println("cfg.DataBase.DatabaseURL:", cfg.DataBase.DatabaseURL)

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
	log.Println("Server exiting")
}
