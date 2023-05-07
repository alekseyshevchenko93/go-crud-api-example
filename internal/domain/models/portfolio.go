package models

import "time"

type Portfolio struct {
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	IsInternal bool       `json:"isInternal"`
	IsFinance  bool       `json:"isFinance"`
	IsActive   bool       `json:"isActive"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}
