package stock

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pact-cdc-example/stock-service/pkg/cerr"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=manager.go -destination=mock_manager.go -package=config
type Handler interface {
	SetupRoutes(fr fiber.Router)
	IsProductAvailableInStockInDesiredQuantity(ctx *fiber.Ctx) error
	CreateStock(ctx *fiber.Ctx) error
	ReserveStock(ctx *fiber.Ctx) error
}

type handler struct {
	logger  *logrus.Logger
	service Service
}

type NewHandlerOpts struct {
	L *logrus.Logger
	S Service
}

func NewHandler(opts *NewHandlerOpts) Handler {
	return &handler{
		logger:  opts.L,
		service: opts.S,
	}
}

func (h *handler) IsProductAvailableInStockInDesiredQuantity(c *fiber.Ctx) error {
	ctx := c.Context()

	var req IsProductAvailableInStockRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(cerr.BodyParser())
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	isAvailableResp, err := h.service.IsProductAvailableInStockInDesiredQuantity(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return c.Status(fiber.StatusOK).JSON(isAvailableResp)
}

func (h *handler) CreateStock(c *fiber.Ctx) error {
	ctx := c.Context()

	var req CreateStockRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(cerr.BodyParser())
	}

	stocks, err := h.service.CreateStock(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(stocks)
}

func (h *handler) ReserveStock(c *fiber.Ctx) error {
	ctx := c.Context()

	var req ReserveStockRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(cerr.BodyParser())
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	stocks, err := h.service.ReserveStock(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(stocks)
}

func (h *handler) SetupRoutes(fr fiber.Router) {
	stocksGroup := fr.Group("/stocks")
	stocksGroup.Post("/availability", h.IsProductAvailableInStockInDesiredQuantity)
	stocksGroup.Post("/", h.CreateStock)
	stocksGroup.Put("/reserve", h.ReserveStock)
}
