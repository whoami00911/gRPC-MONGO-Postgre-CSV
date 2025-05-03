package repository

import (
	"context"
	"gRPC-server/internal/domain"
	"gRPC-server/pkg/logger"
	"gRPC-server/pkg/parseCSV/grpcPb"
)

type Sorting interface {
	Insert(context.Context, []domain.Product) error
	List(context.Context, domain.SortParams) ([]domain.Product, error)
	GetByName(ctx context.Context, product domain.Product) (domain.Product, error)
	UpdateProduct(ctx context.Context, product domain.Product) error
}

type Repository struct {
	Sorting
	logger *logger.Logger
}

func NewRepo(sort Sorting, logger *logger.Logger) *Repository {
	return &Repository{
		Sorting: sort,
		logger:  logger,
	}
}

func (r *Repository) Fetch(ctx context.Context, req []domain.Product) (domain.Status, error) {
	if err := r.Sorting.Insert(ctx, req); err != nil {
		return domain.Status{}, err
	}
	return domain.Status{}, nil
}

func (r *Repository) List(ctx context.Context, req *grpcPb.ListRequest) ([]domain.Product, error) {
	sortParams := domain.SortParams{
		SortField:    req.GetSortField().String(),
		SortAsc:      req.GetSortAsc(),
		PagingOffset: req.GetPagingOffset(),
		PagingLimit:  req.GetPagingLimit(),
	}
	products, err := r.Sorting.List(ctx, sortParams)
	if err != nil {
		r.logger.Errorf("Can't list sortParams: %s", err)
		return []domain.Product{}, err
	}
	return products, nil
}
