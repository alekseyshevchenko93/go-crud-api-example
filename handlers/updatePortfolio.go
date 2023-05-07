package handlers

import (
	"net/http"

	"github.com/alekseyshevchenko93/go-crud-api-example/domain/models"
	requests "github.com/alekseyshevchenko93/go-crud-api-example/domain/requests"
	"github.com/labstack/echo/v4"
)

type PortfolioUpdater interface {
	UpdatePortfolio(string, requests.UpdatePortfolioRequest) (models.Portfolio, error)
}

func NewUpdatePortfolioHandler(portfolioService PortfolioUpdater) func(echo.Context) error {
	return func(ctx echo.Context) error {
		body := requests.UpdatePortfolioRequest{}
		id := ctx.Param("id")

		if err := ctx.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid json body")
		}

		portfolio, err := portfolioService.UpdatePortfolio(id, body)

		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, portfolio)
	}
}
