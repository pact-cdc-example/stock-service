package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server interface {
	Run() error
}

type NewServerOpts struct {
	Port string
}

type server struct {
	app  *fiber.App
	opts *NewServerOpts
}

type RouteHandler interface {
	SetupRoutes(fr fiber.Router)
}

func New(opts *NewServerOpts, routeHandlers []RouteHandler) Server {
	app := fiber.New()

	app.Use(cors.New())

	apiGroup := app.Group("/api")
	v1Group := apiGroup.Group("/v1")

	for _, handler := range routeHandlers {
		handler.SetupRoutes(v1Group)
	}

	s := &server{app: app, opts: opts}

	s.addHealthCheckRoutes()

	return s
}

func (s *server) addHealthCheckRoutes() {
	s.app.Get("/liveness", liveness)
	s.app.Get("/readines", readiness)
}

func liveness(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func readiness(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (s *server) Run() error {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdownChan
		err := s.app.Shutdown()
		if err != nil {
			log.Println("Error on shutdown gracefully")
		}
	}()

	return s.app.Listen(fmt.Sprintf(":%s", s.opts.Port))
}
