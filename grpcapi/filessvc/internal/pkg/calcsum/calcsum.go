package calcsum

import (
	"crypto/sha256"
	"encoding/hex"
)

// const chunkSize = 1024 * 1024 // Размер чанка 1 МБ
// dataChan - канал с бинарными данными
// done - канал-сигнал о завершении
// getSum - канал принимает расчитанную сумму
func CalcSumByGorutine(dataChan <-chan []byte, getSum chan<- string, done <-chan struct{}) {
	hash := sha256.New()
	for {
		select {
		case data := <-dataChan:
			hash.Write(data)
		case <-done:
			getSum <- hex.EncodeToString(hash.Sum(nil))
			return
		default:
		}
	}
}

func RunCalcSum() (chan<- []byte, <-chan string, chan<- struct{}) {
	dataChan := make(chan []byte)
	getSum := make(chan string)
	done := make(chan struct{})

	go CalcSumByGorutine(dataChan, getSum, done)

	return dataChan, getSum, done
}
