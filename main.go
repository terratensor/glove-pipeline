package main

import (
	"flag"
	"fmt"
	"glove-pipeline/pkg/glove"
	"glove-pipeline/pkg/ngrams"
	"glove-pipeline/pkg/textprocessor"
	"log"
)

func main() {
	// Определение флагов
	cleanTextFlag := flag.Bool("clean", false, "Запустить только очистку текста")
	runGloveFlag := flag.Bool("glove", false, "Запустить только GloVe")
	extractNGramsFlag := flag.Bool("ngrams", false, "Запустить только извлечение n-грамм")
	n := flag.Int("n", 2, "Размер n-грамм (2 для биграмм, 3 для триграмм и т.д.)")
	topN := flag.Int("top", 10, "Количество топ-N n-грамм для вывода в консоль")
	useStopwords := flag.Bool("stopwords", false, "Учитывать стоп-слова при формировании n-грамм")
	flag.Parse()

	// Если флаги не указаны, запустить полный pipeline
	if !*cleanTextFlag && !*runGloveFlag && !*extractNGramsFlag {
		fullPipeline(*n, *topN, *useStopwords)
		return
	}

	// Запуск отдельных шагов
	if *cleanTextFlag {
		cleanText()
	}
	if *runGloveFlag {
		runGlove()
	}
	if *extractNGramsFlag {
		extractNGrams(*n, *topN, *useStopwords)
	}
}

// fullPipeline запускает полный pipeline
func fullPipeline(n int, topN int, useStopwords bool) {
	fmt.Println("Запуск полного pipeline...")
	cleanText()
	runGlove()
	extractNGrams(n, topN, useStopwords)
}

// cleanText выполняет очистку текста
func cleanText() {
	fmt.Println("Шаг 1: Очистка текста...")
	inputFile := "data/input.csv"
	outputFile := "data/cleaned_corpus.txt"

	err := textprocessor.ProcessCSV(inputFile, outputFile)
	if err != nil {
		log.Fatalf("Ошибка при очистке текста: %v", err)
	}
}

// runGlove запускает GloVe
func runGlove() {
	fmt.Println("Шаг 2: Запуск GloVe...")
	err := glove.Run()
	if err != nil {
		log.Fatalf("Ошибка при запуске GloVe: %v", err)
	}
}

// extractNGrams извлекает n-граммы
func extractNGrams(n int, topN int, useStopwords bool) {
	fmt.Printf("Шаг 3: Извлечение %d-грамм...\n", n)
	cooccurrenceFile := "data/cooccurrence.bin"
	vocabFile := "data/vocab.txt"
	stopwordsFile := "data/stopwords.txt"
	topNGrams, allNGrams, err := ngrams.ExtractNGrams(cooccurrenceFile, vocabFile, stopwordsFile, n, topN, useStopwords)
	if err != nil {
		log.Fatalf("Ошибка при извлечении n-грамм: %v", err)
	}

	// Сохранение всех n-грамм в файл
	outputFile := fmt.Sprintf("data/%d_grams.txt", n)
	err = ngrams.SaveNGrams(allNGrams, outputFile)
	if err != nil {
		log.Fatalf("Ошибка при сохранении n-грамм: %v", err)
	}

	// Вывод топ-N n-грамм в консоль
	fmt.Printf("Топ-%d %d-грамм:\n", topN, n)
	for _, pair := range topNGrams {
		fmt.Printf("%v: %f\n", pair.Words, pair.Frequency)
	}
}
