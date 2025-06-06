package controllers

import "net/http"

func (ac *AppController) HomeHandler(w http.ResponseWriter, _ *http.Request) {
	ac.render(w, "home.page.gohtml", nil)
}