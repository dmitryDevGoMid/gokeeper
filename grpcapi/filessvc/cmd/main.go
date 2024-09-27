package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/config"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/config/dbmongo"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/config/dbredis"

	pb "github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/pb"

	services "github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/services"

	"google.golang.org/grpc"
)

func main() {

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	mongoDB := dbmongo.NewConnectMongoDB(cfg)
	redisDB := dbredis.NewConnectRedis(cfg)

	lis, err := net.Listen("tcp", cfg.GrpcServerPort.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Product Svc on", cfg.GrpcServerPort.Port)

	s := &services.Server{MongoDB: mongoDB, RedisDB: redisDB}

	grpcServer := grpc.NewServer()

	pb.RegisterDataStreamerServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
