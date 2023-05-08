package handlers

import (
	"encoding/json"
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

func TestGetPortfolioByIdSuccess(t *testing.T) {
	e := echo.New()
	porfoliosRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porfoliosRepository)
	handler := NewGetPortfolioByIdHandler(portfolioService)
	portfolioId := 10
	portfolioIdStr := fmt.Sprintf("%d", portfolioId)
	portfolio := models.Portfolio{
		Name:       "mock-portfolio",
		IsActive:   true,
		IsFinance:  false,
		IsInternal: false,
	}

	porfoliosRepository.EXPECT().GetPortfolioById(portfolioId).Return(portfolio, nil).Once()
	portfolioJson, _ := json.Marshal(portfolio)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/portfolios/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(portfolioIdStr)

	err := handler(ctx)

	assert.NoError(t, err)
	assert.Contains(t, rec.Body.String(), string(portfolioJson))
}

func TestGetPortfolioByIdBadRequest(t *testing.T) {
	e := echo.New()
	porfoliosRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porfoliosRepository)
	handler := NewGetPortfolioByIdHandler(portfolioService)
	tt := []string{"string-instead-of-id", "#$(*)@", ""}

	for _, portfolioId := range tt {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
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

func TestGetPortfolioByIdNotFound(t *testing.T) {
	e := echo.New()
	porfoliosRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porfoliosRepository)
	portfolioId := 1
	portfolioIdStr := fmt.Sprintf("%d", portfolioId)
	porfoliosRepository.EXPECT().GetPortfolioById(portfolioId).Return(models.Portfolio{}, repository.ErrPortfolioNotFound).Once()

	handler := NewGetPortfolioByIdHandler(portfolioService)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
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
