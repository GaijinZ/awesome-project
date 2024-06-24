package repositories

import (
	"awesomeProject/internal/models"
	"awesomeProject/pkg/database"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Categorer interface {
	GetCategory(categoryID string) (*models.CategoryResponse, error)
	UpdateCategory(category models.Category) error
	CreateCategory(category models.Category) error
	DeleteCategory(categoryID string) error
}

type Category struct {
	db database.Database
}

func NewCategory(db database.Database) Categorer {
	return &Category{db: db}
}

func (c *Category) GetCategory(categoryID string) (*models.CategoryResponse, error) {
	category := &models.CategoryResponse{}

	err := c.db.QueryRow(GetCategoryByID, categoryID).
		Scan(&category.Name, &category.ProductID, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category not found: %w", err)
		}

		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

func (c *Category) UpdateCategory(category models.Category) error {
	_, err := c.db.Exec(UpdateCategory, category.ID, category.Name, category.ProductID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	return nil
}

func (c *Category) CreateCategory(category models.Category) error {
	exists, err := checkCategoryExists(category.Name, c.db)
	if err != nil {
		return fmt.Errorf("failed to check category exist: %w", err)
	}

	if exists {
		return errors.New("category already exists")
	}

	_, err = c.db.Exec(CreateCategory, category.Name, category.ProductID)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

func (c *Category) DeleteCategory(categoryID string) error {
	_, err := c.db.Exec(DeleteCategory, categoryID)
	if err != nil {
		return err
	}

	return nil
}

func checkCategoryExists(categoryName string, db database.Database) (bool, error) {
	var exists bool

	err := db.QueryRow(CheckCategoryExists, categoryName).Scan(&exists)
	if err != nil {
		return exists, err
	}

	return exists, nil
}
