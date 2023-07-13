package services

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/models"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/requests"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/repository"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type PortfolioService interface {
	CreatePortfolio(*requests.CreatePortfolioRequest) (*models.Portfolio, error)
	UpdatePortfolio(string, *requests.UpdatePortfolioRequest) (*models.Portfolio, error)
	GetPortfolios() ([]*models.Portfolio, error)
	GetPortfolioById(string) (*models.Portfolio, error)
	DeletePortfolio(string) error
}

type portfolioService struct {
	portfolioRepository repository.PortfolioRepository
}

func (s *portfolioService) CreatePortfolio(body *requests.CreatePortfolioRequest) (*models.Portfolio, error) {
	if err := s.validatePortfolioCreateRequest(body); err != nil {
		return nil, err
	}

	portfolio, err := s.portfolioRepository.CreatePortfolio(body)

	if err != nil {
		if errors.Is(err, repository.ErrPortfolioAlreadyExists) {
			return nil, echo.NewHTTPError(http.StatusConflict, "Portfolio with this name already exists")
		}

		return nil, err
	}

	return portfolio, nil
}

func (s *portfolioService) GetPortfolios() ([]*models.Portfolio, error) {
	portfolios, err := s.portfolioRepository.GetPortfolios()

	if err != nil {
		return nil, err
	}

	return portfolios, nil
}

func (s *portfolioService) GetPortfolioById(id string) (*models.Portfolio, error) {
	if err := s.validatePortfolioId(id); err != nil {
		return nil, err
	}

	idInt, _ := strconv.Atoi(id)
	portfolio, err := s.portfolioRepository.GetPortfolioById(idInt)

	if err != nil {
		if errors.Is(err, repository.ErrPortfolioNotFound) {
			return nil, echo.NewHTTPError(http.StatusNotFound, "Portfolio not found")
		}

		return nil, err
	}

	return portfolio, nil
}

func (s *portfolioService) UpdatePortfolio(id string, body *requests.UpdatePortfolioRequest) (*models.Portfolio, error) {
	if err := s.validatePortfolioId(id); err != nil {
		return nil, err
	}

	if err := s.validatePortfolioUpdateRequest(id, body); err != nil {
		return nil, err
	}

	idInt, _ := strconv.Atoi(id)
	portfolio, err := s.portfolioRepository.GetPortfolioById(idInt)

	if err != nil {
		if errors.Is(err, repository.ErrPortfolioNotFound) {
			return nil, echo.NewHTTPError(http.StatusNotFound, "Portfolio not found")
		}

		return nil, err
	}

	portfolio.Name = body.Name
	portfolio.IsActive = body.IsActive
	portfolio.IsFinance = body.IsFinance
	portfolio.IsInternal = body.IsInternal

	updatedPortfolio, err := s.portfolioRepository.UpdatePortfolio(portfolio)

	if err != nil {
		if errors.Is(err, repository.ErrPortfolioAlreadyExists) {
			return nil, echo.NewHTTPError(http.StatusConflict, "Portfolio with this name already exists")
		}

		return nil, err
	}

	return updatedPortfolio, nil
}

func (s *portfolioService) DeletePortfolio(id string) error {
	if err := s.validatePortfolioId(id); err != nil {
		return err
	}

	idInt, _ := strconv.Atoi(id)
	_, err := s.portfolioRepository.GetPortfolioById(idInt)

	if err != nil {
		if errors.Is(err, repository.ErrPortfolioNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Portfolio not found")
		}

		return err
	}

	if err := s.portfolioRepository.DeletePortfolio(idInt); err != nil {
		return err
	}

	return nil
}

func (s *portfolioService) validatePortfolioCreateRequest(body *requests.CreatePortfolioRequest) error {
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		errors := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, errors[0])
	}

	return nil
}

func (s *portfolioService) validatePortfolioUpdateRequest(id string, body *requests.UpdatePortfolioRequest) error {
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		errors := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, errors[0])
	}

	idInt, _ := strconv.Atoi(id)

	if idInt != body.Id {
		return echo.NewHTTPError(http.StatusBadRequest, "Id in param and body dont match")
	}

	return nil
}

func (s *portfolioService) validatePortfolioId(id string) error {
	validate := validator.New()

	if err := validate.Var(id, "required,numeric"); err != nil {
		errors := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, errors[0])
	}

	return nil
}

func NewPortfolioService(portfolioRepository repository.PortfolioRepository) *portfolioService {
	return &portfolioService{
		portfolioRepository,
	}
}
