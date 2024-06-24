package handlers

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/repositories"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Categorer interface {
	GetCategoryHandler(w http.ResponseWriter, req *http.Request)
	UpdateCategoryHandler(w http.ResponseWriter, req *http.Request)
	CreateCategoryHandler(w http.ResponseWriter, req *http.Request)
	DeleteCategoryHandler(w http.ResponseWriter, req *http.Request)
}

type CategoryHandler struct {
	categoryRepo repositories.Categorer
}

func NewCategoryHandler(categoryRepo repositories.Categorer) Categorer {
	return &CategoryHandler{
		categoryRepo: categoryRepo,
	}
}

func (c *CategoryHandler) GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categoryID := mux.Vars(r)["category_id"]

	category, err := c.categoryRepo.GetCategory(categoryID)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(category)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (c *CategoryHandler) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		http.Error(w, "Request body missing", http.StatusBadRequest)
		return
	}

	category := &models.Category{}

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	category.ID = mux.Vars(r)["category_id"]

	err = c.categoryRepo.UpdateCategory(*category)
	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *CategoryHandler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		http.Error(w, "Request body missing", http.StatusBadRequest)
		return
	}

	category := &models.Category{}

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if category.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	err = c.categoryRepo.CreateCategory(*category)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *CategoryHandler) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categoryID := mux.Vars(r)["category_id"]

	err := c.categoryRepo.DeleteCategory(categoryID)
	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
