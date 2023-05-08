package handlers

import (
	"net/http"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/labstack/echo/v4"
)

func NewGetPortfoliosHandler(portfolioService services.PortfolioService) func(echo.Context) error {
	return func(ctx echo.Context) error {
		portfolios, err := portfolioService.GetPortfolios()

		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, portfolios)
	}
}
