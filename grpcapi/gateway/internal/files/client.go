package files

import (
	"log"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/config"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/files/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Client pb.DataStreamerClient
}

func InitServiceClient(c *config.Config) pb.DataStreamerClient {

	conn, err := grpc.NewClient(c.GrpcServerAdress.AddressGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return pb.NewDataStreamerClient(conn)
}
