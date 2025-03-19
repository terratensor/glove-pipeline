package textcleaner

import (
	"regexp"
	"strings"
)

// CleanText очищает текст от HTML-тегов, ссылок и лишних символов
func CleanText(text string) string {
	// Удаление ссылок, обрамленных тегами <span class="link">...</span>
	re := regexp.MustCompile(`<span class="link">.*?</span>`)
	text = re.ReplaceAllString(text, "")

	// Удаление цитирования (всё, что внутри <blockquote>...</blockquote>)
	re = regexp.MustCompile(`<blockquote[^>]*>.*?</blockquote>`)
	text = re.ReplaceAllString(text, "")

	// Удаление оставшихся HTML-тегов
	re = regexp.MustCompile(`<[^>]+>`)
	text = re.ReplaceAllString(text, "")

	// Удаление обычных ссылок (начинающихся с http или www)
	re = regexp.MustCompile(`https?://\S+|www\.\S+`)
	text = re.ReplaceAllString(text, "")

	// Удаление пунктуации и специальных символов
	re = regexp.MustCompile(`[^a-zA-Zа-яА-Я\s]`)
	text = re.ReplaceAllString(strings.ToLower(text), "")

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
