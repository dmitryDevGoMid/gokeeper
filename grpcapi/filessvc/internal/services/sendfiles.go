package services

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/config/dbmongo"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/config/dbredis"
	pb "github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/pb"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/pkg/calcsum"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/repository/repmongo"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/repository/repredis"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	MongoDB *dbmongo.MongoDBClient
	RedisDB dbredis.RedisDBClient
	pb.UnimplementedDataStreamerServer
}

type CheckWrite struct {
	ClientId string
	Key      string
	Part     string
	Write    bool
}

func (s *Server) SendFiles(stream pb.DataStreamer_SendFilesServer) error {

	checkSumByStream := ""
	calcDataChan, calcGetSum, calcDone := calcsum.RunCalcSum()
	defer close(calcDataChan)
	defer close(calcDone)

	repm := repmongo.NewRepMongo(s.MongoDB)

	reprd := repredis.NewRepRedis(s.RedisDB.GetNewClient(), repm)

	checkWrite := make(map[string]repredis.CheckWrite)

	start := time.Now()

	fmt.Println("Stream...")

	checkAbort := make(chan struct{}, 1)
	bufferChanRedis := make(chan *repredis.Chunk, 10)
	errChan := make(chan error)

	defer close(checkAbort)
	defer close(bufferChanRedis)
	defer close(errChan)

	mutexRedis := new(sync.Mutex)

	go reprd.WriteToDB(mutexRedis, bufferChanRedis, errChan, checkAbort)

	var data []byte

	for {
		select {
		case err := <-errChan:
			fmt.Println(err)
			return err

		default:
			chunk, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("io.EOF STOPPPP ")
				fmt.Println("Данные по клиенту!", chunk)

				calcDone <- struct{}{}
				checkSumCalcByData := <-calcGetSum
				fmt.Println("Stream Sum: ", checkSumByStream)
				fmt.Println("Calc Sum: ", checkSumCalcByData)
				if !CompareCheckSum(checkSumByStream, checkSumCalcByData) {
					checkAbort <- struct{}{}
					fmt.Println("error: checkSum by Stream failed to calculate sum")
					return errors.New("checkSum by Stream failed to calculate sum")
				}

				return nil
			}

			if err != nil {
				checkAbort <- struct{}{}
				fmt.Println("ERRRRR STOPPP")
				return err
			}

			if chunk.Abort {
				fmt.Println("Получили аборт!")
				checkAbort <- struct{}{}
				return nil
			}

			if len(chunk.Data) > 0 {
				//Отправляем данные в канал для расчета checksum
				calcDataChan <- chunk.Data
				checkSumByStream = chunk.Checksum
				fmt.Println("checkSumByStream:", checkSumByStream)

				fmt.Println("ОТПРАВИЛИ В ПОТОК ")
				chunkToChannel := &repredis.Chunk{FileName: chunk.Filename, Check: &checkWrite, Start: start, Buffer: chunk.Data, Uidfile: chunk.Uidfile, Cleintid: chunk.Clientid, CountPart: chunk.Countpart, NumberPart: chunk.Numberpart}
				bufferChanRedis <- chunkToChannel
			}

			// Отправляем ответ клиенту в потоке
			if err := stream.Send(&pb.FilesResponse{
				Message:    fmt.Sprintf("Received part %s", chunk.Numberpart),
				Clientid:   chunk.Clientid,
				Uidfile:    chunk.Uidfile,
				Numberpart: chunk.Numberpart,
				Countpart:  chunk.Countpart,
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to send response: %v", err)
			}

			data = append(data, chunk.Data...)
			fmt.Printf("Received %d bytes\n", len(data))
		}

	}
}

// Выполняем сравнение сумм полученные от клиента и расчитанной на основании данных из потока
func CompareCheckSum(checkSumStream string, checkSumCalcByData string) bool {
	return checkSumStream == checkSumCalcByData
}
