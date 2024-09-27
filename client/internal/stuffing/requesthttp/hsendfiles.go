package requesthttp

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/calcsum"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
	"github.com/schollz/progressbar/v3"
)

func (h RequestHTTP) GetUnicKey() string {
	// Получаем текущее время в секундах Unix
	timestamp := time.Now().Unix()
	//fmt.Println("Timestamp:", timestamp)

	// Преобразуем время в байты
	timestampString := fmt.Sprint(timestamp)

	// Вычисляем MD5 хэш от байтов + соль в виде id_user
	hash := md5.Sum([]byte(timestampString + h.data.User.IDUser))

	// Преобразуем хэш в строку шестнадцатеричных цифр
	hashString := hex.EncodeToString(hash[:])

	return hashString
}

func (h RequestHTTP) RunSendFiles() tea.Msg {

	err := h.SendFiles(h.GetUnicKey(), h.data.User.IDUser, h.data.OptionFiles.SelectedFile)
	if err != nil {
		return errMsg{err}
	}

	return statusMsg(200)
}

func (h RequestHTTP) SendEnd(keyId string, clientId string, fileName string) error {
	// Создаем клиент go-resty
	client := resty.New()

	//Отправляем
	resp, err := client.R().
		SetHeader("X-Filename", fileName).
		SetHeader("Content-Type", "application/octet-stream").
		SetHeader("Token", h.data.User.Token).
		SetHeader("Status-Send-File", "end").
		SetHeader("Key-ID", keyId).
		SetHeader("Client-ID", clientId).
		SetHeader("File-Name", fileName).
		Post(fmt.Sprintf("%s/files/upload", h.data.Config.Server.AddressFileServer))

	if err != nil {
		return fmt.Errorf("error sending file part: %w", err)
	}

	fmt.Println(resp.StatusCode())

	return nil
}

func (h RequestHTTP) SendFail(filename string, keyid string, clientid string) error {
	// Создаем клиент go-resty
	client := resty.New()

	//Отправляем
	resp, err := client.R().
		SetHeader("X-Filename", filename).
		SetHeader("Content-Type", "application/octet-stream").
		SetHeader("Status-Send-File", "abort").
		SetHeader("Token", h.data.User.Token).
		SetHeader("Key-ID", keyid).
		SetHeader("Client-ID", clientid).
		Post(fmt.Sprintf("%s/files/upload", h.data.Config.Server.AddressFileServer))

	if err != nil {
		return fmt.Errorf("error sending file part: %w", err)
	}

	fmt.Println(resp.StatusCode())

	return nil
}

func (h RequestHTTP) SendFiles(keyid string, clientid string, fileName string) error {

	file, err := os.Open(fileName)
	if err != nil {
		// Обрабатываем ошибку
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Определяем размер файла и размер буфера обмена
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}
	fileSize := fileInfo.Size()
	chunkSize := int64(1024 * 1024) // 1 мегабайт

	checkSum, err := calcsum.СalculateChecksum(fileName, int(chunkSize))
	if err != nil {
		fmt.Println("Error calculating checksum:", err)
		return err
	}

	// Рассчитываем количество частей, на которые нужно разделить файл
	chunkCount := fileSize / chunkSize
	if fileSize%chunkSize != 0 {
		chunkCount++
	}

	// Создаем буфер для чтения файла
	const bufferSize = 1024 * 1024 // 1 MB
	buffer := make([]byte, bufferSize)

	// Создаем буферизованный читатель
	reader := bufio.NewReader(file)

	// Читаем файл частями и отправляем их на сервер
	var part int

	fmt.Printf("SendFiles....")

	bar := progressbar.DefaultBytes(
		fileSize,
		"downloading",
	)

	for {
		select {
		case <-h.data.Cancel:
			h.data.NextStep.NextStepByName = "files"
			err := h.SendFail(fileName, keyid, clientid)
			if err != nil {
				return fmt.Errorf("error send info by FAIL to server: %w", err)
			}
			return errors.New("error processing abort")
		default:
			// Читаем очередную часть файла
			n, err := reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Достигнут конец файла")
					err := h.SendEnd(keyid, clientid, fileInfo.Name())
					if err != nil {
						return fmt.Errorf("error send info by END to server: %w", err)
					}
					fmt.Printf("error send info by END to ser....")

					return nil
				} else {
					return fmt.Errorf("error reading file: %w", err)
				}
			}

			url := fmt.Sprintf("%s/files/upload", h.data.Config.Server.AddressFileServer)

			// Отправляем данные на сервер токен устанавливаем через middleware
			resp, err := h.client.R().
				SetHeader("X-Filename", fileName).
				SetHeader("Content-Type", "application/octet-stream").
				SetHeader("Status-Send-File", "done").
				SetHeader("X-Part", fmt.Sprintf("%d", part)).
				SetHeader("Key-ID", keyid).
				SetHeader("Count-Part", fmt.Sprintf("%d", chunkCount)).
				SetHeader("Client-ID", clientid).
				SetHeader("Check-Sum", checkSum).
				SetHeader("File-Name", fileInfo.Name()).
				SetBody(buffer[:n]).
				Post(url)

			if err != nil {
				// Обрабатываем ошибку
				fmt.Printf("error sending file par....%v", err)
				time.Sleep(10 * time.Second)
				return fmt.Errorf("error sending file part: %w", err)
			}

			if resp.StatusCode() != 200 {
				fmt.Println("exit send file, status code:", resp.StatusCode())
			}

			bar.Add(n)

			// Увеличиваем счетчик частей
			part++
		}
	}
}
