package handlers

import (
	"net/http"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/labstack/echo/v4"
)

// GetPortfolioById responds with the list of all portfolios
// @Summary      Gets portfolio by id
// @Tags         Portfolios
// @Produce      json
// @Success      200  {object}  models.Portfolio
// @Param        id path int  true "Portfolio ID"
// @Router       /portfolios/{id} [get]
func NewGetPortfolioByIdHandler(portfolioService services.PortfolioService) func(echo.Context) error {
	return func(ctx echo.Context) error {
		id := ctx.Param("id")
		portfolio, err := portfolioService.GetPortfolioById(id)

		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, portfolio)
	}
}
