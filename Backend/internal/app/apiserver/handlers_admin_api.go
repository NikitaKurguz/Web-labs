package apiserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"primarte/internal/app/model"
	"primarte/internal/app/validation"
)

type AdminProductJSON struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"       validate:"required,max=100"`
	Type     string   `json:"type"       validate:"required,max=50"`
	Color    string   `json:"color"      validate:"required,max=50"`
	Category []string `json:"category"   validate:"dive,max=50"`
	Price    int      `json:"price"      validate:"gte=0"`
	Img      string   `json:"img"        validate:"max=255"`
}

func (s *APIServer) AdminAPIProductsList(w http.ResponseWriter, r *http.Request) {
	products, err := s.store.GetAllProducts()
	if err != nil {
		log.Printf("Список товаров: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "db error")
		return
	}

	out := make([]AdminProductJSON, 0, len(products))
	for _, p := range products {
		out = append(out, AdminProductJSON{
			ID:       p.ID,
			Name:     p.Name,
			Type:     p.Type,
			Color:    p.Color,
			Category: splitCategory(p.Category),
			Price:    p.Price,
			Img:      p.Img,
		})
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *APIServer) AdminAPIProductsCreate(w http.ResponseWriter, r *http.Request) {
	var in AdminProductJSON
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := validation.Validate(&in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	p := &model.Product{
		Name:     in.Name,
		Type:     in.Type,
		Color:    in.Color,
		Category: joinCategory(in.Category),
		Price:    in.Price,
		Img:      in.Img,
	}
	if err := s.store.CreateProduct(p); err != nil {
		log.Printf("admin api create: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "db error")
		return
	}
	in.ID = p.ID

	adminID := 0
	if u, ok := r.Context().Value(userContextKey).(*model.User); ok {
		adminID = u.Id
	}
	if err := s.store.LogAction(adminID, "product_create", in.Name, clientIP(r)); err != nil {
		log.Printf("admin create: log action: %v", err)
	}
	writeJSON(w, http.StatusOK, in)
}

func (s *APIServer) AdminAPIProductsUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var in AdminProductJSON

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := validation.Validate(&in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	p := &model.Product{
		ID:       id,
		Name:     in.Name,
		Type:     in.Type,
		Color:    in.Color,
		Category: joinCategory(in.Category),
		Price:    in.Price,
		Img:      in.Img,
	}
	if err := s.store.UpdateProduct(p); err != nil {
		log.Printf("admin api update: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "db error")
		return
	}
	in.ID = id

	adminID := 0
	if u, ok := r.Context().Value(userContextKey).(*model.User); ok {
		adminID = u.Id
	}
	if err := s.store.LogAction(adminID, "product_update", in.Name, clientIP(r)); err != nil {
		log.Printf("admin update: log action: %v", err)
	}
	writeJSON(w, http.StatusOK, in)
}

func (s *APIServer) AdminAPIProductsDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.store.DeleteProductById(id); err != nil {
		log.Printf("Удаление товара: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "db error")
		return
	}
	adminID := 0
	if u, ok := r.Context().Value(userContextKey).(*model.User); ok {
		adminID = u.Id
	}
	if err := s.store.LogAction(adminID, "product_delete", strconv.Itoa(id), clientIP(r)); err != nil {
		log.Printf("admin delete: log action: %v", err)
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func joinCategory(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	out := ""
	for i, v := range arr {
		if i > 0 {
			out += ","
		}
		out += v
	}
	return out
}
