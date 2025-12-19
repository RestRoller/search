package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SIZE   = 10_000_000
	CHUNKS = 8
)

func generateRandomElements(size int) []int {
	if size <= 0 {
		return []int{}
	}

	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rng.Intn(1000000)
	}

	return data
}

func maximum(data []int) int {
	if len(data) == 0 {
		return 0
	}

	if len(data) == 1 {
		return data[0]
	}

	maxVal := data[0]
	for i := 1; i < len(data); i++ {
		if data[i] > maxVal {
			maxVal = data[i]
		}
	}

	return maxVal
}

func maxChunks(data []int) int {
	if len(data) == 0 {
		return 0
	}

	if len(data) == 1 {
		return data[0]
	}

	if len(data) < CHUNKS {
		return maximum(data)
	}

	chunkSize := len(data) / CHUNKS
	results := make([]int, CHUNKS)

	var wg sync.WaitGroup
	wg.Add(CHUNKS)

	for i := 0; i < CHUNKS; i++ {
		start := i * chunkSize
		end := start + chunkSize

		if i == CHUNKS-1 {
			end = len(data)
		}

		go func(idx int, chunk []int) {
			defer wg.Done()
			results[idx] = maximum(chunk)
		}(i, data[start:end])
	}

	wg.Wait()
	return maximum(results)
}

func main() {
	fmt.Println("=== Поиск максимального значения в большом массиве данных ===")
	fmt.Printf("Размер массива: %d элементов\n", SIZE)
	fmt.Printf("Количество частей для многопоточного поиска: %d\n\n", CHUNKS)

	fmt.Println("Генерация данных...")
	data := generateRandomElements(SIZE)
	fmt.Printf("Данные сгенерированы. Длина слайса: %d\n\n", len(data))

	if len(data) == 0 {
		fmt.Println("Ошибка: сгенерирован пустой массив данных")
		return
	}

	fmt.Println("Запуск однопоточного поиска...")
	startTime := time.Now()
	singleMax := maximum(data)
	singleDuration := time.Since(startTime)

	fmt.Printf("Однопоточный результат: %d\n", singleMax)
	fmt.Printf("Время выполнения: %d микросекунд\n\n", singleDuration.Microseconds())

	fmt.Println("Запуск многопоточного поиска...")
	startTime = time.Now()
	multiMax := maxChunks(data)
	multiDuration := time.Since(startTime)

	fmt.Printf("Многопоточный результат: %d\n", multiMax)
	fmt.Printf("Время выполнения: %d микросекунд\n\n", multiDuration.Microseconds())

	fmt.Println("=== Сравнение результатов ===")
	if singleMax == multiMax {
		fmt.Println("✓ Результаты совпадают!")
	} else {
		fmt.Printf("⚠ Результаты различаются: однопоточный=%d, многопоточный=%d\n", singleMax, multiMax)
	}

	if multiDuration.Microseconds() > 0 {
		speedup := float64(singleDuration.Microseconds()) / float64(multiDuration.Microseconds())
		fmt.Printf("Ускорение: %.2f раз\n", speedup)

		if speedup > 1 {
			fmt.Printf("Многопоточная версия быстрее на %.1f%%\n", (speedup-1)*100)
		}

		if speedup < 1 {
			fmt.Printf("Однопоточная версия быстрее на %.1f%%\n", (1/speedup-1)*100)
		}

		if speedup > 0.999 && speedup < 1.001 {
			fmt.Println("Производительность одинаковая")
		}
	}

	fmt.Println("\n=== Дополнительная информация ===")
	fmt.Printf("Размер каждого чанка: %d элементов\n", len(data)/CHUNKS)
	fmt.Printf("Использовано горутин: %d\n", CHUNKS)
}
