package stock

import "context"

//go:generate mockgen -source=repository.go -destination=mock_repository.go -package=stock
type Repository interface {
	GetStockByProductID(ctx context.Context, productID string) (*Stock, error)
	CreateStock(ctx context.Context, stock *Stock) (*Stock, error)
	UpdateStock(ctx context.Context, stock *Stock) (*Stock, error)
}
