package persistence

import (
	"context"
	"database/sql"
	"github.com/pact-cdc-example/stock-service/app/stock"
	"github.com/sirupsen/logrus"
	"time"
)

//go:generate mockgen -source=postgres.go -destination=mock_postgres_repository.go -package=persistence
type PostgresRepository interface {
	GetStockByProductID(ctx context.Context, productID string) (*stock.Stock, error)
	CreateStock(ctx context.Context, stock *stock.Stock) (*stock.Stock, error)
	UpdateStock(
		ctx context.Context, stock *stock.Stock) (*stock.Stock, error)
}

type postgresRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

type NewPostgresRepositoryOpts struct {
	DB *sql.DB
	L  *logrus.Logger
}

func NewPostgresRepository(opts *NewPostgresRepositoryOpts) PostgresRepository {
	return &postgresRepository{
		db:     opts.DB,
		logger: opts.L,
	}
}

func (p *postgresRepository) GetStockByProductID(
	ctx context.Context, productID string) (*stock.Stock, error) {
	rows := p.db.QueryRowContext(ctx, `
		SELECT id, product_id, quantity, reserved_quantity, created_at, updated_at
		FROM stocks
		WHERE product_id = $1
		LIMIT 1;
	`, productID)

	var stock stock.Stock
	if err := rows.Scan(
		&stock.ID,
		&stock.ProductID,
		&stock.Quantity,
		&stock.ReservedQuantity,
		&stock.CreatedAt,
		&stock.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &stock, nil
}

func (p *postgresRepository) CreateStock(
	ctx context.Context, stock *stock.Stock) (*stock.Stock, error) {
	row := p.db.QueryRowContext(ctx, `
		INSERT INTO stocks (id, product_id, quantity, reserved_quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at;`,
		stock.ID,
		stock.ProductID,
		stock.Quantity,
		stock.ReservedQuantity,
	)

	var createdAt time.Time
	var updatedAt time.Time
	if err := row.Scan(&createdAt, &updatedAt); err != nil {
		return nil, err
	}

	stock.CreatedAt = createdAt
	stock.UpdatedAt = updatedAt

	return stock, nil
}

func (p *postgresRepository) UpdateStock(
	ctx context.Context, stock *stock.Stock) (*stock.Stock, error) {
	row := p.db.QueryRowContext(ctx, `
		UPDATE stocks
		SET quantity = $1, reserved_quantity = $2
		WHERE id = $3
		RETURNING updated_at;`,
		stock.Quantity,
		stock.ReservedQuantity,
		stock.ID,
	)

	var updatedAt time.Time
	if err := row.Scan(&updatedAt); err != nil {
		return nil, err
	}
	stock.UpdatedAt = updatedAt

	return stock, nil
}

func createStockTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS stocks (
			id VARCHAR(255) NOT NULL UNIQUE,
			product_id VARCHAR(255) NOT NULL,
			quantity INT NOT NULL,
			reserved_quantity INT NOT NULL DEFAULT 0,
			place VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		panic(err)
	}
}
