package repmongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/config/dbmongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Выносим через интерфейс доступные для вызова методы обьекта
type IRepmongo interface {
	CreateStreamForSave(cleintID string, mytex *sync.Mutex, fileName string, countPart int, uid string)
	SendPartFileToStream(data []byte, uid string, mytex *sync.Mutex) error
	StopStream(signal string, uid string, mytex *sync.Mutex) error
	GetChunkByIDFile(id string, chanelData chan []byte) error
}

type InitData struct {
	fileName  string
	countPart int
	uid       string
	clientID  string
	mutex     *sync.Mutex
}
type Chanels struct {
	data     chan []byte
	abort    chan struct{}
	done     chan struct{}
	initData *InitData
}

type Repmongo struct {
	MongoDB                *dbmongo.MongoDBClient
	ListChanelForWriteFile map[string]*Chanels
	ChunkSize              int32
	SyncToDelete           *sync.Mutex
}

func NewRepMongo(mongodb *dbmongo.MongoDBClient) *Repmongo {
	//Карта каналов для клиента на запись в базу
	listChan := make(map[string]*Chanels)

	//Размер чанка для записи в базу
	chunkSize := int32((1024 * 1024) * 2)

	//Синхронизация через mutex для удаления данных
	mutex := &sync.Mutex{}

	return &Repmongo{MongoDB: mongodb, ListChanelForWriteFile: listChan, ChunkSize: chunkSize, SyncToDelete: mutex}
}

func (rep *Repmongo) сreateChanels(clientID string, countPart int, uid string, fileName string) *Chanels {
	mutex := new(sync.Mutex)
	initData := &InitData{clientID: clientID, countPart: countPart, uid: uid, fileName: fileName, mutex: mutex}

	data := make(chan []byte)
	done := make(chan struct{})
	abort := make(chan struct{})

	return &Chanels{data: data, done: done, abort: abort, initData: initData}
}

func (rep *Repmongo) CreateStreamForSave(cleintID string, mytex *sync.Mutex, fileName string, countPart int, uid string) {

	mytex.Lock()
	defer mytex.Unlock()

	if _, ok := rep.ListChanelForWriteFile[uid]; !ok {
		chanels := rep.сreateChanels(cleintID, countPart, uid, fileName)
		rep.ListChanelForWriteFile[uid] = chanels
		go rep.writeFileToDB(chanels)
	}
}

// Stop stream to DB
func (rep *Repmongo) StopStream(signal string, uid string, mytex *sync.Mutex) error {

	mytex.Lock()
	defer mytex.Unlock()

	if value, exists := rep.ListChanelForWriteFile[uid]; exists {
		switch signal {
		case "done":
			value.done <- struct{}{}
		case "abort":
			value.abort <- struct{}{}
		default:
			return errors.New("sendPartFileToStream: not found type signal")
		}
	} else {
		return errors.New("sendPartFileToStream: not found key")
	}

	return nil

}

// Send data to stream to DB
func (rep *Repmongo) SendPartFileToStream(data []byte, uid string, mytex *sync.Mutex) error {

	mytex.Lock()
	defer mytex.Unlock()

	// Проверяем наличие ключа "four"
	if value, exists := rep.ListChanelForWriteFile[uid]; exists {
		value.data <- data
	} else {
		return errors.New("sendPartFileToStream: not found key")
	}

	return nil

}

