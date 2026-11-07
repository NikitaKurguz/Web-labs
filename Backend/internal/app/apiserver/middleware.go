package apiserver

import (
	"context"
	"log"
	"net/http"
	"primarte/internal/app/model"
	"strconv"
)

type contextKey string

const userContextKey contextKey = "user"

func (s *APIServer) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idCookie, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		hashCookie, err := r.Cookie("user_hash")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		id, err := strconv.Atoi(idCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		u, err := s.store.GetUserById(id)
		if err != nil {
			log.Printf("auth: get user: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if u.Hash != hashCookie.Value {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *APIServer) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(userContextKey).(*model.User)
		if !ok || !u.IsAdmin {
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
