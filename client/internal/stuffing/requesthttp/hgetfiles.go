package requesthttp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/keeperlog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
	"github.com/schollz/progressbar/v3"
)

func (h RequestHTTP) RunGetFiles() tea.Msg {
	h.logField.Method = "RunGetFiles"
	//Формируем контекс для
	ctx := context.Background()

	//Выполняем запрос на сервер
	err := h.GetFilesFromServer(ctx, h.GetUnicKey(), h.data.User.IDUser)
	if err != nil {
		h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Error("Error:", err.Error())
		return errMsg{err}
	}

	return statusMsg(200)
}

func (h RequestHTTP) RequestToServer(keyid string, clientid string, selectedFile model.OptionFiles) (*resty.Response, error) {
	h.logField.Method = "RequestToServer"

	// Выполняем HTTP POST-запрос для загрузки файла
	resp, err := h.client.R().
		// Отменяем парсинг ответа на стороне фреймворка, чтобы обработать его вручную
		SetDoNotParseResponse(true).
		// Устанавливаем заголовок Accept для указания типа контента, который мы ожидаем получить
		SetHeader("Accept", "application/octet-stream").
		// Устанавливаем заголовок Client-ID для идентификации клиента
		SetHeader("Client-ID", clientid).
		// Устанавливаем заголовок File-ID для идентификации файла, который мы хотим загрузить
		SetHeader("File-ID", selectedFile.ID).
		// Устанавливаем заголовок File-Size для указания размера файла
		SetHeader("File-Size", fmt.Sprintf("%d", selectedFile.Length)).
		// Выполняем POST-запрос на указанный URL для загрузки файла
		//Post("http://localhost:3000/files/download")
		Post(fmt.Sprintf("%s/files/download", h.data.Config.Server.AddressFileServer))

	if err != nil {
		h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Error("Error:", err.Error())
		return nil, err
	}
	return resp, nil

}

func (h RequestHTTP) GetFilesFromServer(ctx context.Context, keyid string, clientid string) error {

	request := h.data.RequestHTTP["fileslist"]

	chocedKey, err := strconv.Atoi(request.Choced)
	if err != nil {
		return err
	}
	filesList := *request.FilesList
	selectedFile := filesList[chocedKey]

	//Запрос файла с сервера
	resp, err := h.RequestToServer(keyid, clientid, selectedFile)
	if err != nil {
		return err
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		// Обновление токена
		if err := h.RefreshToken(); err != nil {
			return err
		}

		//Выполняем повторный запрос файла с сервера после обновления токена
		resp, err = h.RequestToServer(keyid, clientid, selectedFile)
		if err != nil {
			return err
		}
	}

	//outputFileName := fmt.Sprintf("%s/%s", outputDirectory, selectedFile.Filename)
	outputFileName := request.OutputFileName
	if outputFileName == "" {
		return errors.New("OutputFile is empty")
	}

	// Получение размера файла из заголовков ответа
	contentLength := resp.Header().Get("Content-Length")
	if contentLength == "" {
		return errors.New("Content-Length header is missing")
	}

	// Создание файла для записи
	file, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Канал для сигнала об успешной загрузки
	done := make(chan struct{})

	// Создание прогресс-бара
	bar := progressbar.DefaultBytes(
		selectedFile.Length,
		"downloading",
	)

	// Создаем переменную wg типа sync.WaitGroup для синхронизации горутин
	var wg sync.WaitGroup

	// Увеличиваем счетчик группы ожидания на 1, чтобы указать, что мы добавляем одну горутину
	wg.Add(1)

	// Запускаем анонимную функцию в отдельной горутине
	go func() {
		// Отложенный вызов wg.Done() для уменьшения счетчика группы ожидания на 1 при завершении горутины
		defer wg.Done()

		// Отложенный вызов close(done) для закрытия канала done при завершении горутины
		defer close(done)

		// Копируем данные из resp.RawBody() в io.MultiWriter(file, bar)
		// io.MultiWriter создает io.Writer, который записывает данные в несколько io.Writer одновременно (в данном случае, в file и bar)
		_, err = io.Copy(io.MultiWriter(file, bar), resp.RawBody())

		// Проверяем, возникла ли ошибка при копировании данных
		if err != nil {
			// Завершаем выполнение горутины
			h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Error("Error:", err.Error())
			return
		}

		// Выводим сообщение об успешном завершении загрузки файла
		h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info("Gorutine - success donwload file:")

	}()

	// Ожидание завершения загрузки или отмены
	select {
	case <-h.data.Cancel:
		err := resp.RawBody().Close()
		if err != nil {
			return err
		}

		file.Close()

		os.Remove(outputFileName) // Удаляем частично загруженный файл
		h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info("Remove file:", outputFileName)

		// Отправляем запрос на сервер для отмены загрузки
		resp, err = h.client.R().
			SetHeader("File-ID", selectedFile.ID).
			Post(fmt.Sprintf("%s/files/download/cancel", h.data.Config.Server.AddressFileServer))

		if err != nil {
			return err
		} else {
			h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info("Download canceled successfully")
		}
	case <-done:
		h.data.Log.WithFields(keeperlog.ToMap(h.logField)).Info("Download file done")
	}

	wg.Wait()

	return nil
}
