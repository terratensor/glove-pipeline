#!/bin/bash

# Создание папок
mkdir -p data scripts third_party

# Клонирование и компиляция GloVe
if [ ! -d "third_party/glove" ]; then
    echo "Установка GloVe..."
    cd third_party
    git clone https://github.com/stanfordnlp/glove
    cd glove
    make
    cd ../..
else
    echo "GloVe уже установлен."
fi

# Проверка, что GloVe скомпилирован
if [ ! -f "third_party/glove/build/vocab_count" ]; then
    echo "Ошибка: GloVe не скомпилирован."
    exit 1
fi

# Создание скрипта glove.sh
echo "Создание скрипта glove.sh..."
cat << 'EOF' > scripts/glove.sh
#!/bin/bash

# Путь к исполняемым файлам GloVe
GLOVE_PATH="./third_party/glove/build"

# Входной файл с текстом
INPUT_FILE="./data/cleaned_corpus.txt"

# Выходные файлы
OUTPUT_VOCAB="./data/vocab.txt"
OUTPUT_COOCCURRENCE="./data/cooccurrence.bin"
OUTPUT_VECTORS="./data/vectors.txt"

# Параметры GloVe
VOCAB_MIN_COUNT=5  # Минимальная частота слова для включения в словарь
VECTOR_SIZE=100    # Размерность векторов
WINDOW_SIZE=15     # Размер окна для контекста
ITERATIONS=15      # Количество итераций

# Проверка существования входного файла
if [ ! -f "$INPUT_FILE" ]; then
    echo "Ошибка: файл $INPUT_FILE не найден."
    exit 1
fi

# Проверка существования исполняемых файлов GloVe
if [ ! -f "$GLOVE_PATH/vocab_count" ]; then
    echo "Ошибка: исполняемый файл $GLOVE_PATH/vocab_count не найден."
    exit 1
fi

# Запуск GloVe
echo "Создание словаря..."
"$GLOVE_PATH/vocab_count" -min-count $VOCAB_MIN_COUNT -verbose 2 < "$INPUT_FILE" > "$OUTPUT_VOCAB"

echo "Создание файла совместной встречаемости..."
"$GLOVE_PATH/cooccur" -memory 4.0 -vocab-file "$OUTPUT_VOCAB" -verbose 2 -window-size $WINDOW_SIZE < "$INPUT_FILE" > "$OUTPUT_COOCCURRENCE"

echo "Перемешивание данных..."
"$GLOVE_PATH/shuffle" -memory 4.0 -verbose 2 < "$OUTPUT_COOCCURRENCE" > "$OUTPUT_COOCCURRENCE.shuf.bin"

echo "Обучение GloVe..."
"$GLOVE_PATH/glove" -save-file "$OUTPUT_VECTORS" -threads 8 -input-file "$OUTPUT_COOCCURRENCE.shuf.bin" -x-max 10 -iter $ITERATIONS -vector-size $VECTOR_SIZE -binary 2 -vocab-file "$OUTPUT_VOCAB" -verbose 2

echo "GloVe завершен."
EOF

# Делаем скрипт исполняемым
chmod +x scripts/glove.sh

echo "Инициализация завершена."