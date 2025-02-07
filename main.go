package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

// var texts = []string{
// 	"Программирование на Go - это увлекательный процесс...",
// 	"Яндекс предлагает стажерам работать над реальными проектами...",
// 	"Скорость печати важна для разработчика...",
// }

var db *sql.DB
var err error

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func loadConfig() (*Config, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
		// return nil, fmt.Errorf("DB_HOST не установлен")
	}

	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 5432
		// return nil, fmt.Errorf("неверный формат DB_PORT: %v", err)
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "user"
		// return nil, fmt.Errorf("DB_USER не установлен")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "user"
		// return nil, fmt.Errorf("DB_PASSWORD не установлен")
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "db"
		// return nil, fmt.Errorf("DB_NAME не установлен")
	}

	return &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
	}, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	config, err := loadConfig()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to db!")

	// Статические файлы
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API endpoints
	http.HandleFunc("/api/text", enableCORS(textHandler))
	http.HandleFunc("/api/check", enableCORS(checkHandler))
	http.HandleFunc("/api/save", enableCORS(saveHandler))

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func getRandomText() (string, error) {
	query := "SELECT text_for_test FROM speed_texts ORDER BY RANDOM() LIMIT 1;"

	var text string
	err := db.QueryRow(query).Scan(&text)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("тексты не найдены")
		}
		return "", fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	return text, nil
}

// Обработчики API
func textHandler(w http.ResponseWriter, r *http.Request) {
	text, err := getRandomText()
	if err != nil {
		panic(err)
	}
	response := struct {
		Text string `json:"text"`
	}{
		Text: text,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next(w, r)
	}
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Input    string `json:"input"`
		Original string `json:"original"`
	}

	type Response struct {
		WPM      int     `json:"wpm"`
		Accuracy float64 `json:"accuracy"`
		Time     float64 `json:"time"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Логика расчета
	correct := calculateCorrectChars(req.Original, req.Input)
	accuracy := float64(correct) / float64(len(req.Original)) * 100
	timeSpent := 60.0 // Пример значения

	wpm := int(float64(len(req.Input)/5) / (timeSpent / 60))

	json.NewEncoder(w).Encode(Response{
		WPM:      wpm,
		Accuracy: accuracy,
		Time:     timeSpent,
	})
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	// Логика сохранения
	w.WriteHeader(http.StatusOK)
}

// Вспомогательные функции
func calculateCorrectChars(original, input string) int {
	correct := 0
	minLen := len(original)
	if len(input) < minLen {
		minLen = len(input)
	}

	for i := 0; i < minLen; i++ {
		if original[i] == input[i] {
			correct++
		}
	}
	return correct
}
