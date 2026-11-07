package apiserver

import (
	"log"
	"net/http"
	"strconv"

	"primarte/internal/app/store"
)

func (s *APIServer) AdminAPIActions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	action := q.Get("action")

	sortBy := q.Get("sort")
	if sortBy == "" {
		sortBy = "created_at"
	}

	order := q.Get("order")
	if order == "" {
		order = "desc"
	}

	limit := 200
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	rows, err := s.store.GetActions(action, sortBy, order, limit)
	if err != nil {
		log.Printf("admin actions: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "db error")
		return
	}
	if rows == nil {
		rows = []store.ActionRow{}
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *APIServer) AdminAPIStats(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	groupBy := q.Get("group")
	if groupBy == "" {
		groupBy = "day"
	}
	action := q.Get("action")
	sortBy := q.Get("sort")
	if sortBy == "" {
		sortBy = "period"
	}
	order := q.Get("order")
	if order == "" {
		order = "desc"
	}
	limit := 0
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	rows, err := s.store.GetStats(groupBy, action, sortBy, order, limit)
	if err != nil {
		log.Printf("admin stats: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "db error")
		return
	}
	if rows == nil {
		rows = []store.StatsRow{}
	}
	writeJSON(w, http.StatusOK, rows)
}
