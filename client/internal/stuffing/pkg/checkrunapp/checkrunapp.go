package checkrunapp

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dmitryDevGoMid/gokeeper/client/internal/stuffing/pkg/checkdir"
)

func findFreePort(startPort int) (int, error) {
	for port := startPort; port < 65536; port++ {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("не удалось найти свободный порт")
}

// Мьютекс для синхронизации доступа к файлу
var mu sync.Mutex

func runWriteChangeFile(ctx context.Context, hiddenFile string, port int) {
	// Горутина для перезаписи файла каждые 30 секунд
	go func() {

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mu.Lock()
				err := createHiddenFile(hiddenFile, port)
				mu.Unlock()
				if err != nil {
					fmt.Printf("Ошибка при создании файла: %v\n", err)
				}
			}
		}
	}()
}

func getPath() (string, error) {
	path, _, _, err := checkdir.EnsureDirectoryExists("gokeeperspace/gokeeperconfig", ".hiddenfile.txt")
	if err != nil {
		return "", err
	}
	hiddenFile := filepath.Join(path, "/.hiddenfile.txt")

	return hiddenFile, nil
}

func checkFileAndCreateFile(hiddenFile string) (bool, error) {
	exitApp := false
	if _, err := os.Stat(hiddenFile); os.IsNotExist(err) {
		startPort := 8000 // Начальный порт для перебора
		freePort, err := findFreePort(startPort)
		if err != nil {
			log.Fatalf("Not found free port: %v", err)
		}
		// Файл не существует, создаем его и записываем текущее время и номер порта
		mu.Lock()
		if !isPortInUse(freePort) {
			err = createHiddenFile(hiddenFile, freePort)
			if err != nil {
				return exitApp, err
			}
			exitApp = true
			fmt.Printf("Заняли порт %d \n", freePort)
		}

		mu.Unlock()

	}
	return exitApp, nil
}

func checkFileAndReadFile(hiddenFile string) (bool, int, error) {
	exitApp := false
	usePort := 0
	startPort := 8000 // Начальный порт для перебора
	freePort, err := findFreePort(startPort)
	if err != nil {
		log.Fatalf("Not found free port: %v", err)
	}
	// Файл существует, читаем его содержимое
	port, err := readHiddenFilePort(hiddenFile)

	if err != nil {
		log.Fatal("Не удалось прочитать скрытый файл:", err)
		return true, usePort, err
	}

	lastTime, err := parseTimeFromContent(hiddenFile)
	if err != nil {
		log.Fatal("Не удалось распарсить время из файла:", err)
		return true, usePort, err
	}

	elapsedTime := time.Since(lastTime).Seconds()

	if elapsedTime > 10 && !isPortInUse(port) {

		fmt.Printf("Заняли %d порт.\n", port)
		err := createHiddenFile(hiddenFile, port)
		if err != nil {
			log.Fatal("Не удалось прочитать скрытый файл:", err)
			return true, usePort, err
		}
		usePort = port
	} else if elapsedTime > 30 && isPortInUse(port) {

		isPortInUse(freePort)
		fmt.Printf("Заняли %d порт.\n", freePort)
		err := createHiddenFile(hiddenFile, freePort)
		if err != nil {
			log.Fatal("Не удалось прочитать скрытый файл:", err)
			return true, usePort, err
		}
		usePort = freePort

	} else {
		exitApp = true
	}

	return exitApp, usePort, nil
}

func StartCheckRunApp(ctx context.Context) (bool, error) {
	//Выйти из приложения, по умолчанию нет
	exitApp := false
	hiddenFile, err := getPath()
	if err != nil {
		return exitApp, err
	}
	port := -1

	exitApp, err = checkFileAndCreateFile(hiddenFile)
	if err != nil {
		return exitApp, err
	}

	if !exitApp {
		exitApp, port, err = checkFileAndReadFile(hiddenFile)
		if err != nil {
			return exitApp, err
		}
	}

	if !exitApp && port > 0 {
		runWriteChangeFile(ctx, hiddenFile, port)
	}

	return exitApp, nil
}

func createHiddenFile(filePath string, port int) error {
	currentTime := time.Now().Format(time.RFC3339) // Используем формат RFC3339 для точности
	content := fmt.Sprintf("%s - %d", currentTime, port)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Не удалось создать скрытый файл: %v\n", err)
		return err
	}

	return nil
}

func readHiddenFileTime(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	parts := strings.Split(string(content), " - ")
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат содержимого файла")
	}
	timeString := parts[0]
	return timeString, nil
}

func readHiddenFilePort(filePath string) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return -1, err
	}
	parts := strings.Split(string(content), " - ")
	if len(parts) != 2 {
		return -1, fmt.Errorf("неверный формат содержимого файла")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1, err
	}
	return port, nil
}

func isPortInUse(port int) bool {
	_, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	fmt.Println("Пробуем занять порт:", port)
	if err != nil {
		return true
	}
	//ln.Close()
	return false
}

func parseTimeFromContent(path string) (time.Time, error) {
	content, err := readHiddenFileTime(path)
	if err != nil {
		return time.Time{}, err
	}

	lastTime, err := time.Parse(time.RFC3339, content)
	if err != nil {
		return time.Time{}, err
	}

	return lastTime, nil
}
