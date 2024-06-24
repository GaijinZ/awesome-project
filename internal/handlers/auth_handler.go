package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"awesomeProject/internal/models"
	"awesomeProject/internal/repositories"
	"awesomeProject/pkg/utils"
)

type Auther interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct {
	authRepository repositories.AuthRepository
}

func NewAuth(authRepository repositories.AuthRepository) Auther {
	return &AuthHandler{
		authRepository: authRepository,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.Auth

	if r.Body == nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		http.Error(w, "Hashing password error", http.StatusBadRequest)
		return
	}

	err = a.authRepository.Register(&user)
	if err != nil {
		http.Error(w, "Register error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var auth models.Auth

	if r.Body == nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := a.authRepository.Login(&auth)
	if err != nil {
		fmt.Errorf("error logging in: %v", err)
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	passwordCheck := utils.CheckPasswordHash(user.Password, auth.Password)
	if !passwordCheck {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(*user)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   false,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusAccepted)
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/api/v1/login", http.StatusSeeOther)
}
