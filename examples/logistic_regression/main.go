package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// WordVector представляет вектор слова.
type WordVector struct {
	Word   string
	Vector []float64
}

// LoadVectors загружает векторы из файла.
func LoadVectors(filename string) (map[string][]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	vectors := make(map[string][]float64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 2 {
			continue
		}
		word := parts[0]
		vector := make([]float64, len(parts)-1)
		for i := 1; i < len(parts); i++ {
			val, err := strconv.ParseFloat(parts[i], 64)
			if err != nil {
				return nil, err
			}
			vector[i-1] = val
		}
		vectors[word] = vector
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return vectors, nil
}

// TextToVector преобразует текст в средний вектор слов.
func TextToVector(text string, vectors map[string][]float64) []float64 {
	words := strings.Fields(text)
	if len(words) == 0 || len(vectors) == 0 {
		return nil // Возвращаем nil, если текст пуст или векторы не загружены
	}

	// Определяем длину вектора на основе первого слова в vectors
	var vectorLength int
	for _, vec := range vectors {
		vectorLength = len(vec)
		break
	}

	sumVector := make([]float64, vectorLength)
	count := 0

	for _, word := range words {
		if vec, ok := vectors[word]; ok {
			for i := range vec {
				sumVector[i] += vec[i]
			}
			count++
		}
	}

	if count > 0 {
		for i := range sumVector {
			sumVector[i] /= float64(count)
		}
	}

	return sumVector
}

// LogisticRegression реализует простую логистическую регрессию.
type LogisticRegression struct {
	Weights []float64
	Bias    float64
	LR      float64 // Learning rate
}

// Predict предсказывает класс для одного примера.
func (lr *LogisticRegression) Predict(features []float64) float64 {
	score := lr.Bias
	for i := range features {
		score += features[i] * lr.Weights[i]
	}
	return 1 / (1 + math.Exp(-score)) // Сигмоида
}

// Train обучает модель на всех данных за несколько эпох.
func (lr *LogisticRegression) Train(data []struct {
	Text  string
	Label float64
}, vectors map[string][]float64, epochs int) {
	for epoch := 0; epoch < epochs; epoch++ {
		for _, example := range data {
			features := TextToVector(example.Text, vectors)
			if features == nil {
				fmt.Printf("Текст '%s' не содержит слов из векторов.\n", example.Text)
				continue
			}
			prediction := lr.Predict(features)
			error := prediction - example.Label

			// Обновляем веса
			for i := range lr.Weights {
				lr.Weights[i] -= lr.LR * error * features[i]
			}
			lr.Bias -= lr.LR * error
		}
	}
}

func main() {
	// Загрузка векторов
	vectors, err := LoadVectors("../../data/vectors.txt.txt")
	if err != nil {
		fmt.Println("Ошибка загрузки векторов:", err)
		return
	}

	if len(vectors) == 0 {
		fmt.Println("Векторы не загружены или файл пуст.")
		return
	}

	// Пример данных для обучения (текст и метка класса)
	trainingData := []struct {
		Text  string
		Label float64
	}{
		{"выродок пидор тупой капитулянт фашист госдеповский усатый грем конченный клоун", 0},
		{"вакцинаторы антивакцинаторы веганы ололо темная бронй шелковых рефлексы подай критик", 1},
	}

	// Инициализация модели
	lr := LogisticRegression{
		Weights: make([]float64, len(vectors["выродок"])), // Используем первое слово для определения длины
		Bias:    0,
		LR:      0.01,
	}

	// Обучение модели (10 эпох)
	lr.Train(trainingData, vectors, 100)

	// Тестирование модели
	testData := []struct {
		Text  string
		Label float64
	}{
		{"тупой пидор", 0},
		{"антивакцинаторы веганы", 1},
		{"билан", 0},
		{"мясоеды", 1},
	}

	for _, data := range testData {
		features := TextToVector(data.Text, vectors)
		if features == nil {
			fmt.Printf("Текст '%s' не содержит слов из векторов.\n", data.Text)
			continue
		}
		prediction := lr.Predict(features)
		fmt.Printf("Текст: %s, Предсказание: %.2f, Ожидалось: %.0f\n", data.Text, prediction, data.Label)
	}
}
