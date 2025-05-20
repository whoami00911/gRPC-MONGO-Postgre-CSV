package server

import (
	"context"
	"gRPC-server/internal/domain"
	"gRPC-server/pkg/logger"
	"gRPC-server/pkg/parseCSV/grpcPb"
)

//go:generate mockgen -source=sortService.go -destination=mocks/mock.go
type Sorting interface {
	Fetch(ctx context.Context, req *grpcPb.FetchRequest) (domain.Status, error)
	List(ctx context.Context, req *grpcPb.ListRequest) ([]domain.Product, error)
}

type SortServicegRPC struct {
	grpcPb.UnimplementedSortServiceServer
	logger *logger.Logger
	Sorting
}

func NewSortServerService(sortingService Sorting, logger *logger.Logger) *SortServicegRPC {
	return &SortServicegRPC{
		logger:  logger,
		Sorting: sortingService,
	}
}

func (s *SortServicegRPC) Fetch(ctx context.Context, req *grpcPb.FetchRequest) (*grpcPb.FethResponce, error) {
	status, err := s.Sorting.Fetch(ctx, req)
	if err != nil {
		return &grpcPb.FethResponce{
			Status: status.Status,
		}, err
	}
	return &grpcPb.FethResponce{
		Status: status.Status,
	}, nil
}

func (s *SortServicegRPC) List(ctx context.Context, req *grpcPb.ListRequest) (*grpcPb.ListResponce, error) {
	products, err := s.Sorting.List(ctx, req)
	if err != nil {
		return &grpcPb.ListResponce{}, err
	}
	productsGrpc := make([]*grpcPb.Product, len(products))

	for i, product := range products {
		productsGrpc[i] = &grpcPb.Product{
			Id:    int64(product.Id),
			Name:  product.Name,
			Price: product.Price.String(),
		}
	}
	return &grpcPb.ListResponce{
		Product: productsGrpc,
	}, nil
}
