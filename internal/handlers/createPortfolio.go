package handlers

import (
	"net/http"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/requests"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/labstack/echo/v4"
)

// CreatePortfolio create portfolio and responds with new portfolio
// @Summary      Creates portfolio
// @Tags         Portfolios
// @Param        portfolio body requests.CreatePortfolioRequest true "Portfolio Body"
// @Produce      json
// @Success      200  {object}  models.Portfolio
// @Router       /portfolios [post]
func NewCreatePortfolioHandler(portfolioService services.PortfolioService) func(echo.Context) error {
	return func(ctx echo.Context) error {
		body := &requests.CreatePortfolioRequest{}

		if err := ctx.Bind(body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid json body")
		}

		portfolio, err := portfolioService.CreatePortfolio(body)

		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusCreated, portfolio)
	}
}
