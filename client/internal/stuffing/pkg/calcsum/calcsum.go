package calcsum

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func СalculateChecksum(fileName string, chunkSize int) (string, error) {

	file, err := os.Open(fileName)
	if err != nil {
		// Обрабатываем ошибку
		return "", fmt.Errorf("error opening file: %w", err)
	}

	defer file.Close()

	hash := sha256.New()
	buffer := make([]byte, chunkSize)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return "", err
		}

		if n > 0 {
			hash.Write(buffer[:n]) // Обновляем хеш-объект данными из буфера
		}

		if err == io.EOF {
			break
		}
	}

	return hex.EncodeToString((hash.Sum(nil))), nil
}
