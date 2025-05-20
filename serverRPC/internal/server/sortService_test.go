package server

import (
	"context"
	"errors"
	"gRPC-server/internal/domain"
	mock_server "gRPC-server/internal/server/mocks"
	"gRPC-server/pkg/logger"
	"gRPC-server/pkg/parseCSV/grpcPb"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFetch(t *testing.T) {
	logger := logger.GetLogger()

	type mockBehavior func(m *mock_server.MockSorting, ctx context.Context, req *grpcPb.FetchRequest)

	testTables := []struct {
		name         string
		ctx          context.Context
		req          *grpcPb.FetchRequest
		mockBehavior mockBehavior
		want         *grpcPb.FethResponce
		isErr        bool
	}{
		{
			name: "Valid",
			req: &grpcPb.FetchRequest{
				Url: "localhost:8085/products",
			},
			ctx: context.Background(),
			want: &grpcPb.FethResponce{
				Status: "Success",
			},
			mockBehavior: func(m *mock_server.MockSorting, ctx context.Context, req *grpcPb.FetchRequest) {
				m.EXPECT().Fetch(ctx, req).Return(domain.Status{
					Status: "Success",
				}, nil)
			},
			isErr: false,
		},
		{
			name: "Some Error",
			req: &grpcPb.FetchRequest{
				Url: "localhost:8085/products",
			},
			ctx: context.Background(),
			want: &grpcPb.FethResponce{
				Status: "Fail",
			},
			mockBehavior: func(m *mock_server.MockSorting, ctx context.Context, req *grpcPb.FetchRequest) {
				m.EXPECT().Fetch(ctx, req).Return(domain.Status{
					Status: "Fail",
				}, errors.New("Some Error"))
			},
			isErr: true,
		},
	}

	for i := range testTables {
		table := &testTables[i]
		t.Run(table.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockSortingServiceServer := mock_server.NewMockSorting(c)

			serviceServer := NewSortServerService(mockSortingServiceServer, logger)
			table.mockBehavior(mockSortingServiceServer, table.ctx, table.req)
			got, err := serviceServer.Fetch(table.ctx, table.req)

			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, table.want, got)
			}
		})
	}
}

func TestList(t *testing.T) {
	logger := logger.GetLogger()

	type mockBehavior func(m *mock_server.MockSorting, ctx context.Context, req *grpcPb.ListRequest, product []domain.Product)

	testTables := []struct {
		name         string
		ctx          context.Context
		req          *grpcPb.ListRequest
		want         *grpcPb.ListResponce
		mockBehavior mockBehavior
		product      []domain.Product
		isErr        bool
	}{
		{
			name: "Valid",
			req: &grpcPb.ListRequest{
				SortField:    1,
				PagingOffset: 1,
				PagingLimit:  1,
				SortAsc:      1,
			},
			want: &grpcPb.ListResponce{
				Product: []*grpcPb.Product{
					{
						Id:    1,
						Name:  "name",
						Price: "50.00",
					},
					{
						Id:    2,
						Name:  "Name2",
						Price: "60.00",
					},
					{
						Id:    3,
						Name:  "Name3",
						Price: "70.00",
					},
				},
			},
			ctx: context.Background(),
			product: []domain.Product{
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
			mockBehavior: func(m *mock_server.MockSorting, ctx context.Context, req *grpcPb.ListRequest, product []domain.Product) {
				m.EXPECT().List(ctx, req).Return(product, nil)
			},
			isErr: false,
		},
		{
			name: "Service error",
			req: &grpcPb.ListRequest{
				SortField:    1,
				PagingOffset: 1,
				PagingLimit:  1,
				SortAsc:      1,
			},
			want:    &grpcPb.ListResponce{},
			ctx:     context.Background(),
			product: []domain.Product{},
			mockBehavior: func(m *mock_server.MockSorting, ctx context.Context, req *grpcPb.ListRequest, product []domain.Product) {
				m.EXPECT().List(ctx, req).Return(product, errors.New("some error"))
			},
			isErr: true,
		},
	}

	for i := range testTables {
		table := &testTables[i]
		t.Run(table.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockSortingServiceServer := mock_server.NewMockSorting(c)

			serviceServer := NewSortServerService(mockSortingServiceServer, logger)
			table.mockBehavior(mockSortingServiceServer, table.ctx, table.req, table.product)
			got, err := serviceServer.List(table.ctx, table.req)

			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, table.want, got)
			}
		})
	}
}
