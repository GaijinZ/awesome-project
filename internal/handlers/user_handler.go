package handlers

import (
	"encoding/json"
	"net/http"

	"awesomeProject/internal/models"
	"awesomeProject/internal/repositories"
	"awesomeProject/pkg/utils"

	"github.com/gorilla/mux"
)

type Userer interface {
	GetUserByUsername(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	userRepository repositories.UserRepository
}

func NewUserHandler(userRepository repositories.UserRepository) Userer {
	return &UserHandler{
		userRepository: userRepository,
	}
}

func (u *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := mux.Vars(r)["username"]

	userResponse, err := u.userRepository.GetUserByUsername(name)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(userResponse)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (u *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := u.userRepository.GetAllUsers()
	if err != nil {
		http.Error(w, "Users not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		http.Error(w, "Request body missing", http.StatusBadRequest)
		return
	}

	user := &models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to encode request", http.StatusBadRequest)
		return
	}

	userID := mux.Vars(r)["user_id"]

	cookieUserID, err := r.Cookie("userID")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		return
	}

	if cookieUserID.Value != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, err := r.Cookie("username")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		return
	}

	userResponse, err := u.userRepository.GetUserByUsername(username.Value)
	if err != nil || userResponse == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.ID = cookieUserID.Value

	user.Username = defaultIfEmpty(user.Username, userResponse.Username)
	user.Email = defaultIfEmpty(user.Email, userResponse.Email)
	user.Role = defaultIfEmpty(user.Role, userResponse.Role)

	err = u.userRepository.UpdateUser(user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		http.Error(w, "Request body missing", http.StatusBadRequest)
		return
	}

	user := &models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to encode request", http.StatusBadRequest)
		return
	}

	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		http.Error(w, "Hashing password error", http.StatusBadRequest)
		return
	}

	err = u.userRepository.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := mux.Vars(r)["user_id"]

	err := u.userRepository.DeleteUser(userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func defaultIfEmpty(input, defaultValue string) string {
	if input == "" {
		return defaultValue
	}

	return input
}
