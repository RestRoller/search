package main

import (
	"math/rand"
	"testing"
	"time"
)

// TestGenerateRandomElements тестирует функцию generateRandomElements
func TestGenerateRandomElements(t *testing.T) {
	t.Run("Нормальный размер", func(t *testing.T) {
		size := 1000
		data := generateRandomElements(size)

		if len(data) != size {
			t.Errorf("Ожидалась длина %d, получено %d", size, len(data))
		}

		// Проверяем, что все элементы неотрицательные
		for i, val := range data {
			if val < 0 {
				t.Errorf("Элемент [%d] отрицательный: %d", i, val)
			}
		}
	})

	t.Run("Нулевой размер", func(t *testing.T) {
		data := generateRandomElements(0)

		if data == nil {
			t.Error("Для размера 0 не должен возвращаться nil")
		}

		if len(data) != 0 {
			t.Errorf("Для размера 0 ожидался пустой слайс, получено %d элементов", len(data))
		}
	})

	t.Run("Отрицательный размер", func(t *testing.T) {
		data := generateRandomElements(-10)

		if data == nil {
			t.Error("Для отрицательного размера не должен возвращаться nil")
		}

		if len(data) != 0 {
			t.Errorf("Для отрицательного размера ожидался пустой слайс, получено %d элементов", len(data))
		}
	})

	t.Run("Разные результаты", func(t *testing.T) {
		// Два вызова должны генерировать разные последовательности
		data1 := generateRandomElements(100)
		// Небольшая задержка для гарантии разного seed
		time.Sleep(1 * time.Millisecond)
		data2 := generateRandomElements(100)

		// Проверяем, что есть различия
		allEqual := true
		for i := 0; i < 100; i++ {
			if data1[i] != data2[i] {
				allEqual = false
				break
			}
		}

		if allEqual {
			t.Error("Два вызова функции должны генерировать разные данные")
		}
	})
}

// TestMaximum тестирует функцию maximum
func TestMaximum(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{
			name:     "Пустой слайс",
			input:    []int{},
			expected: 0,
		},
		{
			name:     "Один элемент",
			input:    []int{42},
			expected: 42,
		},
		{
			name:     "Несколько элементов, максимум в начале",
			input:    []int{100, 50, 75, 25},
			expected: 100,
		},
		{
			name:     "Несколько элементов, максимум в конце",
			input:    []int{10, 20, 30, 40, 50},
			expected: 50,
		},
		{
			name:     "Несколько элементов, максимум в середине",
			input:    []int{1, 3, 5, 2, 4},
			expected: 5,
		},
		{
			name:     "Все одинаковые элементы",
			input:    []int{7, 7, 7, 7, 7},
			expected: 7,
		},
		{
			name:     "Отрицательные числа",
			input:    []int{-10, -5, -20, -1},
			expected: -1,
		},
		{
			name:     "Смешанные положительные и отрицательные",
			input:    []int{-10, 0, 10, -5, 5},
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maximum(tt.input)
			if result != tt.expected {
				t.Errorf("maximum(%v) = %d, ожидалось %d",
					tt.input, result, tt.expected)
			}
		})
	}
}

// TestMaxChunks тестирует функцию maxChunks
func TestMaxChunks(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{
			name:     "Пустой слайс",
			input:    []int{},
			expected: 0,
		},
		{
			name:     "Один элемент",
			input:    []int{99},
			expected: 99,
		},
		{
			name:     "Меньше элементов, чем чанков",
			input:    []int{1, 2, 3, 4, 5},
			expected: 5,
		},
		{
			name:     "Ровно на один чанк",
			input:    []int{10, 20, 30, 40, 50, 60, 70, 80},
			expected: 80,
		},
		{
			name:     "Большой слайс, максимум в первом чанке",
			input:    createTestSlice(100, 0, 999),
			expected: 999,
		},
		{
			name:     "Большой слайс, максимум в последнем чанке",
			input:    createTestSlice(100, 0, 0, 1000),
			expected: 1000,
		},
		{
			name:     "Большой слайс, максимум в случайном чанке",
			input:    createTestSlice(100, 5000),
			expected: 5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maxChunks(tt.input)
			if result != tt.expected {
				t.Errorf("maxChunks(слайс длиной %d) = %d, ожидалось %d",
					len(tt.input), result, tt.expected)
			}
		})
	}
}

// TestConsistency тестирует согласованность результатов single и multi
func TestConsistency(t *testing.T) {
	// Создаем детерминированный источник для воспроизводимых тестов
	src := rand.NewSource(42)
	rng := rand.New(src)

	// Генерируем тестовые данные
	testData := []struct {
		name string
		data []int
	}{
		{"Маленький слайс", generateDeterministicSlice(rng, 10)},
		{"Средний слайс", generateDeterministicSlice(rng, 1000)},
		{"Большой слайс", generateDeterministicSlice(rng, 10000)},
		{"Очень большой слайс", generateDeterministicSlice(rng, 100000)},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			singleResult := maximum(td.data)
			multiResult := maxChunks(td.data)

			if singleResult != multiResult {
				t.Errorf("Несогласованность результатов для %s: single=%d, multi=%d",
					td.name, singleResult, multiResult)
			}
		})
	}
}

// BenchmarkSingleMax бенчмарк для однопоточного поиска
func BenchmarkSingleMax(b *testing.B) {
	// Используем детерминированный источник для воспроизводимости
	src := rand.NewSource(42)
	rng := rand.New(src)
	data := make([]int, 1_000_000)
	for i := 0; i < len(data); i++ {
		data[i] = rng.Intn(10_000_000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maximum(data)
	}
}

// BenchmarkMultiMax бенчмарк для многопоточного поиска
func BenchmarkMultiMax(b *testing.B) {
	// Используем детерминированный источник для воспроизводимости
	src := rand.NewSource(42)
	rng := rand.New(src)
	data := make([]int, 1_000_000)
	for i := 0; i < len(data); i++ {
		data[i] = rng.Intn(10_000_000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maxChunks(data)
	}
}

// Helper функция для создания тестовых слайсов
func createTestSlice(size int, values ...int) []int {
	if len(values) == 0 {
		// Создаем слайс со значениями от 0 до size-1
		data := make([]int, size)
		for i := 0; i < size; i++ {
			data[i] = i
		}
		return data
	}

	// Создаем слайс и размещаем специальные значения
	data := make([]int, size)

	// Заполняем базовыми значениями
	for i := 0; i < size; i++ {
		data[i] = i % 100
	}

	// Размещаем специальные значения
	for i, val := range values {
		if i < size {
			data[i*size/len(values)] = val
		}
	}

	return data
}

// generateDeterministicSlice создает детерминированный слайс
func generateDeterministicSlice(rng *rand.Rand, size int) []int {
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rng.Intn(size * 10)
	}
	return data
}
