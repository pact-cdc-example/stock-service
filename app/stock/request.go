package stock

import (
	"github.com/pact-cdc-example/stock-service/pkg/cerr"
)

type IsProductAvailableInStockRequest struct {
	ProductID *string `json:"product_id,omitempty"`
	Quantity  *int    `json:"quantity,omitempty"`
}

func (i IsProductAvailableInStockRequest) Validate() error {
	if i.ProductID == nil {
		return cerr.Bag{Code: ProductIDMustBeGivenToStockInquiry, Message: "Product id must be given to stock inquiry."}
	}

	if i.Quantity == nil {
		return cerr.Bag{Code: QuantityMustBeGivenToStockInquiry, Message: "Quantity must be given to stock inquiry."}
	}

	return nil
}

type CreateStockRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ReserveStockRequest struct {
	ProductID string `json:"product_id,omitempty"`
	Quantity  int    `json:"quantity,omitempty"`
}

func (r ReserveStockRequest) Validate() error {
	if r.ProductID == "" {
		return cerr.Bag{Code: ProductIDMustBeGivenToReserveStock, Message: "Product id must be given to reserve stock."}
	}

	if r.Quantity == 0 {
		return cerr.Bag{Code: QuantityMustBeGivenToReserveStock, Message: "Quantity must be given to reserve stock."}
	}

	return nil
}
