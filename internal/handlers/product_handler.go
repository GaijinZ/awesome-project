package handlers

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/repositories"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Producter interface {
	GetProductHandler(w http.ResponseWriter, req *http.Request)
	UpdateProductHandler(w http.ResponseWriter, req *http.Request)
	CreateProductHandler(w http.ResponseWriter, req *http.Request)
	DeleteProductHandler(w http.ResponseWriter, req *http.Request)
}

type ProductHandler struct {
	product repositories.ProductRepository
}

func NewProductHandler(product repositories.ProductRepository) Producter {
	return &ProductHandler{
		product: product,
	}
}

func (p *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productID := mux.Vars(r)["product_id"]

	product, err := p.product.GetProduct(productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (p *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		http.Error(w, "Request body missing", http.StatusBadRequest)
		return
	}

	var product *models.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	productID := mux.Vars(r)["product_id"]

	product.ID = productID

	err = p.product.UpdateProduct(product)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		http.Error(w, "Request body missing", http.StatusBadRequest)
		return
	}

	var product *models.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if product.Name == "" {
		http.Error(w, "Product name must be provided", http.StatusBadRequest)
		return
	}

	err = p.product.CreateProduct(product)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productID := mux.Vars(r)["product_id"]

	err := p.product.DeleteProduct(productID)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
