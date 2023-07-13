package handlers

import (
	"net/http"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/labstack/echo/v4"
)

// DeletePortfolioById deletes portfolio
// @Summary      Deletes portfolio by id
// @Description  Deletes portfolio
// @Tags         Portfolios
// @Produce      json
// @Success      200
// @Param        id path int  true "Portfolio ID"
// @Router       /portfolios/{id} [delete]
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
