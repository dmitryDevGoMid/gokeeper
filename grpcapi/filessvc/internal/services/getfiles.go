package services

import (
	"context"
	"fmt"

	pb "github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/pb"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/pkg/calcsum"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/repository/repmongo"
)

func (s *Server) GetFile(req *pb.FileRequest, stream pb.DataStreamer_GetFileServer) error {
	fileID := req.GetFileId()     // Получаем уникальный номер файла, который предварительно отдали клиент в списке названий файлов
	cleintID := req.GetClientId() // Сам айди клиента для дополнительной фильтрации

	calcDataChan, calcGetSum, calcDone := calcsum.RunCalcSum()
	defer close(calcDataChan)
	defer close(calcDone)

	chanData := make(chan []byte)
	done := make(chan struct{})

	go func(ctx context.Context) error {
		for {
			select {
			case <-stream.Context().Done():
				fmt.Println("stream.Context().Done(): cancelled")
				return stream.Context().Err()
			case <-done:
				calcDone <- struct{}{}
				return nil
			case data := <-chanData:
				chunk := &pb.FileChunk{
					Data: data,
				}
				if err := stream.Send(chunk); err != nil {
					fmt.Println("if err := stream.Send(chunk); err != nil ", err)
					return nil
				}
				//Отправляем данные в канал для подсчета контрольной суммы
				calcDataChan <- chunk.Data
			}
		}
	}(stream.Context())

	fmt.Println("File ID:", fileID, "Client ID:", cleintID)
	repm := repmongo.NewRepMongo(s.MongoDB)

	//Запускаем запрос к базе и отпраку файлов в канал для передачи в поток клиенту
	err := repm.GetChunkByIDFile(fileID, chanData)
	if err != nil {
		fmt.Println("GetChunkByIDFile===>", err)
		return err
	}

	close(done)

	//Получаем контрольную сумму
	checkSumCalcByData := <-calcGetSum

	fmt.Println("checkSumCalcByData===>", checkSumCalcByData)

	return nil
}
