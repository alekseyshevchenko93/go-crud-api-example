package handlers

import (
	"net/http"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/labstack/echo/v4"
)

func NewDeletePortfolioHandler(portfolioService services.PortfolioService) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.Param("id")

		err := portfolioService.DeletePortfolio(id)

		if err != nil {
			return err
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
