package repository

import (
	"context"
	"errors"
	"gRPC-server/internal/domain"
	mock_repository "gRPC-server/internal/repository/mocks"
	"gRPC-server/pkg/logger"
	"gRPC-server/pkg/parseCSV/grpcPb"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFetch(t *testing.T) {
	type mockBehavior func(m *mock_repository.MockSorting, ctx context.Context, products []domain.Product)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		products     []domain.Product
		ctx          context.Context
		status       domain.Status
		isErr        bool
	}{
		{
			name: "Valid",
			mockBehavior: func(m *mock_repository.MockSorting, ctx context.Context, products []domain.Product) {
				m.EXPECT().Insert(ctx, products).Return(nil)
			},
			products: []domain.Product{
				{
					Id:   1,
					Name: "name",
					Price: func() primitive.Decimal128 {
						got, _ := primitive.ParseDecimal128("50.00")
						return got
					}(),
				},
				{
					Id:   2,
					Name: "Name2",
					Price: func() primitive.Decimal128 {
						got, _ := primitive.ParseDecimal128("60.00")
						return got
					}(),
				},
				{
					Id:   3,
					Name: "Name3",
					Price: func() primitive.Decimal128 {
						got, _ := primitive.ParseDecimal128("70.00")
						return got
					}(),
				},
			},
			isErr: false,
		},
		{
			name: "Insert method Error",
			mockBehavior: func(m *mock_repository.MockSorting, ctx context.Context, products []domain.Product) {
				m.EXPECT().Insert(ctx, products).Return(errors.New("some err"))
			},
			products: []domain.Product{
				{
					Id:   1,
					Name: "name",
					Price: func() primitive.Decimal128 {
						got, _ := primitive.ParseDecimal128("50.00")
						return got
					}(),
				},
				{
					Id:   2,
					Name: "Name2",
					Price: func() primitive.Decimal128 {
						got, _ := primitive.ParseDecimal128("60.00")
						return got
					}(),
				},
				{
					Id:   3,
					Name: "Name3",
					Price: func() primitive.Decimal128 {
						got, _ := primitive.ParseDecimal128("70.00")
						return got
					}(),
				},
			},
			isErr: true,
		},
	}
	logger := logger.GetLogger()
	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepo := mock_repository.NewMockSorting(c)
			repo := NewRepo(mockRepo, logger)

			table.mockBehavior(mockRepo, table.ctx, table.products)

			got, err := repo.Fetch(table.ctx, table.products)

			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, table.status, got)
			}

		})
	}

}

func TestList(t *testing.T) {
	type mockBehavior func(m *mock_repository.MockSorting, ctx context.Context, sortParams domain.SortParams, products []domain.Product)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		products     []domain.Product
		sortParams   domain.SortParams
		ctx          context.Context
		req          *grpcPb.ListRequest
		isErr        bool
	}{
		{
			name: "Valid",
			mockBehavior: func(m *mock_repository.MockSorting, ctx context.Context, sortParams domain.SortParams, products []domain.Product) {
				m.EXPECT().List(ctx, sortParams).Return(products, nil)
			},
			products: []domain.Product{
				{
					Id:   1,
					Name: "name",
					Price: func() primitive.Decimal128 {
						got, _ := primitive.ParseDecimal128("50.00")
						return got
					}(),
				},
				{
					Id:   2,
					Name: "Name2",
					Price: func() primitive.Decimal128 {
						got, _ := primitive.ParseDecimal128("60.00")
						return got
					}(),
				},
				{
					Id:   3,
					Name: "Name3",
					Price: func() primitive.Decimal128 {
						got, _ := primitive.ParseDecimal128("70.00")
						return got
					}(),
				},
			},
			sortParams: domain.SortParams{
				SortField:    "name",
				SortAsc:      1,
				PagingOffset: 1,
				PagingLimit:  1,
			},
			req: &grpcPb.ListRequest{
				SortField:    1,
				SortAsc:      1,
				PagingOffset: 1,
				PagingLimit:  1,
			},
			isErr: false,
		},
		{
			name: "Repository error",
			mockBehavior: func(m *mock_repository.MockSorting, ctx context.Context, sortParams domain.SortParams, products []domain.Product) {
				m.EXPECT().List(ctx, sortParams).Return(products, errors.New("some error"))
			},
			products: []domain.Product{},
			sortParams: domain.SortParams{
				SortField:    "name",
				SortAsc:      1,
				PagingOffset: 1,
				PagingLimit:  1,
			},
			req: &grpcPb.ListRequest{
				SortField:    1,
				SortAsc:      1,
				PagingOffset: 1,
				PagingLimit:  1,
			},
			isErr: true,
		},
	}
	logger := logger.GetLogger()
	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepo := mock_repository.NewMockSorting(c)
			repo := NewRepo(mockRepo, logger)

			table.mockBehavior(mockRepo, table.ctx, table.sortParams, table.products)

			got, err := repo.List(table.ctx, table.req)

			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, table.products, got)
			}

		})
	}
}
