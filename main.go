package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/teris-io/shortid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var temp *template.Template
var err error
var db *gorm.DB

type ShortURL struct {
	gorm.Model
	URL  string
	SURL string
}

func main() {
	r := chi.NewRouter()

	db, err = gorm.Open(sqlite.Open("shorts.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&ShortURL{})

	// TODO : Fix template parsing
	temp, err = template.ParseFiles("pages/home.gohtml")
	if err != nil {
		log.Fatal("error parsing home template", err)
	}
	r.Get("/", homeHandler)
	r.Post("/submit", submitHandler)
	r.Get("/urls", URLsHandler)
	r.Post("/delete", deleteHandler)

	log.Println("running server on port :8080")
	if err = http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("error running server")
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
		// Todo: valide string if it valid url
		urlValue := r.FormValue("url")
		if urlValue == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		shortURLValue, err := shortid.Generate()
		if err != nil {
			log.Printf("Error generating random uuid: %v", err)
			http.Error(w, "Failed to generate short url", http.StatusInternalServerError)
			return
		}
		id := db.Create(&ShortURL{URL: urlValue, SURL: shortURLValue})
		fmt.Println("Added url", id)
		http.Redirect(w, r, "/urls", http.StatusMovedPermanently)
		return
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func URLsHandler(w http.ResponseWriter, r *http.Request) {
	var urls []ShortURL
	result := db.Find(&urls)
	if result.Error != nil {
		fmt.Fprint(w, "error fetching urls")
	}

	t, err := template.ParseFiles("pages/urls.gohtml")
	if err != nil {
		fmt.Printf("error parsing template urls.gohtml: %v", err)
		http.Error(w, "error parsing template urls.gohtml", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, struct {
		Host string
		URLS []ShortURL
	}{
		Host: r.Host,
		URLS: urls,
	})
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
	var url ShortURL
	// Todo: check if we need to convert idStr to id
	result := db.First(&url, id)
	if result.Error != nil {
		http.Error(w, "error fetching url", http.StatusNotFound)
		return
	}
	// delete
	result = db.Delete(&ShortURL{}, id)
	if result.Error != nil {
		http.Error(w, "error deleting url", http.StatusInternalServerError)
		return
	}
	// if ok, redirect back
	if result.RowsAffected == 1 {
		http.Redirect(w, r, "/urls", http.StatusMovedPermanently)
		return
	}
	// if not, show error
	http.Error(w, "error deleteing url 2", http.StatusInternalServerError)
	return
}
