package routes

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/files/pb"

	"github.com/gin-gonic/gin"
)

func GetFile(c *gin.Context, client pb.DataStreamerClient, cancelDownloadFile *map[string]chan struct{}) error {

	clientId := c.GetHeader("Client-ID")
	uidFile := c.GetHeader("File-ID")
	sizeFile := c.GetHeader("File-Size")

	chanCancel := make(chan struct{})

	cdf := *cancelDownloadFile
	cdf[uidFile] = chanCancel

	req := &pb.FileRequest{
		ClientId: clientId,
		FileId:   uidFile,
	}

	ctx, cancel := context.WithCancel(context.Background())

	stream, err := client.GetFile(ctx, req)
	if err != nil {
		log.Fatalf("could not get file: %v", err)
	}

	// Устанавливаем заголовки для ответа
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+"file.file")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", sizeFile) // Явно устанавливаем Content-Length

	for {
		select {
		case <-chanCancel:
			cancel()
			delete(cdf, uidFile)
			return nil
		default:
			resp, err := stream.Recv()
			if err == io.EOF {
				delete(cdf, uidFile)
				cancel()
				return nil
			}
			if err != nil {
				fmt.Println("could not receive chunk:", err)
				cancel()
				return nil
			}
			c.Writer.Write(resp.Data)
		}

	}
}
