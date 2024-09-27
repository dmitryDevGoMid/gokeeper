package checkdir

import (
	"fmt"
	"os"
	"path/filepath"
)

// ensureDirectoryExists проверяет наличие папки в домашней директории, создает её, если она не существует,
// и возвращает полный путь к этой папке.
func EnsureDirectoryExists(dirName, fileName string) (string, string, bool, error) {
	var isNotExistFile bool = false
	// Получаем путь к домашней директории
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", isNotExistFile, fmt.Errorf("failed to get home directory: %v", err)
	}

	// Формируем полный путь к папке
	fullPath := filepath.Join(homeDir, dirName)

	// Проверяем, существует ли папка
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Создаем папку, если она не существует
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return "", "", isNotExistFile, fmt.Errorf("failed to create directory: %v", err)
		}
	} else if err != nil {
		return "", "", isNotExistFile, fmt.Errorf("failed to check directory: %v", err)
	}

	// Формируем полный путь к файлу
	filePath := filepath.Join(fullPath, fileName)

	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Создаем файл, если он не существует
		isNotExistFile = true
	} else if err != nil {
		return "", "", isNotExistFile, fmt.Errorf("failed to check file: %v", err)
	}

	// Возвращаем полный путь к папке и файлу
	return fullPath, filePath, isNotExistFile, nil
}

// ensureDirectoryExists проверяет наличие папки в домашней директории, создает её, если она не существует,
// и возвращает полный путь к этой папке + создаем файл если он не существует
func EnsureDirectoryExistsWithCreateFile(dirName, fileName string) (string, string, error) {
	// Получаем путь к домашней директории
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to get home directory: %v", err)
	}

	// Формируем полный путь к папке
	fullPath := filepath.Join(homeDir, dirName)

	// Проверяем, существует ли папка
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Создаем папку, если она не существует
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return "", "", fmt.Errorf("failed to create directory: %v", err)
		}
	} else if err != nil {
		return "", "", fmt.Errorf("failed to check directory: %v", err)
	}

	// Формируем полный путь к файлу
	filePath := filepath.Join(fullPath, fileName)

	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Создаем файл, если он не существует
		file, err := os.Create(filePath)
		if err != nil {
			return "", "", fmt.Errorf("failed to create file: %v", err)
		}
		file.Close()
	} else if err != nil {
		return "", "", fmt.Errorf("failed to check file: %v", err)
	}

	// Возвращаем полный путь к папке и файлу
	return fullPath, filePath, nil
}
