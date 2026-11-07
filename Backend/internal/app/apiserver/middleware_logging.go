package apiserver

import (
	"log"
	"net/http"

	"primarte/internal/app/model"
)

func (s *APIServer) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		path := r.URL.Path
		if isStaticOrAPI(path) {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)

		userID := 0
		if u, ok := r.Context().Value(userContextKey).(*model.User); ok {
			userID = u.Id
		}
		if err := s.store.LogAction(userID, "page_view", path, clientIP(r)); err != nil {
			log.Printf("logging middleware: %v", err)
		}
	})
}

func isStaticOrAPI(path string) bool {
	prefixes := []string{
		"/assets/", "/css/", "/js/", "/uploads/",
		"/api/", "/partials/",
	}
	for _, p := range prefixes {
		if len(path) >= len(p) && path[:len(p)] == p {
			return true
		}
	}
	return false
}
