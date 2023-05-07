package handlers

import (
	"net/http"

	"github.com/alekseyshevchenko93/go-crud-api-example/domain/models"
	"github.com/labstack/echo/v4"
)

type PortfolioByIdGetter interface {
	GetPortfolioById(string) (models.Portfolio, error)
}

func NewGetPortfolioByIdHandler(portfolioService PortfolioByIdGetter) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.Param("id")
		portfolio, err := portfolioService.GetPortfolioById(id)

		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, portfolio)
	}
}
