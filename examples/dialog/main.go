package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
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

// LoadStopWords загружает стоп-слова из файла.
func LoadStopWords(filename string) (map[string]struct{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stopWords := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			stopWords[word] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stopWords, nil
}

// RemovePunctuation удаляет знаки пунктуации из текста.
func RemovePunctuation(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, text)
}

// RemoveStopWords удаляет стоп-слова из текста.
func RemoveStopWords(text string, stopWords map[string]struct{}) string {
	words := strings.Fields(text)
	var filteredWords []string
	for _, word := range words {
		if _, isStopWord := stopWords[word]; !isStopWord {
			filteredWords = append(filteredWords, word)
		}
	}
	return strings.Join(filteredWords, " ")
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

// FindSynonyms находит синонимы для заданного вектора.
func FindSynonyms(targetVector []float64, vectors map[string][]float64, topN int) ([]string, error) {
	type Similarity struct {
		Word       string
		Similarity float64
	}

	var similarities []Similarity
	for word, vec := range vectors {
		sim, err := CosineSimilarity(targetVector, vec)
		if err != nil {
			return nil, err
		}
		similarities = append(similarities, Similarity{Word: word, Similarity: sim})
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

	// Загрузка стоп-слов
	stopWords, err := LoadStopWords("../../data/fct/stopwords.txt")
	if err != nil {
		fmt.Println("Ошибка загрузки стоп-слов:", err)
		return
	}

	// Проверка аргументов командной строки
	if len(os.Args) < 2 {
		fmt.Println("Использование: программа <режим>")
		fmt.Println("Режимы: word (поиск по слову), phrase (поиск по фразе)")
		return
	}

	mode := os.Args[1]
	if mode != "word" && mode != "phrase" {
		fmt.Println("Неправильный режим. Используйте 'word' или 'phrase'.")
		return
	}

	// Бесконечный цикл для интерактивного диалога
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if mode == "word" {
			fmt.Println("Введите слово для поиска синонимов (или 'exit' для выхода):")
		} else {
			fmt.Println("Введите фразу для поиска синонимов (или 'exit' для выхода):")
		}

		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		var targetVector []float64
		if mode == "word" {
			// Поиск синонимов по слову
			vec, ok := vectors[input]
			if !ok {
				fmt.Printf("Слово '%s' не найдено в векторах.\n", input)
				continue
			}
			targetVector = vec
		} else {
			// Поиск синонимов по фразе
			// Приводим фразу к нижнему регистру
			input = strings.ToLower(input)
			// Удаляем знаки пунктуации
			input = RemovePunctuation(input)
			// Удаляем стоп-слова
			filteredPhrase := RemoveStopWords(input, stopWords)
			// Выводим итоговую фразу
			fmt.Printf("Очищенная фраза: '%s'\n", filteredPhrase)
			// Преобразуем фразу в вектор
			targetVector = TextToVector(filteredPhrase, vectors)
			if targetVector == nil {
				fmt.Println("Фраза не содержит слов из векторов.")
				continue
			}
		}

		// Поиск синонимов
		topN := 20
		synonyms, err := FindSynonyms(targetVector, vectors, topN)
		if err != nil {
			fmt.Println("Ошибка поиска синонимов:", err)
			continue
		}

		fmt.Printf("Топ-%d синонимов:\n", topN)
		for i, syn := range synonyms {
			fmt.Printf("%d. %s\n", i+1, syn)
		}
		fmt.Println()
	}
}
