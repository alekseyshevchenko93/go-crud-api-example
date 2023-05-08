package handlers

import (
	"net/http"

	requests "github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/requests"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/labstack/echo/v4"
)

func NewCreatePortfolioHandler(portfolioService services.PortfolioService) func(echo.Context) error {
	return func(ctx echo.Context) error {
		body := requests.CreatePortfolioRequest{}

		if err := ctx.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid json body")
		}

		portfolio, err := portfolioService.CreatePortfolio(body)

		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusCreated, portfolio)
	}
}
