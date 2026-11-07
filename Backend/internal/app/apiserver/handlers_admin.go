package apiserver

import (
	"html/template"
	"log"
	"net/http"
)

func (s *APIServer) AdminProductsPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/admin/layout.html",
		"templates/admin/products.html",
	)
	if err != nil {
		log.Printf("admin products template: %v", err)
		http.Error(w, "ошибка шаблона", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "layout.html", nil); err != nil {
		log.Printf("admin products render: %v", err)
		http.Error(w, "ошибка шаблона", http.StatusInternalServerError)
	}
}

func (s *APIServer) AdminStatsPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/admin/layout.html",
		"templates/admin/stats.html",
	)
	if err != nil {
		log.Printf("admin stats template: %v", err)
		http.Error(w, "ошибка шаблона", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "layout.html", nil); err != nil {
		log.Printf("admin stats render: %v", err)
		http.Error(w, "ошибка шаблона", http.StatusInternalServerError)
	}
}
