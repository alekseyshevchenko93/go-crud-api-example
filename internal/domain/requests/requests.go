package domain

type CreatePortfolioRequest struct {
	Name       string `json:"name" validate:"required,lte=20"`
	IsInternal bool   `json:"isInternal"`
	IsFinance  bool   `json:"isFinance"`
	IsActive   bool   `json:"isActive"`
}

type UpdatePortfolioRequest struct {
	Id         int    `json:"id" validate:"required,numeric"`
	Name       string `json:"name" validate:"required,lte=20"`
	IsInternal bool   `json:"isInternal"`
	IsFinance  bool   `json:"isFinance"`
	IsActive   bool   `json:"isActive"`
}
