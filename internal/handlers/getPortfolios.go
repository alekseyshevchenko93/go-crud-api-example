package handlers

import (
	"net/http"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/models"
	"github.com/labstack/echo/v4"
)

type PortfoliosGetter interface {
	GetPortfolios() ([]models.Portfolio, error)
}

func NewGetPortfoliosHandler(portfolioService PortfoliosGetter) func(echo.Context) error {
	return func(ctx echo.Context) error {
		portfolios, err := portfolioService.GetPortfolios()

		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, portfolios)
	}
}
