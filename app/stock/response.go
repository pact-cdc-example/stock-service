package stock

import "time"

type GetStockResponse struct {
	ID               string    `json:"id"`
	ProductID        string    `json:"product_id"`
	Quantity         int       `json:"quantity"`
	ReservedQuantity int       `json:"reserved_quantity"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewGetStockResponse(stock *Stock) *GetStockResponse {
	return &GetStockResponse{
		ID:               stock.ID,
		ProductID:        stock.ProductID,
		Quantity:         stock.Quantity,
		ReservedQuantity: stock.ReservedQuantity,
		CreatedAt:        stock.CreatedAt,
		UpdatedAt:        stock.UpdatedAt,
	}
}

type IsProductAvailableInStockResponse struct {
	IsAvailable bool `json:"is_available,omitempty"`
}
