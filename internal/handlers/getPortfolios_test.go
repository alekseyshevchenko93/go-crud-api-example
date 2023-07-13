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
	"github.com/stretchr/testify/suite"
)

type GetPortfoliosSuite struct {
	suite.Suite
	portfolioRepository *mocks.PortfolioRepository
	portfolioService    services.PortfolioService
	e                   *echo.Echo
}

func TestGetPortfoliosSuite(t *testing.T) {
	suite.Run(t, new(GetPortfoliosSuite))
}

func (suite *GetPortfoliosSuite) SetupTest() {
	t := suite.T()
	e := echo.New()
	porftolioRepository := mocks.NewPortfolioRepository(t)
	portfolioService := services.NewPortfolioService(porftolioRepository)

	suite.e = e
	suite.portfolioRepository = porftolioRepository
	suite.portfolioService = portfolioService
}

func (suite *GetPortfoliosSuite) TestGetPortfolios() {
	r := suite.Require()
	handler := NewGetPortfoliosHandler(suite.portfolioService)
	portfolios := []*models.Portfolio{
		{Name: "mock-portfolio", IsActive: true, IsFinance: false, IsInternal: false},
		{Name: "mock-portfolio-2", IsActive: false, IsFinance: true},
		{Name: "mock-portfolio-3", IsActive: false, IsFinance: false, IsInternal: true},
	}
	suite.portfolioRepository.EXPECT().GetPortfolios().Return(portfolios, nil).Once()
	portfoliosJson, _ := json.Marshal(portfolios)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, rec)

	err := handler(ctx)

	r.NoError(err)
	r.Contains(rec.Body.String(), string(portfoliosJson))
}
