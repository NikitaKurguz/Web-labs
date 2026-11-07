package apiserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type ProductJSON struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Color      string   `json:"color"`
	Category   []string `json:"category"`
	Price      string   `json:"price"`
	PriceLabel string   `json:"priceLabel"`
	Img        string   `json:"img"`
}

func (s *APIServer) APIProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := s.store.GetAllProducts()
	if err != nil {
		log.Printf("api products: %v", err)
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}

	out := make([]ProductJSON, 0, len(products))
	for _, p := range products {
		out = append(out, ProductJSON{
			Name:       p.Name,
			Type:       p.Type,
			Color:      p.Color,
			Category:   splitCategory(p.Category),
			Price:      p.GetPriceRange(),
			PriceLabel: p.PriceLabel(),
			Img:        p.Img,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Printf("api products: encode: %v", err)
	}
}

func splitCategory(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
