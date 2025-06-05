package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/teris-io/shortid"
)

var temp *template.Template
var err error
var db *sql.DB

func initDB(filepath string) *sql.DB {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Printf("DATABASE file '%s' not found, creating it.", filepath)
	}

	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	createStatement := `
	CREATE TABLE IF NOT EXISTS shorts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_url TEXT NOT NULL,
		shorten_url TEXT NOT NULL
	);
	`

	_, err = database.Exec(createStatement)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	log.Println("Database initialized and 'shorts' table ensured.")
	return database
}

func main() {
	temp, err = template.ParseFiles("pages/home.gohtml")
	if err != nil {
		log.Fatal("error parsing home template", err)
	}

	db = initDB("./shorts.db")
	defer db.Close()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/urls", URLsHandler)
	http.HandleFunc("/delete", deleteHandler)

	fmt.Println("Running server on port :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Couldn't run server...")
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := temp.Execute(w, nil)
	if err != nil {
		log.Fatalln("error executing template home")
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		urlValue := r.FormValue("url")
		if urlValue == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		insertStatemet := `INSERT INTO shorts (original_url, shorten_url) VALUES (?,?)`
		shortURLValue, err := shortid.Generate()
		if err != nil {
			log.Printf("Error generating random uuid: %v", err)
			http.Error(w, "Failed to generate short url", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(insertStatemet, urlValue, shortURLValue)
		if err != nil {
			log.Printf("Error inserting url: %v", err)
			http.Error(w, "Failed to shorten and save url", http.StatusInternalServerError)
			return
		}
		fmt.Println("Added url", urlValue)
		http.Redirect(w, r, "/", http.StatusCreated)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type ShortURL struct {
	ID           int
	Original_url string
	Shorten_url  string
}

func URLsHandler(w http.ResponseWriter, r *http.Request) {
	readStatement := "SELECT * from shorts"
	rows, err := db.Query(readStatement, nil)
	if err != nil {
		log.Fatal("error", err)
	}
	defer rows.Close()
	var urls []ShortURL
	for rows.Next() {
		var current ShortURL
		err := rows.Scan(&current.ID, &current.Original_url, &current.Shorten_url)
		if err != nil {
			log.Printf("Error scanning db row: %v", err)
			http.Error(w, "Error processing database results", http.StatusInternalServerError)
			return
		}
		urls = append(urls, current)
	}

	t, err := template.ParseFiles("pages/urls.gohtml")
	if err != nil {
		fmt.Printf("error parsing template urls.gohtml: %v", err)
		http.Error(w, "error parsing template urls.gohtml", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, urls)
	if err != nil {
		fmt.Printf("error executing template urls.gohtml: %v", err)
		http.Error(w, "error executing template urls.gohtml", http.StatusInternalServerError)
		return
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get the body id
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "error parsing form data", http.StatusInternalServerError)
		return
	}
	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "id can not be empty", http.StatusBadRequest)
		return
	}
	// convert it to num
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "error converting 'string id' to 'int id'", http.StatusInternalServerError)
		return
	}
	// check if id exist
	selectStatement := `select * from shorts where (id) = (?)`
	rows, err := db.Query(selectStatement, id)
	if err != nil {
		http.Error(w, "error fetching urls", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var row ShortURL
		err := rows.Scan(&row.ID, &row.Original_url, &row.Shorten_url)
		if err != nil {
			http.Error(w, "error scanning row", http.StatusInternalServerError)
			return
		}
	}
	// delete
	// if ok, redirect back
	// if not, show error
}
