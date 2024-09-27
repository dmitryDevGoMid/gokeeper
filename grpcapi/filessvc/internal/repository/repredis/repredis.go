package repredis

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/filessvc/internal/repository/repmongo"

	"github.com/go-redis/redis/v8"
)

type IRepredis interface {
	Retry(attempts int, delay time.Duration) RetryDecorator
	ReadPartFile(start time.Time, uid string, countPart int, cleintid string, fileName string)
	WriteToDB(mutex *sync.Mutex, bufferChanRedis chan *Chunk, errChan chan error)
}

type Repredis struct {
	clientRedis *redis.Client
	repMongo    repmongo.IRepmongo
}

type CheckWrite struct {
	ClientId string
	Key      string
	Part     string
	Write    bool
}

type Chunk struct {
	Buffer     []byte
	FileName   string
	Cleintid   string
	Uidfile    string
	NumberPart string
	CountPart  string
	Start      time.Time
	Check      *map[string]CheckWrite
}

func NewRepRedis(clientRedis *redis.Client, repMongo repmongo.IRepmongo) *Repredis {
	return &Repredis{clientRedis: clientRedis, repMongo: repMongo}
}

type RetryDecorator func(func(int, *sync.WaitGroup, []byte, string, string, *sync.Mutex) error) func(int, *sync.WaitGroup, []byte, string, string, *sync.Mutex) error

func (rep *Repredis) Retry(attempts int, delay time.Duration) RetryDecorator {
	return func(f func(int, *sync.WaitGroup, []byte, string, string, *sync.Mutex) error) func(int, *sync.WaitGroup, []byte, string, string, *sync.Mutex) error {
		return func(partID int, wg *sync.WaitGroup, buffer []byte, uid string, clientid string, mytex *sync.Mutex) error {
			for i := 0; i < attempts; i++ {
				err := f(partID, wg, buffer, uid, clientid, mytex)
				if err == nil {
					return nil
				}
				if i == attempts-1 {
					return err
				}
				//fmt.Println("Retry->Ждем перед повторным запуском!", err, clientid, "; uid:", uid)
				time.Sleep(delay)
				//fmt.Println("Retry->Запустили заново сохранение в базе!", err, clientid, "; uid:", uid)
			}
			return nil
		}
	}
}

// Сохраняем данные в базу
func (rep *Repredis) SavePartFile(partID int, wg *sync.WaitGroup, buffer []byte, uid string, clientid string, mytex *sync.Mutex) error {

	start := time.Now()

	key := uid
	//Сохраняем данные в базе и помечаем их удаление через 60 сек.
	err := rep.clientRedis.Set(context.Background(), key, buffer, 10*60*time.Second).Err()
	if err != nil {
		elapsed := time.Since(start)
		seconds := elapsed.Seconds()
		fmt.Printf("Скрипт выполнялся %.2f секунд\n", seconds)
		fmt.Println("Ошибка выполнения скрипта записи в базу REDIS", err)
		return err
	}

	fmt.Println("Write to REDIS DB: ", key, " Client_ID:", clientid, " PartID:", partID, "uid: ", uid)
	wg.Done()
	return nil

}

func (rep *Repredis) StartChangeLogFileRedis() {
	time.AfterFunc(10*60*time.Second, func() {
		clientDBRedis := rep.clientRedis
		err := clientDBRedis.BgRewriteAOF(context.Background()).Err()
		if err != nil {
			fmt.Printf("Error executing BGREWRITEAOF: %v\n", err)
		} else {
			fmt.Println("BGREWRITEAOF command executed successfully.")
		}

		// Закрываем соединение с Redis
		err = clientDBRedis.Close()
		if err != nil {
			fmt.Printf("Error closing Redis connection: %v\n", err)
		} else {
			fmt.Println("Redis connection closed successfully.")
		}
	})
}

