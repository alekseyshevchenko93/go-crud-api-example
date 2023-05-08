package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/repository"
	mocks "github.com/alekseyshevchenko93/go-crud-api-example/internal/repository/mocks"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/alekseyshevchenko93/go-crud-api-example/test/factories"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type GetPortfolioByIdSuite struct {
	suite.Suite
	portfolioRepository *mocks.PortfolioRepository
	portfolioService    services.PortfolioService
	e                   *echo.Echo
}

func TestGetPortfolioByIdSuite(t *testing.T) {
	suite.Run(t, new(GetPortfolioByIdSuite))
}

func (suite *GetPortfolioByIdSuite) SetupSuite() {
	t := suite.T()
	e := echo.New()
	porftolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porftolioRepository)

	suite.e = e
	suite.portfolioRepository = porftolioRepository
	suite.portfolioService = portfolioService
}

func (suite *GetPortfolioByIdSuite) TestGetPortfolioByIdSuccess() {
	r := suite.Require()
	handler := NewGetPortfolioByIdHandler(suite.portfolioService)
	portfolio := factories.GetPortfolio()
	portfolioIdStr := fmt.Sprintf("%d", portfolio.Id)

	suite.portfolioRepository.EXPECT().GetPortfolioById(portfolio.Id).Return(portfolio, nil).Once()
	portfolioJson, _ := json.Marshal(portfolio)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, rec)
	ctx.SetPath("/portfolios/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(portfolioIdStr)

	err := handler(ctx)

	r.NoError(err)
	r.Contains(rec.Body.String(), string(portfolioJson))
}

func (suite *GetPortfolioByIdSuite) TestGetPortfolioByIdBadRequest() {
	r := suite.Require()
	handler := NewGetPortfolioByIdHandler(suite.portfolioService)
	tt := []string{"string-instead-of-id", "#$(*)@", ""}

	for _, portfolioId := range tt {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := suite.e.NewContext(req, rec)
		ctx.SetPath("/portfolios/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(portfolioId)

		err := handler(ctx)

		r.Error(err)
		httpError, ok := err.(*echo.HTTPError)
		r.True(ok)
		r.Equal(httpError.Code, http.StatusBadRequest)
	}
}

func (suite *GetPortfolioByIdSuite) TestGetPortfolioByIdNotFound() {
	r := suite.Require()
	handler := NewGetPortfolioByIdHandler(suite.portfolioService)
	portfolio := factories.GetPortfolio()
	portfolioIdStr := fmt.Sprintf("%d", portfolio.Id)

	suite.portfolioRepository.EXPECT().GetPortfolioById(portfolio.Id).Return(nil, repository.ErrPortfolioNotFound).Once()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, rec)
	ctx.SetPath("/portfolios/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(portfolioIdStr)

	err := handler(ctx)

	r.Error(err)
	httpError, ok := err.(*echo.HTTPError)
	r.True(ok)
	r.Equal(httpError.Code, http.StatusNotFound)
}
