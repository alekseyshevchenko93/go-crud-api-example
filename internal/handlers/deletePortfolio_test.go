package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/repository"
	"github.com/alekseyshevchenko93/go-crud-api-example/test/factories"

	mocks "github.com/alekseyshevchenko93/go-crud-api-example/internal/repository/mocks"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type DeletePortfolioSuite struct {
	suite.Suite
	portfolioRepository *mocks.PortfolioRepository
	portfolioService    services.PortfolioService
	e                   *echo.Echo
}

func TestDeletePortfolioSuite(t *testing.T) {
	suite.Run(t, new(GetPortfolioByIdSuite))
}

func (suite *DeletePortfolioSuite) SetupSuite() {
	t := suite.T()
	e := echo.New()
	porftolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porftolioRepository)

	suite.e = e
	suite.portfolioRepository = porftolioRepository
	suite.portfolioService = portfolioService
}

func (suite *DeletePortfolioSuite) TestDeletePortfolioSuccess() {
	r := suite.Require()
	handler := NewDeletePortfolioHandler(suite.portfolioService)
	portfolio := factories.GetPortfolio()
	portfolioIdStr := fmt.Sprintf("%d", portfolio.Id)

	suite.portfolioRepository.EXPECT().GetPortfolioById(portfolio.Id).Return(portfolio, nil).Once()
	suite.portfolioRepository.EXPECT().DeletePortfolio(portfolio.Id).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, rec)
	ctx.SetPath("/portfolios/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(portfolioIdStr)

	err := handler(ctx)

	r.NoError(err)
}

func (suite *DeletePortfolioSuite) TestDeletePortfolioBadRequest() {
	r := suite.Require()
	handler := NewDeletePortfolioHandler(suite.portfolioService)

	tt := []string{"", "some-string", "#$*(&$#(*!}[-=)"}

	for _, portfolioId := range tt {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
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

func (suite *DeletePortfolioSuite) TestDeletePortfolioNotFound() {
	r := suite.Require()
	handler := NewDeletePortfolioHandler(suite.portfolioService)
	portfolio := factories.GetPortfolio()
	portfolioIdStr := fmt.Sprintf("%d", portfolio.Id)
	suite.portfolioRepository.EXPECT().GetPortfolioById(portfolio.Id).Return(nil, repository.ErrPortfolioNotFound).Once()

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
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
