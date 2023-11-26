package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Структура для представления формата JSON-файла.
type Numbers struct {
	Values []int `json:"values"`
}

func main() {
	// Инициализация параметров командной строки
	fileFlag := flag.String("file", "", "Path to JSON file")
	outputFlag := flag.String("output", "", "Path to output file")
	flag.Parse()

	// Инициализация конфигурации через файл и переменные окружения
	initConfig()

	// Определение источника данных
	var byteValue []byte
	var source string

	if *fileFlag != "" {
		// Если указан файл в аргументах командной строки, читаем из файла
		byteValue = readFile(*fileFlag)
		source = fmt.Sprintf("File: %s", *fileFlag)
	} else {
		// Иначе читаем из стандартного ввода
		byteValue, _ = ioutil.ReadAll(os.Stdin)
		source = "Stdin"
	}

	// Декодируем JSON
	var numbers Numbers
	if err := json.Unmarshal(byteValue, &numbers); err != nil {
		log.Fatalf("Error decoding JSON from %s: %s", source, err)
	}

	// Считаем сумму чисел
	sum := 0
	for _, num := range numbers.Values {
		sum += num
	}

	// Логируем сумму чисел
	log.Printf("Sum of numbers: %d", sum)

	// Выполняем HTTP GET запрос
	url := viper.GetString("http_url")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making HTTP request to %s: %s", url, err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("HTTP request failed to %s. Status: %s", url, resp.Status)
	}

	// Логируем успешный запрос
	log.Printf("HTTP request to %s was successful. Status: %s", url, resp.Status)

	// Выводим результат в файл или стандартный вывод
	writeOutput(*outputFlag, sum)
}

func initConfig() {
	// Инициализация конфигурации из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default values.")
	}

	// Инициализация конфигурации из переменных окружения
	viper.AutomaticEnv()

	// Установка значений по умолчанию
	viper.SetDefault("http_url", "https://example.com")
}

func readFile(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %s", filePath, err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file %s: %s", filePath, err)
	}

	return byteValue
}

func writeOutput(outputFilePath string, result int) {
	// Если указан путь к файлу вывода, записываем туда
	if outputFilePath != "" {
		file, err := os.Create(outputFilePath)
		if err != nil {
			log.Fatalf("Error creating output file %s: %s", outputFilePath, err)
		}
		defer file.Close()

		file.WriteString(fmt.Sprintf("Sum of numbers: %d\n", result))
		log.Printf("Output written to file: %s", outputFilePath)
	} else {
		// Иначе выводим в стандартный вывод
		fmt.Printf("Sum of numbers: %d\n", result)
	}
}
