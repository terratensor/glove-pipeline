Отлично, что у вас получилось с классификацией! Теперь давайте реализуем кластеризацию текстов с использованием алгоритма k-means. Мы будем использовать средние векторы слов для представления текстов и реализуем k-means с нуля на Go.

---

### План реализации:
1. **Подготовка данных**:
- Каждый текст представляется как средний вектор его слов.
- Используем функцию `TextToVector`, которую мы уже реализовали.

2. **Алгоритм k-means**:
- Инициализируем центроиды (начальные центры кластеров).
- На каждом шаге:
- Назначаем каждый текст ближайшему центроиду.
- Пересчитываем центроиды как среднее векторов текстов в кластере.
- Повторяем до сходимости (или до максимального числа итераций).

3. **Оценка результата**:
- Выводим тексты, сгруппированные по кластерам.

---

### Реализация k-means на Go:

```go
package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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

// EuclideanDistance вычисляет евклидово расстояние между двумя векторами.
func EuclideanDistance(vec1, vec2 []float64) float64 {
	if len(vec1) != len(vec2) {
		return math.Inf(1) // Возвращаем +Inf, если векторы разной длины
	}

	var sum float64
	for i := range vec1 {
		sum += math.Pow(vec1[i]-vec2[i], 2)
	}
	return math.Sqrt(sum)
}

// KMeans реализует алгоритм k-means.
func KMeans(data [][]float64, k int, maxIterations int) ([]int, [][]float64) {
	// Инициализация центроидов случайным образом
	rand.Seed(time.Now().UnixNano())
	centroids := make([][]float64, k)
	for i := range centroids {
		centroids[i] = data[rand.Intn(len(data))]
	}

	// Массив для хранения меток кластеров
	labels := make([]int, len(data))

	for iter := 0; iter < maxIterations; iter++ {
		// Шаг 1: Назначение кластеров
		for i, point := range data {
			minDist := math.Inf(1)
			for j, centroid := range centroids {
				dist := EuclideanDistance(point, centroid)
				if dist < minDist {
					minDist = dist
					labels[i] = j
				}
			}
		}

		// Шаг 2: Пересчет центроидов
		newCentroids := make([][]float64, k)
		counts := make([]int, k)
		for i := range newCentroids {
			newCentroids[i] = make([]float64, len(data[0]))
		}

		for i, point := range data {
			cluster := labels[i]
			for j := range point {
				newCentroids[cluster][j] += point[j]
			}
			counts[cluster]++
		}

		for i := range newCentroids {
			if counts[i] > 0 {
				for j := range newCentroids[i] {
					newCentroids[i][j] /= float64(counts[i])
				}
			}
		}

		// Проверка на сходимость
		converged := true
		for i := range centroids {
			if EuclideanDistance(centroids[i], newCentroids[i]) > 1e-6 {
				converged = false
				break
			}
		}

		if converged {
			break
		}

		centroids = newCentroids
	}

	return labels, centroids
}

func main() {
	// Загрузка векторов
	vectors, err := LoadVectors("vectors.txt")
	if err != nil {
		fmt.Println("Ошибка загрузки векторов:", err)
		return
	}

	if len(vectors) == 0 {
		fmt.Println("Векторы не загружены или файл пуст.")
		return
	}

	// Пример текстов для кластеризации
	texts := []string{
		"кот собака мышь",
		"погода солнце дождь",
		"компьютер программа алгоритм",
		"кофе чай напиток",
		"автомобиль дорога скорость",
		"книга библиотека чтение",
	}

	// Преобразуем тексты в векторы
	var data [][]float64
	for _, text := range texts {
		vector := TextToVector(text, vectors)
		if vector != nil {
			data = append(data, vector)
		}
	}

	if len(data) == 0 {
		fmt.Println("Нет данных для кластеризации.")
		return
	}

	// Кластеризация с использованием k-means
	k := 2 // Количество кластеров
	labels, centroids := KMeans(data, k, 100)

	// Вывод результатов
	fmt.Println("Центроиды:")
	for i, centroid := range centroids {
		fmt.Printf("Кластер %d: %v\n", i, centroid)
	}

	fmt.Println("\nРезультаты кластеризации:")
	for i, text := range texts {
		if i < len(labels) {
			fmt.Printf("Текст: '%s' -> Кластер %d\n", text, labels[i])
		}
	}
}

```

---

### Объяснение кода:
1. **EuclideanDistance**:
- Вычисляет евклидово расстояние между двумя векторами. Это метрика, которая используется для определения ближайшего центроида.

2. **KMeans**:
- Реализует алгоритм k-means. На каждом шаге:
- Назначает каждый текст ближайшему центроиду.
- Пересчитывает центроиды как среднее векторов текстов в кластере.
- Алгоритм завершается, когда центроиды перестают изменяться или достигается максимальное число итераций.

3. **main**:
- Загружает векторы слов.
- Преобразует тексты в векторы.
- Выполняет кластеризацию с помощью k-means.
- Выводит результаты.

---

### Пример вывода:
```
Центроиды:
Кластер 0: [0.25 0.35 0.45 0.55]
Кластер 1: [0.65 0.75 0.85 0.95]

Результаты кластеризации:
Текст: 'кот собака мышь' -> Кластер 0
Текст: 'погода солнце дождь' -> Кластер 0
Текст: 'компьютер программа алгоритм' -> Кластер 1
Текст: 'кофе чай напиток' -> Кластер 0
Текст: 'автомобиль дорога скорость' -> Кластер 1
Текст: 'книга библиотека чтение' -> Кластер 1
```

---

### Как улучшить:
1. **Выбор числа кластеров (k)**:
- Используйте метод "локтя" или другие методы для выбора оптимального числа кластеров.

2. **Инициализация центроидов**:
- Используйте более сложные методы инициализации, например, k-means++.

3. **Обработка шума**:
- Если в данных есть шум, можно использовать алгоритм DBSCAN, который не требует задания числа кластеров и устойчив к шуму.

4. **Визуализация**:
- Визуализируйте результаты кластеризации с помощью PCA или t-SNE для уменьшения размерности данных.

---

Теперь у вас есть реализация кластеризации текстов на Go! Вы можете адаптировать этот код для своих данных и задач.