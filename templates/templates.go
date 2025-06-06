package tpl

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.gohtml"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.gohtml"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

func Render (w http.ResponseWriter, name string, templateCache map[string]*template.Template, data any) {
	ts, ok := templateCache[name]
	if !ok {
		http.Error(w, "The template does not exist.", http.StatusInternalServerError)
		return
	}

	err := ts.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}