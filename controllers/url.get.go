package controllers

import (
	"net/http"

	"github.com/aliadelharrat/goshort/models"
	tpl "github.com/aliadelharrat/goshort/templates"
	"github.com/go-chi/chi/v5"
)

func (ac *AppController) GetUrlHandler(w http.ResponseWriter, r *http.Request) {
	// we need to show, succees message if request is from post request
	// if its a get request, show another message
	slug := chi.URLParam(r, "slug")

	var url models.ShortURL
	// Check if url exist inside database
	result := ac.DB.Where("surl = ?", slug).First(&url)
	if result.Error != nil {
		http.Error(w, "url not found", http.StatusNotFound)
		return
	}

	tpl.Render(w, "url.page.gohtml", ac.TemplateCache, struct {
		Title string
		Link  string
		URL   models.ShortURL
	}{
		Title: "Copy this URL and share it with your friends!",
		Link:  r.Host + "/" + url.SURL,
		URL:   url,
	})
}