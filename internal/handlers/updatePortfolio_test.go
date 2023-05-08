package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestUpdatePortfolioSuccess(t *testing.T) {
	e := echo.New()
	portfolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(portfolioRepository)
	portfolioId := 1
	portfolioIdStr := fmt.Sprintf("%d", portfolioId)
	createdAt := time.Now()

	requestBody := requests.UpdatePortfolioRequest{
		Id:       portfolioId,
		Name:     "updated-portfolio",
		IsActive: true,
	}

	portfolio := models.Portfolio{
		Id:         portfolioId,
		Name:       requestBody.Name,
		IsActive:   requestBody.IsActive,
		IsFinance:  false,
		IsInternal: false,
		CreatedAt:  &createdAt,
		UpdatedAt:  &createdAt,
	}

	responseJson, _ := json.Marshal(portfolio)
	requestJson, _ := json.Marshal(requestBody)
	portfolioRepository.EXPECT().GetPortfolioById(portfolioId).Return(&portfolio, nil).Once()
	portfolioRepository.EXPECT().UpdatePortfolio(&portfolio).Return(&portfolio, nil).Once()

	handler := NewUpdatePortfolioHandler(portfolioService)
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(requestJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/portfolios/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(portfolioIdStr)

	err := handler(ctx)

	assert.NoError(t, err)
	assert.Contains(t, rec.Body.String(), string(responseJson))
}

func TestUpdatePortfolioBadRequests(t *testing.T) {
	e := echo.New()
	portfolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(portfolioRepository)
	handler := NewUpdatePortfolioHandler(portfolioService)
	tt := []struct {
		ParamId string
		Body    requests.UpdatePortfolioRequest
	}{
		{
			ParamId: "1",
			Body:    requests.UpdatePortfolioRequest{Id: 1, Name: ""},
		},
		{
			ParamId: "1",
			Body:    requests.UpdatePortfolioRequest{Id: 1, Name: "here-should-be-20-symbols"},
		},
		{
			ParamId: "1",
			Body:    requests.UpdatePortfolioRequest{Id: 2, Name: "random-1"},
		},
		{
			ParamId: "some-test-string",
			Body:    requests.UpdatePortfolioRequest{Id: 2, Name: "random-2"},
		},
	}

	for _, testcase := range tt {
		paramId := testcase.ParamId
		body := testcase.Body
		bodyJson, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(bodyJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/portfolios/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(paramId)

		err := handler(ctx)

		assert.Error(t, err)
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, httpError.Code, http.StatusBadRequest)
	}
}

func TestUpdatePortfolioNotFound(t *testing.T) {
	e := echo.New()
	portfolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(portfolioRepository)
	handler := NewUpdatePortfolioHandler(portfolioService)
	portfolioId := 1
	portfolioIdStr := fmt.Sprintf("%d", portfolioId)

	requestBody := requests.UpdatePortfolioRequest{
		Id:       portfolioId,
		Name:     "updated-portfolio",
		IsActive: true,
	}

	responseErr := echo.NewHTTPError(http.StatusNotFound)
	portfolioRepository.EXPECT().GetPortfolioById(portfolioId).Return(nil, responseErr).Once()
	bodyJson, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(bodyJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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

func TestUpdatePortfolioConflict(t *testing.T) {
	e := echo.New()
	portfolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(portfolioRepository)
	handler := NewUpdatePortfolioHandler(portfolioService)
	portfolioId := 1
	portfolioIdStr := fmt.Sprintf("%d", portfolioId)
	createdAt := time.Now()

	requestBody := requests.UpdatePortfolioRequest{
		Id:       portfolioId,
		Name:     "updated-portfolio",
		IsActive: true,
	}

	portfolio := models.Portfolio{
		Id:         portfolioId,
		Name:       requestBody.Name,
		IsActive:   requestBody.IsActive,
		IsFinance:  false,
		IsInternal: false,
		CreatedAt:  &createdAt,
		UpdatedAt:  &createdAt,
	}

	portfolioRepository.EXPECT().GetPortfolioById(portfolioId).Return(&portfolio, nil).Once()
	portfolioRepository.EXPECT().UpdatePortfolio(&portfolio).Return(nil, repository.ErrPortfolioAlreadyExists).Once()
	bodyJson, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(bodyJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/portfolios/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(portfolioIdStr)

	err := handler(ctx)

	assert.Error(t, err)
	httpError, ok := err.(*echo.HTTPError)

	assert.True(t, ok)
	assert.Equal(t, httpError.Code, http.StatusConflict)
}
