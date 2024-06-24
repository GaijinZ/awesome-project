package models

import "time"

type Category struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	ProductID int       `json:"product_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type CategoryResponse struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	ProductID int       `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
