package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Список стоп-слов для русского языка
var stopWords = map[string]bool{
	"и": true, "в": true, "не": true, "на": true, "я": true,
	"с": true, "что": true, "а": true, "по": true, "как": true,
	"но": true, "к": true, "у": true, "же": true, "вы": true,
	"то": true, "о": true, "из": true, "за": true, "от": true,
	"для": true, "это": true, "так": true, "все": true, "его": true,
	"она": true, "они": true, "мы": true, "бы": true, "ее": true,
	"ещё": true, "еще": true, "уже": true, "или": true, "если": true,
	"только": true, "когда": true, "даже": true, "нет": true,
}

func main() {
	// Параметры
	n := 5            // Размер n-граммы
	minFrequency := 5 // Минимальная частота для сохранения

	// Загрузка файла
	file, err := os.Open("../data/cleaned_corpus.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Чтение файла и разбиение на слова
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	var words []string
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text()) // Приводим к нижнему регистру
		if !stopWords[word] {                   // Игнорируем стоп-слова
			words = append(words, word)
		}
	}

	// Составление словаря n-грамм
	ngramFrequency := make(map[string]int)
	for i := 0; i <= len(words)-n; i++ {
		ngram := strings.Join(words[i:i+n], " ")
		ngramFrequency[ngram]++
	}

	// Фильтрация n-грамм по частоте
	var filteredNgrams []string
	for ngram, freq := range ngramFrequency {
		if freq >= minFrequency {
			filteredNgrams = append(filteredNgrams, ngram)
		}
	}

	// Сортировка по частоте
	sort.Slice(filteredNgrams, func(i, j int) bool {
		return ngramFrequency[filteredNgrams[i]] > ngramFrequency[filteredNgrams[j]]
	})

	// Сохранение результата в файл
	outputFile, err := os.Create("ngram_output.txt")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outputFile.Close()

	for _, ngram := range filteredNgrams {
		line := fmt.Sprintf("%s: %d\n", ngram, ngramFrequency[ngram])
		outputFile.WriteString(line)
	}

	fmt.Println("Готово! Результат сохранен в ngram_output.txt")
}
