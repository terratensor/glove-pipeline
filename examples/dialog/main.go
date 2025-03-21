package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
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

// CosineSimilarity вычисляет косинусное сходство между двумя векторами.
func CosineSimilarity(vec1, vec2 []float64) (float64, error) {
	if len(vec1) != len(vec2) {
		return 0, fmt.Errorf("векторы должны быть одинаковой длины")
	}

	var dotProduct, magnitude1, magnitude2 float64
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
		magnitude1 += vec1[i] * vec1[i]
		magnitude2 += vec2[i] * vec2[i]
	}

	magnitude1 = math.Sqrt(magnitude1)
	magnitude2 = math.Sqrt(magnitude2)

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0, fmt.Errorf("один из векторов имеет нулевую длину")
	}

	return dotProduct / (magnitude1 * magnitude2), nil
}

// FindSynonyms находит синонимы для заданного слова.
func FindSynonyms(word string, vectors map[string][]float64, topN int) ([]string, error) {
	targetVector, ok := vectors[word]
	if !ok {
		return nil, fmt.Errorf("слово '%s' не найдено в векторах", word)
	}

	type Similarity struct {
		Word       string
		Similarity float64
	}

	var similarities []Similarity
	for w, vec := range vectors {
		if w == word {
			continue
		}
		sim, err := CosineSimilarity(targetVector, vec)
		if err != nil {
			return nil, err
		}
		similarities = append(similarities, Similarity{Word: w, Similarity: sim})
	}

	// Сортировка по убыванию сходства
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	// Выбор topN синонимов
	var synonyms []string
	for i := 0; i < topN && i < len(similarities); i++ {
		synonyms = append(synonyms, similarities[i].Word)
	}

	return synonyms, nil
}

func main() {
	// Загрузка векторов
	vectors, err := LoadVectors("../../data/fct/vectors.txt.txt")
	if err != nil {
		fmt.Println("Ошибка загрузки векторов:", err)
		return
	}

	if len(vectors) == 0 {
		fmt.Println("Векторы не загружены или файл пуст.")
		return
	}

	// Канал для обработки сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Канал для завершения работы
	done := make(chan struct{})

	// Запуск горутины для обработки ввода
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Введите слово (или 'exit' для выхода): ")
			if !scanner.Scan() {
				break
			}
			word := scanner.Text()
			if word == "exit" {
				close(done)
				return
			}
			word = strings.ToLower(word)

			topN := 30
			synonyms, err := FindSynonyms(word, vectors, topN)
			if err != nil {
				fmt.Printf("Ошибка поиска синонимов для слова '%s': %v\n", word, err)
				continue
			}

			fmt.Printf("Топ-%d синонимов для слова '%s':\n", topN, word)
			for i, syn := range synonyms {
				fmt.Printf("%d. %s\n", i+1, syn)
			}
			fmt.Println()
		}
	}()

	// Ожидание сигнала или завершения работы
	select {
	case <-sigChan:
		fmt.Println("\nПолучен сигнал завершения. Завершение работы...")
	case <-done:
		fmt.Println("Завершение работы...")
	}
}
