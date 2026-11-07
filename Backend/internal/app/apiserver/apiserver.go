package apiserver

import (
	"html/template"
	"net/http"
	"primarte/internal/app/store"
)

type APIServer struct {
	config *Config
	store  *store.Store

	loginTmpl    *template.Template
	registerTmpl *template.Template
	profileTmpl  *template.Template
}

func New(config *Config, st *store.Store) *APIServer {
	return &APIServer{
		config: config,
		store:  st,
	}
}
func (s *APIServer) Start() error {
	s.loginTmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/login.html"))
	s.registerTmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/register.html"))
	s.profileTmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/profile.html"))

	mux := http.NewServeMux()

	mux.Handle("/assets/", http.FileServer(http.Dir("../L1-site")))
	mux.Handle("/css/", http.FileServer(http.Dir("../L1-site")))
	mux.Handle("/js/", http.FileServer(http.Dir("../L1-site")))

	mux.Handle("/uploads/", http.FileServer(http.Dir(".")))

	mux.HandleFunc("GET /register", s.RegisterPageHandler)
	mux.HandleFunc("POST /register", s.RegisterHandler)
	mux.HandleFunc("GET /login", s.LoginPageHandler)
	mux.HandleFunc("POST /login", s.LoginHandler)
	mux.HandleFunc("GET /logout", s.LogoutHandler)
	mux.Handle("GET /profile", s.AuthMiddleware(http.HandlerFunc(s.ProfileHandler)))
	mux.Handle("POST /profile", s.AuthMiddleware(http.HandlerFunc(s.UpdateProfileHandler)))
	mux.HandleFunc("GET /api/products", s.APIProductsHandler)
	mux.Handle("GET /api/admin/products",
		s.AuthMiddleware(s.AdminMiddleware(http.HandlerFunc(s.AdminAPIProductsList))))
	mux.Handle("POST /api/admin/products",
		s.AuthMiddleware(s.AdminMiddleware(http.HandlerFunc(s.AdminAPIProductsCreate))))
	mux.Handle("PUT /api/admin/products/{id}",
		s.AuthMiddleware(s.AdminMiddleware(http.HandlerFunc(s.AdminAPIProductsUpdate))))
	mux.Handle("DELETE /api/admin/products/{id}",
		s.AuthMiddleware(s.AdminMiddleware(http.HandlerFunc(s.AdminAPIProductsDelete))))
	mux.Handle("GET /admin/products",
		s.AuthMiddleware(s.AdminMiddleware(http.HandlerFunc(s.AdminProductsPageHandler))))
	mux.Handle("GET /admin/stats",
		s.AuthMiddleware(s.AdminMiddleware(http.HandlerFunc(s.AdminStatsPageHandler))))
	mux.Handle("GET /api/admin/actions",
		s.AuthMiddleware(s.AdminMiddleware(http.HandlerFunc(s.AdminAPIActions))))
	mux.Handle("GET /api/admin/stats",
		s.AuthMiddleware(s.AdminMiddleware(http.HandlerFunc(s.AdminAPIStats))))

	mux.Handle("/", http.FileServer(http.Dir("../L1-site")))

	return http.ListenAndServe(s.config.Server.BindAddr,
		s.LoggingMiddleware(mux))
}
