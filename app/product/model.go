package product

import "time"

type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Color     string    `json:"color,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Price     float64   `json:"price"`
	ImageURL  string    `json:"image_url,omitempty"`
	Type      string    `json:"type"`
}

type GetProductByIDsRequest struct {
	IDs []string `json:"ids"`
}

type GetProductResponse struct {
	Product Product `json:"product"`
}

type GetProductsResponse struct {
	Products []Product `json:"products"`
}
