package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/requests"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/repository"
	mocks "github.com/alekseyshevchenko93/go-crud-api-example/internal/repository/mocks"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/services"
	"github.com/alekseyshevchenko93/go-crud-api-example/test/factories"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type CreatePortfolioSuite struct {
	suite.Suite
	portfolioRepository *mocks.PortfolioRepository
	portfolioService    services.PortfolioService
	e                   *echo.Echo
}

func TestCreatePortfolioSuite(t *testing.T) {
	suite.Run(t, new(CreatePortfolioSuite))
}

func (suite *CreatePortfolioSuite) SetupTest() {
	t := suite.T()
	e := echo.New()
	porftolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porftolioRepository)

	suite.e = e
	suite.portfolioRepository = porftolioRepository
	suite.portfolioService = portfolioService
}

func (suite *CreatePortfolioSuite) TestCreatePortfolioSuccess() {
	r := suite.Require()
	handler := NewCreatePortfolioHandler(suite.portfolioService)
	portfolio := factories.GetPortfolio()
	requestBody := requests.CreatePortfolioRequest{
		Name:     portfolio.Name,
		IsActive: portfolio.IsActive,
	}
	responseJson, _ := json.Marshal(portfolio)
	requestJson, _ := json.Marshal(requestBody)

	suite.portfolioRepository.EXPECT().CreatePortfolio(&requestBody).Return(portfolio, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, rec)

	err := handler(ctx)

	r.NoError(err)
	r.Contains(rec.Body.String(), string(responseJson))
}

func (suite *CreatePortfolioSuite) TestCreatePortfolioBadRequests() {
	r := suite.Require()
	handler := NewCreatePortfolioHandler(suite.portfolioService)

	tt := []interface{}{
		&requests.CreatePortfolioRequest{Name: ""},
		&requests.CreatePortfolioRequest{Name: "here-should-be-20-symbols", IsActive: true, IsFinance: false, IsInternal: false},
		"invalid json body",
	}

	for _, body := range tt {
		bodyJson, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := suite.e.NewContext(req, rec)

		err := handler(ctx)

		r.Error(err)
		httpError, ok := err.(*echo.HTTPError)
		r.True(ok)
		r.Equal(httpError.Code, http.StatusBadRequest)
	}
}

func (suite *CreatePortfolioSuite) TestCreatePortfolioConflict() {
	r := suite.Require()
	handler := NewCreatePortfolioHandler(suite.portfolioService)
	requestBody := factories.GetCreatePortfolioRequest()
	bodyJson, _ := json.Marshal(requestBody)

	suite.portfolioRepository.EXPECT().CreatePortfolio(requestBody).Return(nil, repository.ErrPortfolioAlreadyExists).Once()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, rec)

	err := handler(ctx)

	r.Error(err)
	httpError, ok := err.(*echo.HTTPError)
	r.True(ok)
	r.Equal(httpError.Code, http.StatusConflict)
}
