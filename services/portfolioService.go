package services

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alekseyshevchenko93/go-crud-api-example/domain/models"
	requests "github.com/alekseyshevchenko93/go-crud-api-example/domain/requests"
	"github.com/alekseyshevchenko93/go-crud-api-example/repository"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type PortfolioService struct {
	portfolioRepository repository.PortfolioRepository
}

func (s *PortfolioService) CreatePortfolio(body requests.CreatePortfolioRequest) (models.Portfolio, error) {
	if err := s.ValidatePortfolioCreateRequest(body); err != nil {
		return models.Portfolio{}, err
	}

	portfolio, err := s.portfolioRepository.CreatePortfolio(body)

	if err != nil {
		if errors.Is(err, repository.ErrPortfolioAlreadyExists) {
			return models.Portfolio{}, echo.NewHTTPError(http.StatusConflict, "Portfolio with this name already exists")
		}

		return models.Portfolio{}, err
	}

	return portfolio, nil
}

func (s *PortfolioService) GetPortfolios() ([]models.Portfolio, error) {
	portfolios, err := s.portfolioRepository.GetPortfolios()

	if err != nil {
		return nil, err
	}

	return portfolios, nil
}

func (s *PortfolioService) GetPortfolioById(id string) (models.Portfolio, error) {
	if err := s.ValidatePortfolioId(id); err != nil {
		return models.Portfolio{}, err
	}

	portfolio, err := s.portfolioRepository.GetPortfolioById(id)

	if err != nil {
		if errors.Is(err, repository.ErrPortfolioNotFound) {
			return models.Portfolio{}, echo.NewHTTPError(http.StatusNotFound, "Portfolio not found")
		}

		return models.Portfolio{}, err
	}

	return portfolio, nil
}

func (s *PortfolioService) UpdatePortfolio(id string, body requests.UpdatePortfolioRequest) (models.Portfolio, error) {
	if err := s.ValidatePortfolioId(id); err != nil {
		return models.Portfolio{}, err
	}

	if err := s.ValidatePortfolioUpdateRequest(id, body); err != nil {
		return models.Portfolio{}, err
	}

	portfolio, err := s.portfolioRepository.GetPortfolioById(id)

	if err != nil {
		if errors.Is(err, repository.ErrPortfolioNotFound) {
			return models.Portfolio{}, echo.NewHTTPError(http.StatusNotFound, "Portfolio not found")
		}

		return models.Portfolio{}, err
	}

	portfolio.Name = body.Name
	portfolio.IsActive = body.IsActive
	portfolio.IsFinance = body.IsFinance
	portfolio.IsInternal = body.IsInternal

	updatedPortfolio, err := s.portfolioRepository.UpdatePortfolio(portfolio)

	if err != nil {
		return models.Portfolio{}, err
	}

	return updatedPortfolio, nil
}

func (s *PortfolioService) DeletePortfolio(id string) error {
	if err := s.ValidatePortfolioId(id); err != nil {
		return err
	}

	_, err := s.portfolioRepository.GetPortfolioById(id)

	if err != nil {
		if errors.Is(err, repository.ErrPortfolioNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Portfolio not found")
		}

		return err
	}

	if err := s.portfolioRepository.DeletePortfolio(id); err != nil {
		return err
	}

	return nil
}

func (s *PortfolioService) ValidatePortfolioCreateRequest(body requests.CreatePortfolioRequest) error {
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		errors := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, errors[0])
	}

	return nil
}

func (s *PortfolioService) ValidatePortfolioUpdateRequest(id string, body requests.UpdatePortfolioRequest) error {
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

func (s *PortfolioService) ValidatePortfolioId(id string) error {
	validate := validator.New()

	if err := validate.Var(id, "required,numeric"); err != nil {
		errors := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, errors[0])
	}

	return nil
}

func NewPortfolioService(portfolioRepository repository.PortfolioRepository) *PortfolioService {
	return &PortfolioService{
		portfolioRepository,
	}
}
