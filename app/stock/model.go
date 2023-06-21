package stock

import "time"

type Stock struct {
	ID               string    `json:"-"`
	ProductID        string    `json:"-"`
	Quantity         int       `json:"-"`
	ReservedQuantity int       `json:"-"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}
