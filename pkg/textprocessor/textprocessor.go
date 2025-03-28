package textprocessor

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// CleanText очищает текст от HTML-тегов, ссылок и лишних символов
func CleanText(text string) string {
	// Удаление HTML-сущностей (включая &quot;, &#34; и другие)
	re := regexp.MustCompile(`&[a-z]+;|&#\d+;`)
	text = re.ReplaceAllString(text, " ")

	// Удаление ссылок, обрамленных тегами <span class="link">...</span>
	re = regexp.MustCompile(`<span class="link">.*?</span>`)
	text = re.ReplaceAllString(text, " ")

	// Удаление цитирования (всё, что внутри <blockquote>...</blockquote>)
	re = regexp.MustCompile(`<blockquote[^>]*>.*?</blockquote>`)
	text = re.ReplaceAllString(text, " ")

	// Удаление оставшихся HTML-тегов
	re = regexp.MustCompile(`<[^>]+>`)
	text = re.ReplaceAllString(text, " ")

	// Удаление обычных ссылок (начинающихся с http или www)
	re = regexp.MustCompile(`https?://\S+|www\.\S+`)
	text = re.ReplaceAllString(text, " ")

	// Удаление пунктуации и специальных символов (оставляем только буквы, цифры и пробелы)
	// Сохраняем букву ё (Ё) явно в регулярном выражении
	re = regexp.MustCompile(`[^a-zA-Zа-яА-ЯёЁ0-9\s]`)
	text = re.ReplaceAllString(text, " ")

	// Приведение текста к нижнему регистру с сохранением буквы ё
	text = strings.Map(func(r rune) rune {
		if r == 'Ё' {
			return 'ё'
		}
		return unicode.ToLower(r)
	}, text)

	return text
}

// RemoveExcessNewlines удаляет лишние переносы строк и пробелы
func RemoveExcessNewlines(text string) string {
	// Удаление множественных пробелов
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	// Удаление множественных переносов строк
	re = regexp.MustCompile(`\n+`)
	text = re.ReplaceAllString(text, "\n")

	// Удаление пробелов и переносов строк в начале и конце текста
	text = strings.TrimSpace(text)

	return text
}

// ProcessCSV обрабатывает CSV-файл, очищает текст и сохраняет результат в файл
func ProcessCSV(inputFile, outputFile string) error {
	// Открытие CSV-файла
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла: %v", err)
	}
	defer file.Close()

	// Создание CSV-ридера
	reader := csv.NewReader(file)
	reader.Comma = ','       // Указываем разделитель
	reader.LazyQuotes = true // Разрешаем "ленивые" кавычки

	// Открытие файла для записи очищенного текста
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %v", err)
	}
	defer output.Close()

	// Пропуск первой строки (заголовка)
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("ошибка при чтении заголовка CSV: %v", err)
	}

	// Чтение и обработка данных
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ошибка при чтении CSV: %v", err)
		}

		// Предполагаем, что текст находится в первом столбце
		if len(record) > 0 {
			text := record[0]
			cleanedText := CleanText(text)
			cleanedText = RemoveExcessNewlines(cleanedText)

			// Запись в файл, если текст не пустой
			if cleanedText != "" {
				_, err := output.WriteString(cleanedText + "\n")
				if err != nil {
					return fmt.Errorf("ошибка при записи в файл: %v", err)
				}
			}
		}
	}

	log.Printf("Очищенный корпус сохранен в файл %s\n", outputFile)
	return nil
}
