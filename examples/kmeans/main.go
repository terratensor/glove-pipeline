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
	vectors, err := LoadVectors("../../data/vectors.txt.txt")
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
		"возможно необходима параллельная спецоперация по зачистке сми от окопавшихся там власовцев сейчас многие из них могут ярче проявиться а мб и какието из спящих активизируются новая старая тема",
		"так у большинства калейдоскоп кал лей до скоп или мозайка коб это сложенная картина мира из культурного наследия где в мире вс управляемо берите и исследуйте он видоизменяется сми кино музыка может видоизменяться неоднократно не было изложено структурированно с примерами и доказательствами на практике в письменном виде скорее это бога закон в коб идт многое изложенное пушкиным как и других классиков эволюция не стоит на месте новое время новые события новые достижения новые писатели скорее воссоединялись потому как их всю дорогу разъединяют скорее не при завоеваниях а при возвращениях когото смог тех кто когдато тоже были мы вероятно следует вносить отдельный урок по типу соборной темы кто знает тот поймет кто не знает тому и понимать нечего с",
		"компьютер программа алгоритм",
		"кофе чай напиток",
		"заместитель председателя комитета госдумы по обороне алексей журавлв заявил о необходимости подготовки мужского населения россии к возможной мобилизации в интервью абзацу он подчеркнул что страны европы и их союзники могут быть готовы к войне с россией к годам по его словам в настоящий момент стране хватает добровольцев но необходимо учитывать долгосрочные угрозы политик считает что подготовка должна включать чткое функционирование военкоматов и создание эффективного мобилизационного резерва журавлв отметил что запад сам заявляет о своих военных намерениях а потому россия должна быть готова к защите своей территории он добавил что такие вопросы нельзя замалчивать поскольку они касаются будущей безопасности страны депутат госдумы андрей гурулев ранее рассказал что несмотря на возможные сюрпризы со стороны штатов и незалежной в россии в году не будут объявлять мобилизацию по его словам сейчас нет предпосылок к мобилизации он уточнил что сегодня люди массово идут добровольцами на контрактную службу очередные модули о войне в интересах сша в европе по аналогу мировой подпиндосники готовы участвовать в договорняке и клянутся в верности хозяину в сша причем даат указана когда срок трампа будет подходить к концу вот такие депутаты имеются провокаторы на марше",
		"это у них есть такая иллюзия что это в принципе возможно окопаться и закрепиться за суток оно возможно только если у тебя есть готовое выстроенное логистическое плечо отуда оно у франции или германии с британией это смехотворно в ес сегодня даже нет так называемого военного шенгена то есть пока грузы едут из одной директории в другую буквально каждые километров они останавливаются на сутки не меньше просто потому что так устроено в европейской бюрократии и изменить это положение дел невозможно абсолютно для этого потребуется много бюрократической возни на годы вперд нет тридцати дней не хватит на окопаться но перегруппироваться укомплектовать более боеспособные подразделения в целом можно если есть человеческий ресурс и основной вопрос путин уже спросил притом он сразу же ответил на него чтобы у запада не было иллюзий кто определит где и кто нарушил возможную договорнность о прекращении огня на протяжении двух тысяч километров и потом кто на кого будет сваливать нарушение этой договорнности если ктото хочет чемто воспользоваться и прямо так обмануть россию то пусть не удивляются почему это поначалу оказалось так просто сделать а потом тысяч натовских трупов и приказ конечно берите в плен если возможно и соответствует боевой обстановке",
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
	k := 4 // Количество кластеров
	labels, centroids := KMeans(data, k, 1000)

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
