package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aliadelharrat/goshort/models"
	"github.com/teris-io/shortid"
)

func (ac *AppController) SubmitHandler(w http.ResponseWriter, r *http.Request) {
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

		shortURLValue, err := shortid.Generate()
		if err != nil {
			log.Printf("Error generating random uuid: %v", err)
			http.Error(w, "Failed to generate short url", http.StatusInternalServerError)
			return
		}
		url := models.ShortURL{
			URL:    urlValue,
			SURL:   shortURLValue,
			Visits: 0,
		}
		result := ac.DB.Create(&url)
		if result.Error != nil {
			log.Printf("Error creating short url: %v", result.Error)
			http.Error(w, "Failed to create short url", http.StatusInternalServerError)
			return
		}
		link := fmt.Sprintf("/url/%s", url.SURL)
		http.Redirect(w, r, link, http.StatusMovedPermanently)
		return
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}