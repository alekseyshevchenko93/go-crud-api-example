package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alekseyshevchenko93/go-crud-api-example/config"
	"github.com/alekseyshevchenko93/go-crud-api-example/handlers"
	"github.com/alekseyshevchenko93/go-crud-api-example/middlewares"
	"github.com/alekseyshevchenko93/go-crud-api-example/repository"
	"github.com/alekseyshevchenko93/go-crud-api-example/services"
	"github.com/labstack/echo/v4"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	appConfig, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.HTTPErrorHandler = middlewares.ErrorHandler

	portfolioRepository := repository.NewPortfolioRepository()
	portfolioService := services.NewPortfolioService(portfolioRepository)

	e.GET("/portfolios/:id", handlers.NewGetPortfolioByIdHandler(portfolioService))
	e.GET("/portfolios", handlers.NewGetPortfoliosHandler(portfolioService))
	e.PUT("/portfolios/:id", handlers.NewUpdatePortfolioHandler(portfolioService))
	e.POST("/portfolios", handlers.NewCreatePortfolioHandler(portfolioService))
	e.DELETE("/portfolios/:id", handlers.NewDeletePortfolioHandler(portfolioService))

	go func() {
		address := fmt.Sprintf(":%s", appConfig.HttpPort)

		if err := e.Start(address); err != nil {
			panic(err)
		}
	}()

	select {
	case <-ctx.Done():
		os.Exit(0)
	}
}
