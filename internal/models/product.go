package models

import "time"

type Product struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name"`
	CategoryID int       `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProductResponse struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name"`
	CategoryID int       `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
