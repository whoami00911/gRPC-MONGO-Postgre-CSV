package service

import (
	"context"
	"errors"
	"gRPC-server/internal/domain"
	mock_service "gRPC-server/internal/service/mocks"
	"gRPC-server/pkg/logger"
	"gRPC-server/pkg/parseCSV/grpcPb"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* //переделать
func TestFetch(t *testing.T) {
	logger := logger.GetLogger()
	type mockBehavior func(m *mock_service.MockSorting, ctx context.Context, req *grpcPb.FetchRequest, products []domain.Product)
	//type Handler func(products []domain.Product, c *gin.Context)
	testTables := []struct {
		name         string
		ctx          context.Context
		req          *grpcPb.FetchRequest
		mockBehavior mockBehavior
		Url          string
		want         *grpcPb.FethResponce
		products     []domain.Product
		product      domain.Product
		requestBody  string
		Handler      func(c *gin.Context)
		isErr        bool
	}{
		{
			name: "Valid",
			Url:  "http://localhost:8888/products",
			req: &grpcPb.FetchRequest{
				Url: "http://localhost:8888/products",
			},
			ctx: context.Background(),
			want: &grpcPb.FethResponce{
				Status: "Success",
			},
			Handler: func(c *gin.Context) {
				products := []domain.Product{
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
				}
				writer := csv.NewWriter(c.Writer)
				// Устанавливаем разделитель ';'
				writer.Comma = ';'
				defer writer.Flush()

				// Записываем строки данных
				for _, product := range products {
					record := []string{
						strconv.Itoa(product.Id),
						product.Name,
						product.Price.String(),
					}
					if err := writer.Write(record); err != nil {
						c.JSON(500, gin.H{"error": "Internal Server Error"})
						return
					}
				}
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
			requestBody: "1;name;50.00\n2;Name2;60.00\n3;Name3;70.00",
			mockBehavior: func(m *mock_service.MockSorting, ctx context.Context, req *grpcPb.FetchRequest, products []domain.Product) {
				for _, product := range products {
					m.EXPECT().GetByName(ctx, product).Return(product, mongo.ErrNoDocuments)
				}

				m.EXPECT().Fetch(ctx, products).Return(domain.Status{
					Status: "Success",
				}, nil)
			},
			isErr: false,
		},
	}

	r := gin.New()

	go r.Run("localhost:8888")

	for i := range testTables {
		table := &testTables[i]
		t.Run(table.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			new_service := mock_service.NewMockSorting(c)

			service := NewService(new_service, logger)
			table.mockBehavior(new_service, table.ctx, table.req, table.products)

			time.Sleep(1 * time.Second)

			r.POST("/products/", table.Handler)
			//req := httptest.NewRequest("POST", "http://localhost:8888/poducts/", bytes.NewBufferString(table.requestBod

			got, err := service.Fetch(table.ctx, table.req)
			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, table.want, got)
			}
		})
	}
}
*/

func TestList(t *testing.T) {
	type mockBehavior func(m *mock_service.MockSorting, ctx context.Context, req *grpcPb.ListRequest, products []domain.Product)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		products     []domain.Product
		ctx          context.Context
		req          *grpcPb.ListRequest
		isErr        bool
	}{
		{
			name: "Valid",
			mockBehavior: func(m *mock_service.MockSorting, ctx context.Context, req *grpcPb.ListRequest, products []domain.Product) {
				m.EXPECT().List(ctx, req).Return(products, nil)
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
			mockBehavior: func(m *mock_service.MockSorting, ctx context.Context, req *grpcPb.ListRequest, products []domain.Product) {
				m.EXPECT().List(ctx, req).Return(products, errors.New("some error"))
			},
			req: &grpcPb.ListRequest{
				SortField:    1,
				SortAsc:      1,
				PagingOffset: 1,
				PagingLimit:  1,
			},
			products: []domain.Product{},
			isErr:    true,
		},
	}

	logger := logger.GetLogger()
	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			mockService := mock_service.NewMockSorting(c)

			service := NewService(mockService, logger)

			table.mockBehavior(mockService, table.ctx, table.req, table.products)

			got, err := service.List(table.ctx, table.req)

			if table.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, table.products, got)
			}
		})
	}
}
