package ngrams

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// Pair представляет n-грамму и её частоту
type Pair struct {
	Words     []string
	Frequency float64
}

// ByFrequency реализует sort.Interface для сортировки пар по частоте
type ByFrequency []Pair

func (a ByFrequency) Len() int           { return len(a) }
func (a ByFrequency) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFrequency) Less(i, j int) bool { return a[i].Frequency > a[j].Frequency } // Сортировка по убыванию

// ExtractNGrams извлекает топ-N n-грамм из бинарного файла совместной встречаемости
func ExtractNGrams(cooccurrenceFile string, vocabFile string, n int, topN int) ([]Pair, []Pair, error) {
	// Открытие файла совместной встречаемости
	cooccurFile, err := os.Open(cooccurrenceFile)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка при открытии файла совместной встречаемости: %v", err)
	}
	defer cooccurFile.Close()

	// Загрузка словаря
	vocab, err := loadVocab(vocabFile)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка при загрузке словаря: %v", err)
	}

	// Логирование начала обработки
	log.Println("Начало обработки файла совместной встречаемости...")
	startTime := time.Now()

	// Чтение бинарного файла
	cooccurrenceData := make(map[string]float64)
	var lineCount int
	var prevWords []int32 // Хранение предыдущих слов для формирования n-грамм

	for {
		var word1, word2 int32
		var val float64

		// Чтение word1 (int32)
		if err := binary.Read(cooccurFile, binary.LittleEndian, &word1); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, nil, fmt.Errorf("ошибка при чтении word1: %v", err)
		}

		// Чтение word2 (int32)
		if err := binary.Read(cooccurFile, binary.LittleEndian, &word2); err != nil {
			return nil, nil, fmt.Errorf("ошибка при чтении word2: %v", err)
		}

		// Чтение val (float64)
		if err := binary.Read(cooccurFile, binary.LittleEndian, &val); err != nil {
			return nil, nil, fmt.Errorf("ошибка при чтении val: %v", err)
		}

		// Добавление слов в историю
		prevWords = append(prevWords, word1, word2)

		// Формирование n-грамм
		if len(prevWords) >= n {
			// Получение n-граммы
			ngramWords := prevWords[len(prevWords)-n:]

			// Проверка, что все слова есть в словаре
			valid := true
			ngramStr := make([]string, n)
			for i, word := range ngramWords {
				wordStr, ok := vocab[word]
				if !ok {
					valid = false
					break
				}
				ngramStr[i] = wordStr
			}

			if valid {
				// Сохранение n-граммы и её частоты
				key := strings.Join(ngramStr, " ")
				cooccurrenceData[key] += val
			}
		}

		// Логирование прогресса
		lineCount++
		if lineCount%1000000 == 0 {
			log.Printf("Обработано %d строк...", lineCount)
		}
	}

	// Логирование завершения обработки
	log.Printf("Обработка завершена. Всего обработано %d строк за %v.", lineCount, time.Since(startTime))

	// Преобразование мапы в срез пар для сортировки
	var pairs []Pair
	for key, freq := range cooccurrenceData {
		words := strings.Split(key, " ")
		pairs = append(pairs, Pair{Words: words, Frequency: freq})
	}

	// Сортировка пар по частоте
	sort.Sort(ByFrequency(pairs))

	// Ограничение результата топ-N n-граммами
	if topN > len(pairs) {
		topN = len(pairs)
	}

	return pairs[:topN], pairs, nil
}

// loadVocab загружает словарь из файла vocab.txt
func loadVocab(vocabFile string) (map[int32]string, error) {
	file, err := os.Open(vocabFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла словаря: %v", err)
	}
	defer file.Close()

	vocab := make(map[int32]string)
	scanner := bufio.NewScanner(file)
	var index int32 = 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 1 {
			continue
		}
		word := parts[0]
		vocab[index] = word
		index++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении файла словаря: %v", err)
	}

	return vocab, nil
}

// SaveNGrams сохраняет n-граммы в файл
func SaveNGrams(ngrams []Pair, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %v", err)
	}
	defer file.Close()

	for _, pair := range ngrams {
		line := fmt.Sprintf("%v: %f\n", pair.Words, pair.Frequency)
		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("ошибка при записи в файл: %v", err)
		}
	}

	return nil
}
