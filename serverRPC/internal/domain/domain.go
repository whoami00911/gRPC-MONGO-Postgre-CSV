package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id    int                  `json:"Id" binding:"required"`
	Name  string               `json:"Name" binding:"required"`
	Price primitive.Decimal128 `json:"Price" binding:"required"`
}

type URL struct {
	Url string
}

type Status struct {
	Status string
}

type SortParams struct {
	SortField    string
	SortAsc      int32
	PagingOffset int32
	PagingLimit  int32
}
