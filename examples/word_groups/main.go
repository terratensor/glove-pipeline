package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Jeffail/tunny"
	"github.com/cheggaaa/pb/v3"
)

type WordVector struct {
	Word   string
	Vector []float64
}

var (
	vectors     map[string][]float64
	vectorsLock sync.RWMutex
	progressBar *pb.ProgressBar
)

func loadGloveVectors(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	tempVectors := make(map[string][]float64)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 2 {
			continue
		}

		word := parts[0]
		vector := make([]float64, len(parts)-1)

		for i := 1; i < len(parts); i++ {
			_, err := fmt.Sscanf(parts[i], "%f", &vector[i-1])
			if err != nil {
				return err
			}
		}
		tempVectors[word] = vector
	}

	vectorsLock.Lock()
	vectors = tempVectors
	vectorsLock.Unlock()

	return scanner.Err()
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func cleanText(text string) []string {
	reg := regexp.MustCompile(`[^а-яА-ЯёЁa-zA-Z\s]`)
	cleaned := reg.ReplaceAllString(text, " ")
	return strings.Fields(strings.ToLower(cleaned))
}

var stopWords = map[string]bool{
	"и": true, "в": true, "не": true, "по": true, "же": true, "с": true, "о": true,
}

func filterStopWords(words []string) []string {
	filtered := make([]string, 0, len(words))
	for _, word := range words {
		if !stopWords[word] {
			filtered = append(filtered, word)
		}
	}
	return filtered
}

func cosineSimilarity(a, b []float64) float64 {
	dotProduct := 0.0
	magA := 0.0
	magB := 0.0

	for i := range a {
		dotProduct += a[i] * b[i]
		magA += a[i] * a[i]
		magB += b[i] * b[i]
	}

	magA = math.Sqrt(magA)
	magB = math.Sqrt(magB)

	if magA == 0 || magB == 0 {
		return 0.0
	}

	return dotProduct / (magA * magB)
}

func processChunk(chunk []string, threshold float64) [][]string {
	localGroups := [][]string{}
	used := make(map[string]bool)

	for _, word1 := range chunk {
		if used[word1] {
			continue
		}

		vectorsLock.RLock()
		vec1, exists := vectors[word1]
		vectorsLock.RUnlock()

		if !exists {
			continue
		}

		group := []string{word1}
		used[word1] = true

		for _, word2 := range chunk {
			if word1 == word2 || used[word2] {
				continue
			}

			vectorsLock.RLock()
			vec2, exists := vectors[word2]
			vectorsLock.RUnlock()

			if !exists {
				continue
			}

			if cosineSimilarity(vec1, vec2) > threshold {
				group = append(group, word2)
				used[word2] = true
			}
		}

		if len(group) > 0 {
			localGroups = append(localGroups, group)
		}
	}

	progressBar.Add(len(chunk))
	return localGroups
}

func groupSimilarWords(words []string, threshold float64) [][]string {
	var (
		groups [][]string
		mu     sync.Mutex
	)

	progressBar = pb.StartNew(len(words))
	progressBar.SetTemplateString(`{{counters . }} {{ bar . "[" "=" ">" " " "]" }} {{percent . }} {{etime . }}`)

	pool := tunny.NewFunc(32, func(payload interface{}) interface{} {
		chunk := payload.([]string)
		return processChunk(chunk, threshold)
	})
	defer pool.Close()

	chunkSize := (len(words) + 31) / 32
	var wg sync.WaitGroup
	results := make(chan [][]string, 32)

	for i := 0; i < 32; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(words) {
			end = len(words)
		}
		if start >= end {
			continue
		}

		wg.Add(1)
		go func(chunk []string) {
			defer wg.Done()
			result := pool.Process(chunk).([][]string)
			results <- result
		}(words[start:end])
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		mu.Lock()
		groups = append(groups, result...)
		mu.Unlock()
	}

	progressBar.Finish()
	return groups
}

func saveGroupsToFile(groups [][]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, group := range groups {
		line := strings.Join(group, ", ")
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func main() {
	startTime := time.Now()

	fmt.Println("Загрузка векторов GloVe...")
	if err := loadGloveVectors("../../data/vectors.txt.txt"); err != nil {
		fmt.Println("Ошибка загрузки векторов:", err)
		return
	}

	fmt.Println("Чтение входного файла...")
	lines, err := readLines("/home/audetv/go/src/github.com/terratensor/glove-pipeline/data/svodd/cleaned_corpus.txt")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	fmt.Println("Очистка текста...")
	var allWords []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	cleanBar := pb.StartNew(len(lines))
	cleanBar.SetTemplateString(`{{counters . }} {{ bar . "[" "=" ">" " " "]" }} {{percent . }} {{etime . }}`)

	cleanPool := tunny.NewFunc(32, func(payload interface{}) interface{} {
		line := payload.(string)
		words := filterStopWords(cleanText(line))
		mu.Lock()
		allWords = append(allWords, words...)
		mu.Unlock()
		cleanBar.Increment()
		return nil
	})
	defer cleanPool.Close()

	for _, line := range lines {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			cleanPool.Process(line)
		}(line)
	}

	wg.Wait()
	cleanBar.Finish()

	fmt.Println("\nГруппировка слов...")
	groups := groupSimilarWords(allWords, 0.7)

	fmt.Println("\nСохранение результатов...")
	if err := saveGroupsToFile(groups, "word_groups.txt"); err != nil {
		fmt.Println("Ошибка сохранения:", err)
		return
	}

	fmt.Printf("\nОбработка завершена за %v\n", time.Since(startTime))
	fmt.Println("Результат сохранен в word_groups.txt")
}
