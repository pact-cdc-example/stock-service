package stock

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/pact-cdc-example/stock-service/app/product"
	"github.com/pact-cdc-example/stock-service/pkg/cerr"
	"github.com/sirupsen/logrus"
)

type Service interface {
	IsProductAvailableInStockInDesiredQuantity(
		ctx context.Context, req IsProductAvailableInStockRequest,
	) (*IsProductAvailableInStockResponse, error)
	CreateStock(ctx context.Context, req CreateStockRequest) (*GetStockResponse, error)
	ReserveStock(ctx context.Context, req ReserveStockRequest) (*GetStockResponse, error)
}

type service struct {
	logger        *logrus.Logger
	repository    Repository
	productClient product.Client
}

type NewServiceOpts struct {
	L  *logrus.Logger
	R  Repository
	PC product.Client
}

func NewService(opts *NewServiceOpts) Service {
	return &service{
		logger:        opts.L,
		repository:    opts.R,
		productClient: opts.PC,
	}
}

func (s *service) IsProductAvailableInStockInDesiredQuantity(
	ctx context.Context, req IsProductAvailableInStockRequest,
) (*IsProductAvailableInStockResponse, error) {
	stock, err := s.repository.GetStockByProductID(ctx, *req.ProductID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Errorf("error while getting stocks by product id: %v", err)
		return nil, cerr.Processing()
	}

	if stock == nil {
		return nil, cerr.Bag{Code: NoStockInformationFoundAboutGivenProduct,
			Message: "No stock information found for given product id."}
	}

	if (stock.Quantity - stock.ReservedQuantity) < *req.Quantity {
		return &IsProductAvailableInStockResponse{IsAvailable: false}, nil
	}

	return &IsProductAvailableInStockResponse{IsAvailable: true}, nil
}

func (s *service) CreateStock(
	ctx context.Context, req CreateStockRequest) (*GetStockResponse, error) {
	_, err := s.productClient.GetProductByID(ctx, req.ProductID)
	if err != nil {
		s.logger.Errorf("error while getting product by id: %v", err)
		return nil, cerr.Processing()
	}

	stock := &Stock{
		ID:        uuid.New().String(),
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	stock, err = s.repository.CreateStock(ctx, stock)
	if err != nil {
		s.logger.Errorf("error while creating stock: %v", err)
		return nil, cerr.Processing()
	}

	return NewGetStockResponse(stock), nil
}

func (s *service) ReserveStock(
	ctx context.Context, req ReserveStockRequest) (*GetStockResponse, error) {
	stock, err := s.repository.GetStockByProductID(ctx, req.ProductID)
	if err != nil {
		s.logger.Errorf("error while getting stocks by product id: %v", err)
		return nil, cerr.Processing()
	}

	if stock.Quantity < req.Quantity {
		return nil, cerr.Bag{Code: NotEnoughStockToReserve,
			Message: "Not enough stock for given product."}
	}

	stock.ReservedQuantity = stock.ReservedQuantity + req.Quantity

	stock, err = s.repository.UpdateStock(ctx, stock)
	if err != nil {
		s.logger.Errorf("error while updating stock: %v", err)
		return nil, cerr.Processing()
	}

	return NewGetStockResponse(stock), nil
}

//func availableStockQuantity(stocks []Stock) int {
//	availableStock := 0
//	for _, stock := range stocks {
//		availableStock = availableStock + stock.Quantity - stock.ReservedQuantity
//	}
//
//	return availableStock
//}
