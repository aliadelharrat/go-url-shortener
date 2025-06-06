package controllers

import (
	"net/http"

	"github.com/aliadelharrat/goshort/models"
)

func (ac *AppController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "id can not be empty", http.StatusBadRequest)
		return
	}
	var url models.ShortURL
	// Todo: check if we need to convert idStr to id
	result := ac.DB.First(&url, id)
	if result.Error != nil {
		http.Error(w, "error fetching url", http.StatusNotFound)
		return
	}
	// delete
	result = ac.DB.Delete(&models.ShortURL{}, id)
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
}