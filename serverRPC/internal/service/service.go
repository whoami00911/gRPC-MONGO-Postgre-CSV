package service

import (
	"context"
	"encoding/csv"
	"gRPC-server/internal/domain"
	"gRPC-server/pkg/logger"
	"gRPC-server/pkg/parseCSV/grpcPb"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Sorting interface {
	Fetch(ctx context.Context, product []domain.Product) (domain.Status, error)
	List(ctx context.Context, product *grpcPb.ListRequest) ([]domain.Product, error)
	GetByName(ctx context.Context, product domain.Product) (domain.Product, error)
	UpdateProduct(ctx context.Context, product domain.Product) error
}

type Service struct {
	logger *logger.Logger
	Sorting
}

func NewService(sortService Sorting, logger *logger.Logger) *Service {
	return &Service{
		logger:  logger,
		Sorting: sortService,
	}
}

func (s *Service) Fetch(ctx context.Context, req *grpcPb.FetchRequest) (domain.Status, error) {
	var products []domain.Product

	resp, err := http.Get(req.GetUrl())
	if err != nil {
		s.logger.Errorf("Get URL request error: %s", err)
		return domain.Status{
			Status: "Fail",
		}, err
	}

	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		s.logger.Errorf("Read csv error: %s", err)
		return domain.Status{
			Status: "Fail",
		}, err
	}

	for i, v := range records {
		if len(v) < 3 {
			s.logger.Warnf("skipping invalid record #%d: %v", i, v)
			continue
		}
		Id, _ := strconv.Atoi(v[0])
		price, err := primitive.ParseDecimal128(v[2])
		if err != nil {
			s.logger.Errorf("Decimal parse error: %s", err)
			return domain.Status{
				Status: "Fail",
			}, err
		}

		product := domain.Product{
			Id:    Id,
			Name:  v[1],
			Price: price,
		}

		exists, err := s.Sorting.GetByName(ctx, product)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				products = append(products, product)
			} else {
				return domain.Status{
					Status: "Fail",
				}, err
			}
		}
		if exists.Price != product.Price {
			if err := s.Sorting.UpdateProduct(ctx, product); err != nil {
				return domain.Status{
					Status: "Fail",
				}, err
			}
		}
	}
	status, err := s.Sorting.Fetch(ctx, products)
	if err != nil {
		if err == domain.ErrNoProducts {
			status.Status = "Success"
			return status, nil
		}
		s.logger.Errorf("Fetch request error: %s", err)
		return domain.Status{
			Status: "Fail",
		}, err
	}
	status.Status = "Success"
	return status, nil
}

func (s *Service) List(ctx context.Context, req *grpcPb.ListRequest) ([]domain.Product, error) {
	products, err := s.Sorting.List(ctx, req)
	if err != nil {
		return []domain.Product{}, err
	}
	return products, nil
}
