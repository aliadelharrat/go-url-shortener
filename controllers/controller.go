package controllers

import (
	"html/template"
	"net/http"

	tpl "github.com/aliadelharrat/goshort/templates"
	"gorm.io/gorm"
)

type AppController struct {
	DB *gorm.DB
	TemplateCache map[string]*template.Template
}

func NewAppController (db *gorm.DB, tc map[string]*template.Template) *AppController {
	return &AppController{
		DB: db,
		TemplateCache: tc,
	}
}

func (ac *AppController) render (w http.ResponseWriter, name string, data any) {
	tpl.Render(w, name, ac.TemplateCache, data)
}