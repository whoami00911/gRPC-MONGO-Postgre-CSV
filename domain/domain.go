package domain

import "github.com/shopspring/decimal"

type Product struct {
	Id    int             `json:"Id" binding:"required"`
	Name  string          `json:"Name" binding:"required"`
	Price decimal.Decimal `json:"Price" binding:"required"`
}
