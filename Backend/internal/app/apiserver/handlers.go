package apiserver

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"primarte/internal/app/model"
	"primarte/internal/app/validation"
)

type RegisterInput struct {
	Login    string `validate:"required,alphanum,min=3,max=30"`
	Password string `validate:"required,min=6"`
	Fullname string `validate:"max=100"`
	Email    string `validate:"omitempty,email,max=100"`
}

type UpdateProfileInput struct {
	Fullname    string `validate:"max=100"`
	Email       string `validate:"omitempty,email,max=100"`
	NewPassword string `validate:"omitempty,min=6"`
}

func (s *APIServer) renderRegister(w http.ResponseWriter, errMsg string) {
	data := map[string]any{"Error": errMsg}
	if err := s.registerTmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Printf("template register: %v", err)
		http.Error(w, "ошибка шаблона", http.StatusInternalServerError)
	}
}

func saveAvatar(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		return "", errors.New("недопустимый формат файла")
	}

	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join("uploads", "avatars", name)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/" + filepath.ToSlash(path), nil
}

func (s *APIServer) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		s.renderRegister(w, "Ошибка загрузки формы")
		return
	}

	login := strings.TrimSpace(r.FormValue("login"))
	password := r.FormValue("password")
	fullname := strings.TrimSpace(r.FormValue("fullname"))
	email := strings.TrimSpace(r.FormValue("email"))

	input := RegisterInput{
		Login:    login,
		Password: password,
		Fullname: fullname,
		Email:    email,
	}
	if err := validation.Validate(&input); err != nil {
		s.renderRegister(w, err.Error())
		return
	}

	if _, err := s.store.GetUserByLogin(login); err == nil {
		s.renderRegister(w, "Пользователь с таким логином уже существует")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("register: bcrypt: %v", err)
		s.renderRegister(w, "Ошибка сервера")
		return
	}

	avatarPath := ""
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		avatarPath, err = saveAvatar(file, header)
		if err != nil {
			log.Printf("register: save avatar: %v", err)
			s.renderRegister(w, "Ошибка загрузки аватара")
			return
		}
	}

	u := &model.User{
		Login:    login,
		Password: string(hash),
		FullName: fullname,
		Email:    email,
		Avatar:   avatarPath,
	}
	if err := s.store.CreateUser(u); err != nil {
		log.Printf("register: create user: %v", err)
		s.renderRegister(w, "Ошибка сервера")
		return
	}
	if err := s.store.LogAction(u.Id, "register", u.Login, clientIP(r)); err != nil {
		log.Printf("register: log action: %v", err)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *APIServer) RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.registerTmpl.ExecuteTemplate(w, "layout.html", nil); err != nil {
		log.Printf("template register: %v", err)
		http.Error(w, "ошибка шаблона", http.StatusInternalServerError)
	}
}

func (s *APIServer) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.loginTmpl.ExecuteTemplate(w, "layout.html", nil); err != nil {
		log.Printf("template login: %v", err)
		http.Error(w, "ошибка шаблона", http.StatusInternalServerError)
	}
}

func (s *APIServer) renderLogin(w http.ResponseWriter, errMsg string) {
	data := map[string]any{"Error": errMsg}
	if err := s.loginTmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Printf("template login: %v", err)
		http.Error(w, "ошибка шаблона", http.StatusInternalServerError)
	}
}

func generateSessionHash() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *APIServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderLogin(w, "Ошибка формы")
		return
	}

	login := strings.TrimSpace(r.FormValue("login"))
	password := r.FormValue("password")

	if login == "" || password == "" {
		s.renderLogin(w, "Заполните логин и пароль")
		return
	}

	u, err := s.store.GetUserByLogin(login)
	if err != nil {
		log.Printf("login: get user: %v", err)
		s.renderLogin(w, "Неверный логин или пароль")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		s.renderLogin(w, "Неверный логин или пароль")
		return
	}

	sessionHash, err := generateSessionHash()
	if err != nil {
		log.Printf("login: generate hash: %v", err)
		s.renderLogin(w, "Ошибка сервера")
		return
	}

	if err := s.store.UpdateHash(u.Id, sessionHash); err != nil {
		log.Printf("login: update hash: %v", err)
		s.renderLogin(w, "Ошибка сервера")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(u.Id),
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_hash",
		Value:    sessionHash,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	if err := s.store.LogAction(u.Id, "login", u.Login, clientIP(r)); err != nil {
		log.Printf("login: log action: %v", err)
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (s *APIServer) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if idCookie, err := r.Cookie("user_id"); err == nil {
		id, err := strconv.Atoi(idCookie.Value)
		if err == nil {
			if err := s.store.UpdateHash(id, ""); err != nil {
				log.Printf("logout: update hash: %v", err)
			}
			if err := s.store.LogAction(id, "logout", "", clientIP(r)); err != nil {
				log.Printf("logout: log action: %v", err)
			}
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_hash",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *APIServer) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value(userContextKey).(*model.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]any{
		"User": u,
	}

	if err := s.profileTmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Printf("template profile: %v", err)
		http.Error(w, "ошибка шаблона", http.StatusInternalServerError)
	}
}

func (s *APIServer) renderProfile(w http.ResponseWriter, u *model.User, success, errMsg string) {
	data := map[string]any{
		"User":    u,
		"Success": success,
		"Error":   errMsg,
	}
	if err := s.profileTmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Printf("template profile: %v", err)
		http.Error(w, "ошибка шаблона", http.StatusInternalServerError)
	}
}

func (s *APIServer) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value(userContextKey).(*model.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		s.renderProfile(w, u, "", "Ошибка загрузки формы")
		return
	}

	fullname := strings.TrimSpace(r.FormValue("fullname"))
	email := strings.TrimSpace(r.FormValue("email"))
	newPassword := r.FormValue("new_password")

	input := UpdateProfileInput{
		Fullname:    fullname,
		Email:       email,
		NewPassword: newPassword,
	}
	if err := validation.Validate(&input); err != nil {
		s.renderProfile(w, u, "", err.Error())
		return
	}

	u.FullName = fullname
	u.Email = email

	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		avatarPath, err := saveAvatar(file, header)
		if err != nil {
			log.Printf("update profile: save avatar: %v", err)
			s.renderProfile(w, u, "", "Ошибка загрузки аватара")
			return
		}
		u.Avatar = avatarPath
	}

	if err := s.store.UpdateUser(u); err != nil {
		log.Printf("update profile: %v", err)
		s.renderProfile(w, u, "", "Ошибка сохранения")
		return
	}

	if newPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("update profile: bcrypt: %v", err)
			s.renderProfile(w, u, "", "Ошибка сервера")
			return
		}
		if err := s.store.UpdatePassword(u.Id, string(hash)); err != nil {
			log.Printf("update profile: update password: %v", err)
			s.renderProfile(w, u, "", "Ошибка сохранения пароля")
			return
		}
	}
	if err := s.store.LogAction(u.Id, "profile_update", "", clientIP(r)); err != nil {
		log.Printf("update profile: log action: %v", err)
	}
	s.renderProfile(w, u, "Данные сохранены", "")
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.Index(xff, ","); i != -1 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i != -1 {
		host = host[:i]
	}
	return host
}
