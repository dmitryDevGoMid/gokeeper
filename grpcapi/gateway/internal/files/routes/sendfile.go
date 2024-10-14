package routes

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/files/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const chunkSize = 1024 * 1024 // 1 MB

func setFolder(name string) string {
	folderPath := fmt.Sprintf("/%s", name)

	d, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	dir := filepath.Dir(folderPath)

	err = os.MkdirAll(d+dir, 0775)

	if err != nil {
		fmt.Println(err)
	}

	filename := d + folderPath

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	if err != nil {
		fmt.Println(err)
	}

	f.Close()

	return filename
}

var StreamMap map[string]pb.DataStreamer_SendFilesClient

func SendFiles(c *gin.Context, client pb.DataStreamerClient) error {

	fmt.Println("SIZE MAP:", len(StreamMap))

	checkStatus := c.GetHeader("Status-Send-File")

	// Получаем номер части из заголовка X-Part
	part_file := c.GetHeader("X-Part")
	fmt.Println("Получили номер файла: ", part_file)

	// Получаем имя файла из заголовка X-Filename
	filename := c.GetHeader("X-Filename")

	clientid := c.GetHeader("Client-ID")
	uidfile := c.GetHeader("Key-ID")
	countPart := c.GetHeader("Count-Part")
	checkSum := c.GetHeader("Check-Sum")
	fileName := c.GetHeader("File-Name")
	fmt.Println("FileName: ", fileName)

	if _, ok := StreamMap[uidfile]; !ok {
		stream, err := client.SendFiles(context.Background())
		if err != nil {
			fmt.Println("Create stream 5=>", err)
			return err
		}
		StreamMap[uidfile] = stream
		go GetByStreamDataFromServer(stream)
	}

	stream := StreamMap[uidfile]

	if checkStatus == "abort" {
		err := SendAbotr(stream, uidfile)

		if err != nil {
			fmt.Println("Stream Send Abort:", filename, "10=>", err)
			return err
		}

		err = StreamClose(stream, uidfile)
		if err != nil {
			fmt.Println("Stream Close:", filename, "10=>", err)
			return err
		}
		fmt.Println("Stream Close:", uidfile, "ClientID:", clientid, "filename:", filename)
		return nil
	}

	buffer, err := io.ReadAll(c.Request.Body)

	fmt.Println("Размер данных:", len(buffer))
	fmt.Println("Ошибка при чтении данных:", err)
	fmt.Println("checkStatus:", checkStatus)

	if err != nil {
		fmt.Println(filename, "7=>", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return nil
	}

	streamClose := false

	if len(buffer) == 0 && checkStatus == "end" {

		fmt.Println("Закрываем поток!")
		streamClose = true

		c.JSON(200, gin.H{"message": "File received"})
	}

	if !streamClose {
		if err := stream.Send(&pb.FilesRequest{Filename: fileName, Checksum: checkSum, Data: buffer, Countpart: countPart, Uidfile: uidfile, Clientid: clientid, Numberpart: part_file}); err != nil {
			fmt.Println(filename, "9=>", err)
			return err
		}
	} else {

		err := StreamClose(stream, uidfile)
		if err != nil {
			fmt.Println(filename, "10=>", err)
			return err
		}
	}
	return nil
}

func SendAbotr(stream pb.DataStreamer_SendFilesClient, uidfile string) error {
	if err := stream.Send(&pb.FilesRequest{Abort: true, Uidfile: uidfile}); err != nil {
		fmt.Println(uidfile, "9=>", err)
		return err
	}

	return nil
}

func StreamClose(stream pb.DataStreamer_SendFilesClient, uidfile string) error {
	fmt.Println("Закрываем этот поток!")
	err := stream.CloseSend()
	if err != nil {
		fmt.Println(uidfile, "99=>", err)
		return err
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Удалили этот поток!")
	delete(StreamMap, uidfile)
	return nil
}

func GetByStreamDataFromServer(stream pb.DataStreamer_SendFilesClient) error {

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Canceled {
				break
			}
			log.Printf("failed to receive response: %v", err)
			return err
		}
		log.Printf("Received response: %s", resp.Message)
	}

	log.Printf("Close Stream Stop End Done....")
	return nil
}
