package controllers

import (
	"fmt"
	"net/http"

	"github.com/aliadelharrat/goshort/models"
	tpl "github.com/aliadelharrat/goshort/templates"
)

func (ac *AppController) URLsHandler(w http.ResponseWriter, r *http.Request) {
	var urls []models.ShortURL
	result := ac.DB.Find(&urls)
	if result.Error != nil {
		fmt.Fprint(w, "error fetching urls")
	}

	tpl.Render(w, "urls.page.gohtml", ac.TemplateCache, struct {
		URLS []models.ShortURL
	}{
		URLS: urls,
	})
}