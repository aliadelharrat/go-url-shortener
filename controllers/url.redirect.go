package controllers

import (
	"fmt"
	"net/http"

	"github.com/aliadelharrat/goshort/models"
	"github.com/go-chi/chi/v5"
)

func (ac *AppController) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	var url models.ShortURL
	result := ac.DB.Where("surl = ?", slug).First(&url)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("error fetching short url: %v", result.Error), http.StatusNotFound)
		return
	}
	// Increment number of visits
	url.Visits++
	ac.DB.Save(&url)
	http.Redirect(w, r, "//"+url.URL, http.StatusFound)
}