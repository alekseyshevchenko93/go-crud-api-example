package factories

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/models"
	"github.com/alekseyshevchenko93/go-crud-api-example/internal/domain/requests"
)

const minPortfolioId = 1
const maxPortfolioId = 10

func GetRandomName() string {
	n := rand.Intn(maxPortfolioId-minPortfolioId) + minPortfolioId
	return fmt.Sprintf("portfolio-%d", n)
}

func GetRandomBool() bool {
	return rand.Intn(2) == 1
}

func GetPortfolio() *models.Portfolio {
	rand.Seed(time.Now().UnixNano())
	now := time.Now()

	portfolio := models.Portfolio{
		Id:         rand.Intn(maxPortfolioId-minPortfolioId) + minPortfolioId,
		Name:       GetRandomName(),
		IsActive:   GetRandomBool(),
		IsFinance:  GetRandomBool(),
		IsInternal: GetRandomBool(),
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	return &portfolio
}

func GetCreatePortfolioRequest() *requests.CreatePortfolioRequest {
	request := requests.CreatePortfolioRequest{
		Name:       GetRandomName(),
		IsActive:   GetRandomBool(),
		IsFinance:  GetRandomBool(),
		IsInternal: GetRandomBool(),
	}

	return &request
}

func GetUpdatePortfolioRequest() *requests.UpdatePortfolioRequest {
	portfolio := GetPortfolio()

	request := requests.UpdatePortfolioRequest{
		Id:         portfolio.Id,
		Name:       portfolio.Name,
		IsActive:   portfolio.IsActive,
		IsFinance:  portfolio.IsFinance,
		IsInternal: portfolio.IsInternal,
	}

	return &request
}
