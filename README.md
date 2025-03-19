# GloVe Pipeline Project

Этот проект реализует pipeline для обработки текста, обучения модели GloVe и извлечения n-грамм (биграмм, триграмм и т.д.). Проект написан на Go и использует утилиту GloVe для создания векторных представлений слов.

---

## Оглавление
1. [Установка](#установка)
2. [Запуск](#запуск)
- [Полный pipeline](#полный-pipeline)
- [Отдельные шаги](#отдельные-шаги)
3. [Структура проекта](#структура-проекта)
4. [Файл `cooccurrence.bin`](#файл-cooccurrencebin)
5. [Примеры использования](#примеры-использования)
6. [Особенности](#особенности)

---

## Установка

### 1. Установка Go
Убедитесь, что у вас установлен Go (версия 1.16 или выше). Если нет, скачайте и установите его с [официального сайта](https://golang.org/dl/).

### 2. Клонирование репозитория
Клонируйте репозиторий с проектом:
```bash
git clone https://github.com/terratensor/glove-pipeline.git
cd glove-pipeline
```

### 3. Установка GloVe
GloVe используется для создания векторных представлений слов. Установите его:
```bash
mkdir -p third_party
cd third_party
git clone https://github.com/stanfordnlp/glove
cd glove
make
cd ../..
```

### 4. Установка зависимостей
Инициализируйте Go модуль и установите зависимости:
```bash
go mod init glove-pipeline
go mod tidy
```

---

## Запуск

### Полный pipeline
Запуск полного pipeline (очистка текста, обучение GloVe, извлечение n-грамм):
```bash
go run main.go -n 2
```
- `-n`: Размер n-грамм (2 для биграмм, 3 для триграмм и т.д.).

### Отдельные шаги
Вы можете запускать отдельные шаги pipeline:

1. **Очистка текста**:
```bash
go run main.go -clean
```

2. **Обучение GloVe**:
```bash
go run main.go -glove
```

3. **Извлечение n-грамм**:
```bash
go run main.go -ngrams -n 3 -top 10
```
- `-n`: Размер n-грамм.
- `-top`: Количество топ-N n-грамм для вывода в консоль.

---

## Структура проекта

```
glove-pipeline/
├── data/ # Входные и выходные данные
│ ├── input.csv # Входной CSV-файл с текстом
│ ├── cleaned_corpus.txt # Очищенный текст
│ ├── vocab.txt # Словарь, созданный GloVe
│ ├── cooccurrence.bin # Файл совместной встречаемости
│ ├── vectors.txt # Векторные представления слов
│ └── {n}_grams.txt # Файл с n-граммами (например, 2_grams.txt)
├── scripts/ # Скрипты для запуска GloVe
│ └── glove.sh # Скрипт для запуска GloVe
├── third_party/ # Сторонние зависимости
│ └── glove/ # Исходный код GloVe
├── pkg/ # Пакеты Go
│ ├── textprocessor/ # Очистка текста
│ ├── glove/ # Запуск GloVe
│ └── ngrams/ # Извлечение n-грамм
├── main.go # Основной файл для запуска pipeline
├── init.sh # Скрипт инициализации проекта
└── README.md # Документация
```

---

## Файл `cooccurrence.bin`

### Структура
Файл `cooccurrence.bin` содержит данные о совместной встречаемости слов в формате:
```
| word1 (int32, 4 байта) | word2 (int32, 4 байта) | val (float64, 8 байт) |
```
- `word1`, `word2`: Индексы слов в словаре.
- `val`: Частота совместной встречаемости слов.

### Использование
Этот файл используется для извлечения n-грамм. Мы читаем его, чтобы получить частоту совместной встречаемости слов и формировать n-граммы.

---

## Примеры использования

### 1. Очистка текста и обучение GloVe
```bash
go run main.go -clean
go run main.go -glove
```

### 2. Извлечение биграмм
```bash
go run main.go -ngrams -n 2 -top 10
```

### 3. Извлечение триграмм
```bash
go run main.go -ngrams -n 3 -top 10
```

### 4. Полный pipeline для триграмм
```bash
go run main.go -n 3
```

---

## Особенности

1. **Очистка текста**:
- Удаление HTML-тегов, ссылок, пунктуации и специальных символов.
- Приведение текста к нижнему регистру.

2. **Обучение GloVe**:
- Используется утилита GloVe для создания векторных представлений слов.
- Генерируются файлы `vocab.txt`, `cooccurrence.bin` и `vectors.txt`.

3. **Извлечение n-грамм**:
- Поддерживаются n-граммы любого порядка (биграммы, триграммы и т.д.).
- Все n-граммы сохраняются в файл `data/{n}_grams.txt`.
- Топ-N n-грамм выводятся в консоль.

4. **Логирование**:
- Программа логирует прогресс обработки файлов, что помогает отслеживать выполнение.

---

## Лицензия
Этот проект распространяется под лицензией MIT. Подробности см. в файле [LICENSE](LICENSE).