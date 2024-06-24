package repositories

import (
	"awesomeProject/internal/models"
	"awesomeProject/pkg/database"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ProductRepository interface {
	GetProduct(productID string) (*models.ProductResponse, error)
	UpdateProduct(product *models.Product) error
	CreateProduct(product *models.Product) error
	DeleteProduct(id string) error
}

type Product struct {
	db database.Database
}

func NewProduct(db database.Database) ProductRepository {
	return &Product{db: db}
}

func (p *Product) GetProduct(productID string) (*models.ProductResponse, error) {
	product := &models.ProductResponse{}

	err := p.db.QueryRow(GetProduct, productID).
		Scan(&product.Name, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (p *Product) UpdateProduct(product *models.Product) error {
	_, err := p.db.Exec(UpdateProduct, product.ID, product.Name, product.CategoryID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (p *Product) CreateProduct(product *models.Product) error {

	exists, err := checkProductExists(product.Name, p.db)
	if err != nil {
		return fmt.Errorf("failed to check product exist: %w", err)
	}

	if exists {
		return errors.New("product already exists")
	}

	_, err = p.db.Exec(CreateProduct, product.Name, product.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (p *Product) DeleteProduct(id string) error {
	_, err := p.db.Exec(DeleteProduct, id)
	if err != nil {
		return err
	}

	return nil
}

func checkProductExists(categoryName string, db database.Database) (bool, error) {
	var exists bool

	err := db.QueryRow(CheckProductExists, categoryName).Scan(&exists)
	if err != nil {
		return exists, err
	}

	return exists, nil
}
