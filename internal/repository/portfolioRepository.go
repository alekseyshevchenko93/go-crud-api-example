package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/models"
	requests "github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/requests"
)

type portfolioRepository struct {
	storage map[int]models.Portfolio
	counter int
	mu      sync.RWMutex
}

var (
	ErrPortfolioNotFound      = errors.New("Portfolio not found")
	ErrPortfolioAlreadyExists = errors.New("Portfolio already exists")
)

//go:generate mockery --name PortfolioRepository
type PortfolioRepository interface {
	GetPortfolios() ([]models.Portfolio, error)
	GetPortfolioById(int) (models.Portfolio, error)
	CreatePortfolio(requests.CreatePortfolioRequest) (models.Portfolio, error)
	UpdatePortfolio(models.Portfolio) (models.Portfolio, error)
	DeletePortfolio(int) error
}

func (p *portfolioRepository) GetPortfolios() ([]models.Portfolio, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	items := make([]models.Portfolio, 0, len(p.storage))

	for _, v := range p.storage {
		items = append(items, v)
	}

	return items, nil
}

func (p *portfolioRepository) CreatePortfolio(body requests.CreatePortfolioRequest) (models.Portfolio, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	name := body.Name

	for _, v := range p.storage {
		if v.Name == name {
			return models.Portfolio{}, ErrPortfolioNotFound
		}
	}

	now := time.Now()
	p.counter = p.counter + 1
	id := p.counter

	model := models.Portfolio{
		Id:         id,
		Name:       body.Name,
		IsFinance:  body.IsFinance,
		IsActive:   body.IsActive,
		IsInternal: body.IsInternal,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	p.storage[id] = model

	return model, nil
}

func (p *portfolioRepository) GetPortfolioById(id int) (models.Portfolio, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	model, ok := p.storage[id]

	if !ok {
		return models.Portfolio{}, ErrPortfolioNotFound
	}

	return model, nil
}

func (p *portfolioRepository) DeletePortfolio(id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.storage[id]; !ok {
		return ErrPortfolioNotFound
	}

	delete(p.storage, id)

	return nil
}

func (p *portfolioRepository) UpdatePortfolio(model models.Portfolio) (models.Portfolio, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, v := range p.storage {
		if v.Name == model.Name && v.Id != model.Id {
			return models.Portfolio{}, ErrPortfolioAlreadyExists
		}
	}

	now := time.Now()
	model.UpdatedAt = &now

	p.storage[model.Id] = model

	return model, nil
}

func NewPortfolioRepository() *portfolioRepository {
	return &portfolioRepository{
		storage: make(map[int]models.Portfolio),
	}
}
