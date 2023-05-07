package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/models"
	requests "github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/requests"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/repository"
	mocks "github.com/alekseyshevchenko93/go-crud-api-example/internal/repository/mocks"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreatePortfolioSuccess(t *testing.T) {
	e := echo.New()
	porfoliosRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porfoliosRepository)
	createdAt := time.Now()

	requestBody := requests.CreatePortfolioRequest{
		Name:     "created-portfolio",
		IsActive: true,
	}

	portfolio := models.Portfolio{
		Id:         1,
		Name:       requestBody.Name,
		IsActive:   requestBody.IsActive,
		IsFinance:  false,
		IsInternal: false,
		CreatedAt:  &createdAt,
		UpdatedAt:  &createdAt,
	}

	responseJson, _ := json.Marshal(portfolio)
	requestJson, _ := json.Marshal(requestBody)
	porfoliosRepository.EXPECT().CreatePortfolio(requestBody).Return(portfolio, nil).Once()

	handler := NewCreatePortfolioHandler(portfolioService)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler(ctx)

	assert.NoError(t, err)
	assert.Contains(t, rec.Body.String(), string(responseJson))
}

func TestCreatePortfolioBadRequests(t *testing.T) {
	e := echo.New()
	porfoliosRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porfoliosRepository)
	handler := NewCreatePortfolioHandler(portfolioService)
	tt := []requests.CreatePortfolioRequest{
		{Name: ""},
		{Name: "here-should-be-20-symbols", IsActive: true, IsFinance: false, IsInternal: false},
	}

	for _, body := range tt {
		bodyJson, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		err := handler(ctx)
		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, httpError.Code, http.StatusBadRequest)
	}
}

func TestCreatePortfolioConflict(t *testing.T) {
	e := echo.New()
	porfoliosRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porfoliosRepository)
	handler := NewCreatePortfolioHandler(portfolioService)
	requestBody := requests.CreatePortfolioRequest{
		Name:     "portfolio",
		IsActive: true,
	}
	porfoliosRepository.EXPECT().CreatePortfolio(requestBody).Return(models.Portfolio{}, repository.ErrPortfolioAlreadyExists).Once()
	bodyJson, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler(ctx)

	assert.Error(t, err)
	httpError, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, httpError.Code, http.StatusConflict)
}
