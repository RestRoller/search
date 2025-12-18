package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Константы проекта
const (
	SIZE   = 10_000_000 // Размер основного массива
	CHUNKS = 8          // Количество частей для многопоточного поиска
)

// generateRandomElements генерирует слайс случайных положительных целых чисел
func generateRandomElements(size int) []int {
	// Обработка крайнего случая: некорректный размер
	if size <= 0 {
		return []int{}
	}

	// Создаем новый источник случайных чисел с текущим временем
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// Создание и заполнение слайса
	data := make([]int, size)
	for i := 0; i < size; i++ {
		// Генерируем числа от 0 до size*10 для разнообразия
		data[i] = rng.Intn(size * 10)
	}

	return data
}

// maximum находит максимальное значение в слайсе (однопоточная версия)
func maximum(data []int) int {
	// Обработка крайних случаев
	if len(data) == 0 {
		return 0
	}

	if len(data) == 1 {
		return data[0]
	}

	// Поиск максимума
	maxVal := data[0]
	for i := 1; i < len(data); i++ {
		if data[i] > maxVal {
			maxVal = data[i]
		}
	}

	return maxVal
}

// maxChunks находит максимальное значение многопоточно, разделяя слайс на части
func maxChunks(data []int) int {
	// Обработка крайних случаев
	if len(data) == 0 {
		return 0
	}

	if len(data) == 1 {
		return data[0]
	}

	// Если данных меньше, чем чанков, используем однопоточный подход
	if len(data) < CHUNKS {
		return maximum(data)
	}

	// Вычисляем размер каждого чанка
	chunkSize := len(data) / CHUNKS

	// Слайс для хранения максимумов каждого чанка
	chunkMaxes := make([]int, CHUNKS)

	// WaitGroup для синхронизации горутин
	var wg sync.WaitGroup
	wg.Add(CHUNKS)

	// Мьютекс для безопасного доступа к результатам (альтернативный подход)
	var mu sync.Mutex

	// Запускаем горутины для поиска максимумов в каждом чанке
	for i := 0; i < CHUNKS; i++ {
		// Вычисляем границы чанка
		start := i * chunkSize
		end := start + chunkSize

		// Для последнего чанка берем все оставшиеся элементы
		if i == CHUNKS-1 {
			end = len(data)
		}

		// Запускаем горутину для обработки чанка
		go func(chunkIndex int, chunk []int) {
			defer wg.Done()

			// Находим максимум в чанке
			chunkMax := chunk[0]
			for _, value := range chunk[1:] {
				if value > chunkMax {
					chunkMax = value
				}
			}

			// Сохраняем результат с использованием мьютекса
			mu.Lock()
			chunkMaxes[chunkIndex] = chunkMax
			mu.Unlock()
		}(i, data[start:end])
	}

	// Ожидаем завершения всех горутин
	wg.Wait()

	// Находим максимальное значение среди всех чанков
	return maximum(chunkMaxes)
}

func main() {
	fmt.Println("=== Поиск максимального значения в большом массиве данных ===")
	fmt.Printf("Размер массива: %d элементов\n", SIZE)
	fmt.Printf("Количество частей для многопоточного поиска: %d\n\n", CHUNKS)

	// Генерация данных
	fmt.Println("Генерация данных...")
	data := generateRandomElements(SIZE)
	fmt.Printf("Данные сгенерированы. Длина слайса: %d\n\n", len(data))

	// Проверка, что данные не пустые
	if len(data) == 0 {
		fmt.Println("Ошибка: сгенерирован пустой массив данных")
		return
	}

	// Однопоточный поиск максимума
	fmt.Println("Запуск однопоточного поиска...")
	startTime := time.Now()
	singleMax := maximum(data)
	singleDuration := time.Since(startTime)

	fmt.Printf("Однопоточный результат: %d\n", singleMax)
	fmt.Printf("Время выполнения: %d микросекунд\n\n", singleDuration.Microseconds())

	// Многопоточный поиск максимума
	fmt.Println("Запуск многопоточного поиска...")
	startTime = time.Now()
	multiMax := maxChunks(data)
	multiDuration := time.Since(startTime)

	fmt.Printf("Многопоточный результат: %d\n", multiMax)
	fmt.Printf("Время выполнения: %d микросекунд\n\n", multiDuration.Microseconds())

	// Сравнение результатов и производительности
	fmt.Println("=== Сравнение результатов ===")
	if singleMax == multiMax {
		fmt.Println("✓ Результаты совпадают!")
	} else {
		fmt.Printf("⚠ Результаты различаются: однопоточный=%d, многопоточный=%d\n",
			singleMax, multiMax)
	}

	// Расчет ускорения
	if multiDuration > 0 {
		speedup := float64(singleDuration.Microseconds()) / float64(multiDuration.Microseconds())
		fmt.Printf("Ускорение: %.2f раз\n", speedup)

		if speedup > 1 {
			fmt.Printf("Многопоточная версия быстрее на %.1f%%\n",
				(speedup-1)*100)
		} else if speedup < 1 {
			fmt.Printf("Однопоточная версия быстрее на %.1f%%\n",
				(1/speedup-1)*100)
		} else {
			fmt.Println("Производительность одинаковая")
		}
	}

	// Дополнительная статистика
	fmt.Println("\n=== Дополнительная информация ===")
	fmt.Printf("Размер каждого чанка: %d элементов\n", len(data)/CHUNKS)
	fmt.Printf("Использовано горутин: %d\n", CHUNKS)

	// Демонстрация работы с маленькими массивами
	demoSmallArrays()
}

// demoSmallArrays демонстрирует работу с небольшими массивами
func demoSmallArrays() {
	fmt.Println("\n=== Демонстрация с небольшими массивами ===")

	testCases := []struct {
		name string
		size int
	}{
		{"Пустой массив", 0},
		{"Один элемент", 1},
		{"Меньше чанков", 5},
		{"Ровно чанков", 8},
		{"Немного больше чанков", 12},
	}

	for _, tc := range testCases {
		data := generateRandomElements(tc.size)
		singleResult := maximum(data)
		multiResult := maxChunks(data)

		status := "✓"
		if singleResult != multiResult {
			status = "⚠"
		}

		fmt.Printf("%s %s (размер=%d): single=%d, multi=%d\n",
			status, tc.name, len(data), singleResult, multiResult)
	}
}
