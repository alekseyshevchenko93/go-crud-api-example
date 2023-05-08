package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/models"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/repository"

	mocks "github.com/alekseyshevchenko93/go-crud-api-example/internal/repository/mocks"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDeletePortfolioSuccess(t *testing.T) {
	e := echo.New()
	portfolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(portfolioRepository)
	portfolioId := 1
	portfolioIdStr := fmt.Sprintf("%d", portfolioId)

	portfolio := models.Portfolio{
		Name:      "Alex",
		IsActive:  true,
		IsFinance: true,
	}

	portfolioRepository.EXPECT().GetPortfolioById(portfolioId).Return(&portfolio, nil).Once()
	portfolioRepository.EXPECT().DeletePortfolio(portfolioId).Return(nil).Once()
	handler := NewDeletePortfolioHandler(portfolioService)

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/portfolios/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(portfolioIdStr)

	err := handler(ctx)

	assert.NoError(t, err)
}

func TestDeletePortfolioBadRequest(t *testing.T) {
	e := echo.New()
	portfolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(portfolioRepository)
	tt := []string{"", "some-string", "#$*(&$#(*!}[-=)"}
	handler := NewDeletePortfolioHandler(portfolioService)

	for _, portfolioId := range tt {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/portfolios/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(portfolioId)

		err := handler(ctx)

		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, httpError.Code, http.StatusBadRequest)
	}
}

func TestDeletePortfolioNotFound(t *testing.T) {
	e := echo.New()
	portfolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(portfolioRepository)
	portfolioId := 1
	portfolioIdStr := fmt.Sprintf("%d", portfolioId)
	portfolioRepository.EXPECT().GetPortfolioById(portfolioId).Return(nil, repository.ErrPortfolioNotFound).Once()

	handler := NewDeletePortfolioHandler(portfolioService)
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/portfolios/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(portfolioIdStr)

	err := handler(ctx)

	assert.Error(t, err)

	httpError, ok := err.(*echo.HTTPError)
	assert.True(t, ok)

	assert.Equal(t, httpError.Code, http.StatusNotFound)
}
