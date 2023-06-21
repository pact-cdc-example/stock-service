package main

import (
	"github.com/pact-cdc-example/stock-service/app/persistence"
	"github.com/pact-cdc-example/stock-service/app/stock"
	"github.com/pact-cdc-example/stock-service/config"
	"github.com/pact-cdc-example/stock-service/pkg/postgres"
	"github.com/pact-cdc-example/stock-service/pkg/server"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	c := config.New()

	db := postgres.New(&postgres.NewPostgresOpts{
		Host:     c.Postgres().Host,
		Port:     c.Postgres().Port,
		DBName:   c.Postgres().DBName,
		Password: c.Postgres().Password,
		Username: c.Postgres().Username,
	})

	logger := logrus.New()

	stockRepository := persistence.NewPostgresRepository(&persistence.NewPostgresRepositoryOpts{
		DB: db,
		L:  logger,
	})

	stockService := stock.NewService(&stock.NewServiceOpts{
		R: stockRepository,
		L: logger,
	})

	stockHandler := stock.NewHandler(&stock.NewHandlerOpts{
		S: stockService,
		L: logger,
	})

	app := server.New(&server.NewServerOpts{
		Port: c.Server().Port,
	}, []server.RouteHandler{
		stockHandler,
	})

	if err := app.Run(); err != nil {
		log.Fatalf("server is closed: %v", err)
	}
}
