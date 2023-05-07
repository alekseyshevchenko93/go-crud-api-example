package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/models"
	mocks "github.com/alekseyshevchenko93/go-crud-api-example/internal/repository/mocks"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetPortfolios(t *testing.T) {
	e := echo.New()
	porfoliosRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porfoliosRepository)

	portfolios := []models.Portfolio{
		{Name: "mock-portfolio", IsActive: true, IsFinance: false, IsInternal: false},
		{Name: "mock-portfolio-2", IsActive: false, IsFinance: true},
		{Name: "mock-portfolio-3", IsActive: false, IsFinance: false, IsInternal: true},
	}
	porfoliosRepository.EXPECT().GetPortfolios().Return(portfolios, nil).Once()
	portfoliosJson, _ := json.Marshal(portfolios)

	handler := NewGetPortfoliosHandler(portfolioService)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler(ctx)

	assert.NoError(t, err)
	assert.Contains(t, rec.Body.String(), string(portfoliosJson))
}
