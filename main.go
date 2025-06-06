package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/aliadelharrat/go-url-shortener/controllers"
	"github.com/aliadelharrat/go-url-shortener/models"
	tpl "github.com/aliadelharrat/go-url-shortener/templates"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var templateCache map[string]*template.Template

func main() {
	r := chi.NewRouter()

	db, err := gorm.Open(sqlite.Open("shorts.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.ShortURL{})

	templateCache, err = tpl.NewTemplateCache("./templates/")
	if err != nil {
		log.Fatalf("Could not create template cache: %v", err)
	}

	appController := controllers.NewAppController(db, templateCache)

	r.Get("/", appController.HomeHandler)
	r.Post("/submit", appController.SubmitHandler)
	r.Get("/urls", appController.URLsHandler)
	r.Post("/delete", appController.DeleteHandler)
	r.Get("/{slug}", appController.RedirectHandler)
	r.Get("/url/{slug}", appController.GetUrlHandler)

	log.Println("running server on port :8080")
	if err = http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("error running server: ", err)
	}
}