// Gorutine to stream
func (rep *Repmongo) writeFileToDB(chanels *Chanels) error {

	//Получаем бакет для записи в базу файла
	bucket, err := rep.MongoDB.GetBucket()
	if err != nil {
		return err
	}
	// Создаем опции для загрузки файла
	uploadOpts := options.GridFSUpload().SetChunkSizeBytes(rep.ChunkSize)

	//Устанавливаем метаданные для файла
	uploadOpts.SetMetadata(bson.D{{"file_name", chanels.initData.fileName}, {"client_id", chanels.initData.clientID}, {"count_part", chanels.initData.countPart}, {"uid", chanels.initData.uid}})

	// Создаем поток для загрузки файла
	stream, err := bucket.OpenUploadStream(chanels.initData.fileName, uploadOpts)
	if err != nil {
		return err
	}

	defer stream.Close()
	chunkNumber := 0

	// Читаем данные из канала и записываем их в поток
	for {
		select {
		case data, ok := <-chanels.data:
			chunkNumber++
			if !ok {
				// Канал закрыт, завершаем запись файла
				if err := stream.Close(); err != nil {
					return err
				}
				fmt.Printf("New file uploaded with ID %d\n", chunkNumber)
				fmt.Println("File_ID_1:", stream.FileID)
				return nil
			}
			if _, err := stream.Write(data); err != nil {
				return err
			}
			fmt.Println("File_ID_2:", stream.FileID)
		case <-chanels.done:
			// Получен сигнал о завершении работы, завершаем запись файла, успешно завершаем запись и создание файла
			if err := stream.Close(); err != nil {
				return err
			}
			err := rep.deleteClientIntoListChannel(chanels.initData.uid)
			if err != nil {
				return err
			}
			fmt.Println("File_ID_3:", stream.FileID)
			fmt.Printf("Close Last chunkNumber %d\n", chunkNumber)
			return nil
		case <-chanels.abort:
			// Получен сигнал о завершении работы, завершаем запись файла обрываем запись файла, данные сохранены не будут
			if err := stream.Abort(); err != nil {
				return err
			}
			fmt.Println("File_ID_4:", stream.FileID)
			err := rep.deleteClientIntoListChannel(chanels.initData.uid)
			if err != nil {
				return err
			}
			fmt.Printf("Abort Last chunkNumber %d\n", chunkNumber)
			return nil
		}
	}
}

// Dalete stream
func (rep *Repmongo) deleteClientIntoListChannel(uid string) error {
	rep.SyncToDelete.Lock()
	defer rep.SyncToDelete.Unlock()

	if _, exists := rep.ListChanelForWriteFile[uid]; exists {

		delete(rep.ListChanelForWriteFile, uid)

	} else {
		return errors.New("sendPartFileToStream: not found key")
	}

	fmt.Println("Удалили стрим для записи в базу")

	return nil

}

func (rep *Repmongo) GetChunkByIDFile(id string, chanelData chan []byte) error {
	var err error

	// Укажите ID файла, для которого нужно получить чанки
	fileIDObject := returnHexById(id)

	// Определяем фильтр для поиска чанков по ID файла
	filter := bson.D{{Key: "files_id", Value: fileIDObject}}

	// Определяем опции для сортировки чанков по возрастанию
	sortOpts := options.Find().SetSort(bson.D{{Key: "n", Value: 1}})

	// Получаем коллекцию chunks
	chunksCollection := rep.MongoDB.GetCollectionFiles("custom name.chunks")

	// Находим чанки, соответствующие фильтру и сортируем их по возрастанию
	cursor, err := chunksCollection.Find(context.TODO(), filter, sortOpts)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer cursor.Close(context.TODO())

	fmt.Println("2GetChunkByIDFile:", id)

	// Обрабатываем найденные чанки
	for cursor.Next(context.TODO()) {
		fmt.Println("1cursor.Next(context.TODO()):", id)
		var chunk bson.M
		if err := cursor.Decode(&chunk); err != nil {
			log.Fatal(err)
			return err
		}
		// Извлекаем данные из поля data
		data, ok := chunk["data"].(primitive.Binary)
		if !ok {
			log.Fatal(err)
			return err
		}

		// Преобразуем данные в []byte
		dataBytes := data.Data
		fmt.Println("2cursor.Next(context.TODO()):", id)
		// Разбиваем данные на порции размером по 100 КБ и отправляем их в канал
		for i := 0; i < len(dataBytes); i += 100 * 1024 {
			time.Sleep(10 * time.Millisecond)
			end := i + 100*1024
			if end > len(dataBytes) {
				end = len(dataBytes)
			}
			chanelData <- dataBytes[i:end]
		}

		fmt.Printf("Chunk: %d; Size: %d bytes\n", chunk["n"], len(dataBytes))
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

func returnHexById(id string) primitive.ObjectID {
	//Обявляем перменную типа primitive.ObjectID
	var ID primitive.ObjectID

	// Выполняем преобразование в Hex - шестнадцатиричное кодирование
	ID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		panic(err)
	}

	//Возвращаем primitive.ObjectID
	return ID
}
