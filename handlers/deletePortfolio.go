package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PortfolioDeleter interface {
	DeletePortfolio(string) error
}

func NewDeletePortfolioHandler(portfolioService PortfolioDeleter) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.Param("id")

		err := portfolioService.DeletePortfolio(id)

		if err != nil {
			return err
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
