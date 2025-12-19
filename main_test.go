package main

import (
	"testing"
)

func TestGenerateRandomElements(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantLen int
	}{
		{"Нормальный размер", 1000, 1000},
		{"Нулевой размер", 0, 0},
		{"Отрицательный размер", -10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := generateRandomElements(tt.size)

			if len(data) != tt.wantLen {
				t.Errorf("generateRandomElements(%d) длина = %d, ожидалось %d", tt.size, len(data), tt.wantLen)
			}

			for i, val := range data {
				if val < 0 {
					t.Errorf("Элемент [%d] отрицательный: %d", i, val)
				}
			}
		})
	}
}

func TestMaximum(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"Пустой слайс", []int{}, 0},
		{"Один элемент", []int{42}, 42},
		{"Несколько элементов, максимум в начале", []int{100, 50, 75, 25}, 100},
		{"Несколько элементов, максимум в конце", []int{10, 20, 30, 40, 50}, 50},
		{"Несколько элементов, максимум в середине", []int{1, 3, 5, 2, 4}, 5},
		{"Все одинаковые элементы", []int{7, 7, 7, 7, 7}, 7},
		{"Только положительные числа", []int{1, 5, 3, 9, 2}, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maximum(tt.input)
			if result != tt.expected {
				t.Errorf("maximum(%v) = %d, ожидалось %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaxChunks(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"Пустой слайс", []int{}, 0},
		{"Один элемент", []int{99}, 99},
		{"Меньше элементов, чем чанков", []int{1, 2, 3, 4, 5}, 5},
		{"Ровно на один чанк", []int{10, 20, 30, 40, 50, 60, 70, 80}, 80},
		{"Большой слайс", createTestSlice(100, 999), 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maxChunks(tt.input)
			if result != tt.expected {
				t.Errorf("maxChunks(слайс длиной %d) = %d, ожидалось %d", len(tt.input), result, tt.expected)
			}
		})
	}
}

func createTestSlice(size int, maxVal int) []int {
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = i + 1
	}

	if maxVal > 0 && size > 0 {
		data[size-1] = maxVal
	}

	return data
}

// Все исправил
