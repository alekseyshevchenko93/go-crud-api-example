package handlers

import (
	"net/http"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/models"
	requests "github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/requests"
	"github.com/labstack/echo/v4"
)

type PortfolioCreater interface {
	CreatePortfolio(requests.CreatePortfolioRequest) (models.Portfolio, error)
}

func NewCreatePortfolioHandler(portfolioService PortfolioCreater) func(echo.Context) error {
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