func (rep *Repredis) WriteToDB(mutex *sync.Mutex, bufferChanRedis chan *Chunk, errChan chan error, checkAbort chan struct{}) {

	var check *map[string]CheckWrite
	var uid string
	var countPart string
	var cleintid string
	var fileName string = "no.name"
	var start time.Time

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		wg.Wait()
		countParts, err := strconv.Atoi(countPart)
		if err != nil {
			// ... handle errors
			panic(err)
		}
		//Запускаем обновления файла после загрузки данных чтобы скорректировать его размер
		rep.StartChangeLogFileRedis()

		fmt.Println("Горутина записи в REDIS завершилась!")

		_, ok := <-checkAbort
		if !ok {
			fmt.Println("Горутина чтения в REDIS запустилась!")
			go rep.ReadPartFile(start, uid, countParts, cleintid, fileName)
		} else {
			fmt.Println("Горутина чтения REDIS не запустилась!")
		}

		fmt.Println(*check)
	}()

	for msg := range bufferChanRedis {
		if len(msg.Buffer)/1024/1024 > 0 {
			fmt.Println(len(msg.Buffer)/1024/1024, "mb")
		} else {
			fmt.Println(len(msg.Buffer)/1024, "kb")
		}

		check = msg.Check
		checkWrite := *check
		fileName = msg.FileName
		uid = msg.Uidfile

		countPart = msg.CountPart
		uidFull := fmt.Sprintf("%s-%s-%s", msg.Uidfile, msg.Cleintid, msg.NumberPart)
		checkWrite[msg.NumberPart] = CheckWrite{ClientId: msg.Cleintid, Key: uidFull, Part: msg.NumberPart}

		uidFull = GetMD5Hash(uidFull)
		start = msg.Start
		cleintid = msg.Cleintid

		//fmt.Println("Reids-uidFull", uidFull)
		//fmt.Println("Reids-cleintid", msg.Cleintid)
		//fmt.Println("Reids-numberPart", msg.NumberPart)
		wg.Add(1)

		decoratedAcceptPartFile := rep.Retry(10, 1*time.Second)(rep.SavePartFile)

		i, err := strconv.Atoi(msg.NumberPart)
		if err != nil {
			// ... handle error
			panic(err)
		}
		err = decoratedAcceptPartFile(i, wg, msg.Buffer, uidFull, msg.Cleintid, mutex)
		if err != nil {
			//fmt.Println("decoratedAcceptPartFile =>")
			panic(err)
			// обработка ошибки
		}

		//fmt.Println("Записали в базу!", msg.NumberPart)
	}

	//fmt.Println("Вышли из цикла обработки данных из канал!")

	wg.Done()
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (rep *Repredis) ReadPartFile(start time.Time, uid string, countPart int, cleintid string, fileName string) {

	mutexCreate := new(sync.Mutex)
	mutexSend := new(sync.Mutex)
	mutexStop := new(sync.Mutex)

	//Создаем каналы и горутину для сохранения данных в базе в потоке
	rep.repMongo.CreateStreamForSave(cleintid, mutexCreate, fileName, countPart, uid)

	for i := 0; i < countPart; i++ {
		key := fmt.Sprintf("%s-%s-%d", uid, cleintid, i)
		key = GetMD5Hash(key)

		val, err := rep.clientRedis.Get(context.Background(), key).Result()
		if err != nil {
			//fmt.Println("uid:", uid, "countPart:", countPart, "cleintid", cleintid, "PartID:", i)
			//fmt.Println("Пробуем читать ключ:", key)
			rep.repMongo.StopStream("abort", uid, mutexStop)
			panic(err)
		}

		rep.repMongo.SendPartFileToStream([]byte(val), uid, mutexSend)

		//fmt.Printf("Прочитали часть: %d ключ: %s: %d\n", i, key, len([]byte(val)))
	}

	// sgnal: done, abort
	rep.repMongo.StopStream("done", uid, mutexStop)

	// Получаем время, прошедшее с момента start
	elapsed := time.Since(start)

	// Переводим время в секунды с помощью метода Seconds()
	seconds := elapsed.Seconds()

	fmt.Printf("Скрипт выполнялся %.2f секунд\n", seconds)

	err := rep.clientRedis.Close()
	if err != nil {
		fmt.Println("Ошибка закрытия соединения с базой:", err)
	}
	fmt.Printf("Закрыли соединение с базой успешно!")
}
