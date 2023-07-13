package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/alekseyshevchenko93/go-crud-api-example/docs"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/config"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/handlers"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/middlewares"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/repository"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

// @title           Swagger Example API
// @version         0.1
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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

	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
