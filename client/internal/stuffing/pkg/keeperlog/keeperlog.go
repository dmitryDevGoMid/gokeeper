package keeperlog

import (
	"os"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/checkdir"

	"github.com/sirupsen/logrus"
)

type ContextLogger struct {
	*logrus.Logger
	UserID string
}

type LogField struct {
	Action         string
	Method         string
	NextStepByName string
	RequestByName  string
	UserID         string
}

func ToMap(l LogField) map[string]interface{} {
	return map[string]interface{}{
		"Action":         l.Action,
		"NextStepByName": l.NextStepByName,
		"RequestByName":  l.RequestByName,
		"UserID":         l.UserID,
	}
}

// Создаем обвертку надо логгером и инициализируем поле UserID
func (cl *ContextLogger) WithFields(fields map[string]interface{}) *logrus.Entry {
	fields["UserID"] = cl.UserID
	return cl.Logger.WithFields(fields)
}

func (cl *ContextLogger) SetUserID(userID string) {
	cl.UserID = userID
}

const logFilePath = "gokeeperspace/gokeeperlog/gokeeper.log"

func SetFolder() (*os.File, error) {
	_, pathWithFile, err := checkdir.EnsureDirectoryExistsWithCreateFile("gokeeperspace/gokeeperlog", "gokeeper.log")
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(pathWithFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// CustomHook структура для хука
type CustomHook struct{}

// Levels метод возвращает уровни логирования, которые будут обрабатываться этим хуком
func (hook *CustomHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire метод выполняет определенный код при логировании ошибок
func (hook *CustomHook) Fire(entry *logrus.Entry) error {
	//entryData := entry.Data
	//logrus.Infof("Custom error handling code executed: %v", entryData["Action"])
	return nil
}

func initLogger(log *logrus.Logger) (*os.File, error) {
	file, err := SetFolder()
	if err != nil {
		log.Errorf("Failed to open log file: %v", err)
		return nil, err
	}

	log.SetOutput(file)

	return file, nil
}

func NewLogger() (*logrus.Logger, *os.File, error) {
	customHook := &CustomHook{}
	log := logrus.New()
	log.AddHook(customHook)
	file, err := initLogger(log)
	if err != nil {
		return nil, nil, err
	}
	return log, file, nil
}

func NewContextLogger(userID string, outPutToFile bool) (*ContextLogger, *os.File, error) {
	// Создание экземпляра пользовательского хука
	customHook := &CustomHook{}

	// Создание нового экземпляра логгера logrus
	log := logrus.New()

	// Добавление пользовательского хука к логгеру
	log.AddHook(customHook)
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	var file *os.File
	var err error
	// Инициализация логгера и открытие файла для записи логов
	if outPutToFile {
		file, err = initLogger(log)
		if err != nil {
			// Если произошла ошибка при инициализации логгера, возвращаем nil и ошибку
			return nil, nil, err
		}
	}

	// Возвращаем экземпляр ContextLogger с установленным UserID и файлом для записи логов
	return &ContextLogger{Logger: log, UserID: userID}, file, nil
}
