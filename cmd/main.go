package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		name := r.FormValue("name")
		age := r.FormValue("age")
		description := r.FormValue("description")
		fmt.Println("Введенные данные: ", name, age, description, "\n")
		fmt.Printf("Введенные данные: %t, %t, %t \n", name, age, description)
		err := addUser(db, name, age, description)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("***********")
		fmt.Println(r.Method)
		fmt.Println(r.Header)
		fmt.Println(r.Body)
		fmt.Println(r.URL)
		fmt.Println(r.Host)
		fmt.Println("***********")

	}

	tmplParsed, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		fmt.Println("Ошибка загрузки шаблона:", err)
		return
	}
	tmplParsed.Execute(w, nil)
}

func addUser(db *sql.DB, name, age, description string) error {
	insert := `INSERT INTO users (name, age, description) VALUES ($1, $2, $3)`
	_, err := db.Exec(insert, name, age, description)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

func clearHandler(w http.ResponseWriter, r *http.Request) {
	delete := `DELETE FROM users`
	if _, err := db.Exec(delete); err != nil {
		http.Error(w, "Failed to clear database", http.StatusInternalServerError)
		return
	}
	tmplParsed, err := template.ParseFiles("./templates/clear.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		fmt.Println("Ошибка загрузки шаблона:", err)
		return
	}
	tmplParsed.Execute(w, nil)
}

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "user"
	dbname   = "db"
)

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
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
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/clear", clearHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
