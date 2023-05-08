package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	requests "github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/requests"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/repository"
	mocks "github.com/alekseyshevchenko93/go-crud-api-example/internal/repository/mocks"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/alekseyshevchenko93/go-crud-api-example/test/factories"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type UpdatePortfolioSuite struct {
	suite.Suite
	portfolioRepository *mocks.PortfolioRepository
	portfolioService    services.PortfolioService
	e                   *echo.Echo
}

func TestUpdatePortfolioSuite(t *testing.T) {
	suite.Run(t, new(UpdatePortfolioSuite))
}

func (suite *UpdatePortfolioSuite) SetupSuite() {
	t := suite.T()
	e := echo.New()
	porftolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porftolioRepository)

	suite.e = e
	suite.portfolioRepository = porftolioRepository
	suite.portfolioService = portfolioService
}

func (suite *UpdatePortfolioSuite) TestUpdatePortfolioSuccess() {
	r := suite.Require()
	handler := NewUpdatePortfolioHandler(suite.portfolioService)
	portfolio := factories.GetPortfolio()
	portfolioIdStr := fmt.Sprintf("%d", portfolio.Id)

	portfolioCopy := *portfolio
	updatedPortfolio := &portfolioCopy
	updatedPortfolio.Name = "updated-portfolio"
	updatedPortfolio.IsActive = true

	requestBody := requests.UpdatePortfolioRequest{
		Id:         updatedPortfolio.Id,
		Name:       updatedPortfolio.Name,
		IsActive:   updatedPortfolio.IsActive,
		IsFinance:  updatedPortfolio.IsFinance,
		IsInternal: updatedPortfolio.IsInternal,
	}

	responseJson, _ := json.Marshal(updatedPortfolio)
	requestJson, _ := json.Marshal(requestBody)
	suite.portfolioRepository.EXPECT().GetPortfolioById(portfolio.Id).Return(portfolio, nil).Once()
	suite.portfolioRepository.EXPECT().UpdatePortfolio(updatedPortfolio).Return(updatedPortfolio, nil).Once()

	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(requestJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, rec)
	ctx.SetPath("/portfolios/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(portfolioIdStr)

	err := handler(ctx)

	r.NoError(err)
	r.Contains(rec.Body.String(), string(responseJson))
}

func (suite *UpdatePortfolioSuite) TestUpdatePortfolioBadRequests() {
	r := suite.Require()
	handler := NewUpdatePortfolioHandler(suite.portfolioService)

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
		ctx := suite.e.NewContext(req, rec)
		ctx.SetPath("/portfolios/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(paramId)

		err := handler(ctx)

		r.Error(err)
		httpError, ok := err.(*echo.HTTPError)
		r.True(ok)
		r.Equal(httpError.Code, http.StatusBadRequest)
	}
}

func (suite *UpdatePortfolioSuite) TestUpdatePortfolioNotFound() {
	r := suite.Require()
	handler := NewUpdatePortfolioHandler(suite.portfolioService)
	portfolio := factories.GetPortfolio()
	portfolioIdStr := fmt.Sprintf("%d", portfolio.Id)

	requestBody := requests.UpdatePortfolioRequest{
		Id:         portfolio.Id,
		Name:       "updated-portfolio",
		IsActive:   !portfolio.IsActive,
		IsFinance:  !portfolio.IsFinance,
		IsInternal: !portfolio.IsInternal,
	}

	responseErr := echo.NewHTTPError(http.StatusNotFound)
	suite.portfolioRepository.EXPECT().GetPortfolioById(portfolio.Id).Return(nil, responseErr).Once()
	bodyJson, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(bodyJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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

func (suite *UpdatePortfolioSuite) TestUpdatePortfolioConflict() {
	r := suite.Require()
	handler := NewUpdatePortfolioHandler(suite.portfolioService)
	portfolio := factories.GetPortfolio()
	portfolioIdStr := fmt.Sprintf("%d", portfolio.Id)

	portfolioCopy := *portfolio
	updatedPortfolio := &portfolioCopy
	updatedPortfolio.Name = "updated-portfolio"

	requestBody := requests.UpdatePortfolioRequest{
		Id:         updatedPortfolio.Id,
		Name:       "updated-portfolio",
		IsActive:   updatedPortfolio.IsActive,
		IsFinance:  updatedPortfolio.IsFinance,
		IsInternal: updatedPortfolio.IsInternal,
	}

	suite.portfolioRepository.EXPECT().GetPortfolioById(portfolio.Id).Return(portfolio, nil).Once()
	suite.portfolioRepository.EXPECT().UpdatePortfolio(updatedPortfolio).Return(nil, repository.ErrPortfolioAlreadyExists).Once()

	bodyJson, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(bodyJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, rec)
	ctx.SetPath("/portfolios/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(portfolioIdStr)

	err := handler(ctx)

	r.Error(err)
	httpError, ok := err.(*echo.HTTPError)

	r.True(ok)
	r.Equal(httpError.Code, http.StatusConflict)
}
