package handlers

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"
	"webApp/domain"
	"webApp/internal/service"
	mock_service "webApp/internal/service/mocks"
	"webApp/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	type mockBehavior func(m *mock_service.MockProduct, product domain.Product)

	testTables := []struct {
		name                 string
		mockBehavior         mockBehavior
		inputBody            string
		input                domain.Product
		expectedStatusCode   int
		expectedResponceBody string
	}{
		{
			name: "Valid",
			mockBehavior: func(m *mock_service.MockProduct, product domain.Product) {
				m.EXPECT().AddProduct(product).Return(nil)
			},
			inputBody: `{"name":"name","price":"50.00"}`,
			input: domain.Product{
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			expectedStatusCode:   200,
			expectedResponceBody: `{"message":"Product created"}`,
		},
		{
			name: "Invalid Input or empty",
			mockBehavior: func(m *mock_service.MockProduct, product domain.Product) {
			},
			inputBody:            `{"name": "", "price":"50.00"}`, // or just `{"price":"50.00"}`
			expectedStatusCode:   400,
			expectedResponceBody: `{"error":"Bad Request"}`,
		},
		{
			name: "Product exists error",
			mockBehavior: func(m *mock_service.MockProduct, product domain.Product) {
				m.EXPECT().AddProduct(product).Return(domain.ErrProductExists)
			},
			inputBody: `{"name": "name", "price":"50.00"}`,
			input: domain.Product{
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			expectedStatusCode:   400,
			expectedResponceBody: `{"error":"Product Exists"}`,
		},
		{
			name: "Internal server error",
			mockBehavior: func(m *mock_service.MockProduct, product domain.Product) {
				m.EXPECT().AddProduct(product).Return(errors.New("some error"))
			},
			inputBody: `{"name": "name", "price":"50.00"}`,
			input: domain.Product{
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			expectedStatusCode:   500,
			expectedResponceBody: `{"error":"Internal Server Error"}`,
		},
	}

	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockProduct(c)
			table.mockBehavior(mockService, table.input)

			service := service.Service{
				Product: mockService,
			}
			logger := logger.GetLogger()
			handler := ProductForHandlers{
				service: &service,
				logger:  logger,
			}

			r := gin.New()
			r.POST("/products", handler.CreateHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/products", bytes.NewBufferString(table.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, table.expectedStatusCode, w.Code)
			assert.Equal(t, table.expectedResponceBody, w.Body.String())
		})
	}
}

func TestGetAllHandler(t *testing.T) {
	type mockBehavior func(m *mock_service.MockProduct, products []domain.Product)

	testTables := []struct {
		name                 string
		mockBehavior         mockBehavior
		input                []domain.Product
		expectedStatusCode   int
		expectedResponceBody string
	}{
		{
			name: "Valid",
			mockBehavior: func(m *mock_service.MockProduct, products []domain.Product) {
				m.EXPECT().GetAllProducts().Return(products, nil)
			},
			input: []domain.Product{
				{
					Name:  "name",
					Price: decimal.RequireFromString("50.00"),
				},
				{
					Name:  "name2",
					Price: decimal.RequireFromString("61.00"),
				},
				{
					Name:  "name3",
					Price: decimal.RequireFromString("72.00"),
				},
			},
			expectedStatusCode: 200,
			expectedResponceBody: `0;name;50
0;name2;61
0;name3;72
`,
		},
		{
			name: "No data found",
			mockBehavior: func(m *mock_service.MockProduct, products []domain.Product) {
				m.EXPECT().GetAllProducts().Return(nil, nil)
			},
			expectedStatusCode:   404,
			expectedResponceBody: `{"Error":"No data found"}`,
		},
		{
			name: "Internal server error",
			mockBehavior: func(m *mock_service.MockProduct, products []domain.Product) {
				m.EXPECT().GetAllProducts().Return(nil, errors.New("some error"))
			},
			input: []domain.Product{
				{
					Name:  "name",
					Price: decimal.RequireFromString("50.00"),
				},
				{
					Name:  "name2",
					Price: decimal.RequireFromString("61.00"),
				},
				{
					Name:  "name3",
					Price: decimal.RequireFromString("72.00"),
				},
			},
			expectedStatusCode:   500,
			expectedResponceBody: `{"Internal server error":"error"}`,
		},
	}

	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockProduct(c)
			table.mockBehavior(mockService, table.input)

			service := service.Service{
				Product: mockService,
			}
			logger := logger.GetLogger()
			handler := ProductForHandlers{
				service: &service,
				logger:  logger,
			}

			r := gin.New()
			r.GET("/products", handler.GetAllHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/products", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, table.expectedStatusCode, w.Code)
			assert.Equal(t, table.expectedResponceBody, w.Body.String())
		})
	}
}

func TestDeleteAllHandler(t *testing.T) {
	type mockBehavior func(m *mock_service.MockProduct)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		//input               int
		expectedStatusCode   int
		expectedResponceBody string
	}{
		{
			name: "Valid",
			mockBehavior: func(m *mock_service.MockProduct) {
				m.EXPECT().DeleteAllProducts().Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponceBody: `{"message":"All products deleted"}`,
		},
		{
			name: "Service error",
			mockBehavior: func(m *mock_service.MockProduct) {
				m.EXPECT().DeleteAllProducts().Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponceBody: `{"error":"Internal Server Error"}`,
		},
	}

	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockProduct(c)
			table.mockBehavior(mockService)

			service := service.Service{
				Product: mockService,
			}
			logger := logger.GetLogger()
			handler := ProductForHandlers{
				service: &service,
				logger:  logger,
			}

			r := gin.New()
			r.POST("/products", handler.DeleteAllHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/products", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, table.expectedStatusCode, w.Code)
			assert.Equal(t, table.expectedResponceBody, w.Body.String())
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	type mockBehavior func(m *mock_service.MockProduct, id int)

	testTables := []struct {
		name                 string
		mockBehavior         mockBehavior
		id                   int
		expectedStatusCode   int
		expectedResponceBody string
	}{
		{
			name: "Valid",
			mockBehavior: func(m *mock_service.MockProduct, id int) {
				m.EXPECT().DeleteProduct(id).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponceBody: `{"message":"Product deleted"}`,
			id:                   1,
		},
		{
			/*Когда обработчик вызвал бы p.service.DeleteProduct(0),
			gomock обнаружила бы неожиданный вызов метода DeleteProduct,
			так как для него не было установлено EXPECT(). Тест бы упал с ошибкой. */
			name: "Invalid ID Format",
			mockBehavior: func(m *mock_service.MockProduct, id int) {
			},
			expectedStatusCode:   400,
			expectedResponceBody: `{"error":"Invalid product ID"}`,
		},
		{
			name: "Product Not Found In Service",
			mockBehavior: func(m *mock_service.MockProduct, id int) {
				m.EXPECT().DeleteProduct(id).Return(domain.ErrProductNotFound)
			},
			expectedStatusCode:   400,
			expectedResponceBody: `{"error":"Product Not Found"}`,
			id:                   1,
		},
		{
			name: "Service error",
			mockBehavior: func(m *mock_service.MockProduct, id int) {
				m.EXPECT().DeleteProduct(id).Return(errors.New("some error"))
			},
			expectedStatusCode:   500,
			expectedResponceBody: `{"error":"Internal Server Error"}`,
			id:                   1,
		},
	}

	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockProduct(c)

			if table.name != "Invalid ID Format" {
				table.mockBehavior(mockService, table.id)
			}
			service := service.Service{
				Product: mockService,
			}

			logger := logger.GetLogger()
			handler := ProductForHandlers{
				service: &service,
				logger:  logger,
			}

			r := gin.New()
			r.DELETE("/products/:id", handler.DeleteHandler)

			w := httptest.NewRecorder()
			url := "/products/"
			if table.name == "Invalid ID Format" {
				url += "invalid-id"
			} else {
				url += strconv.Itoa(table.id)
			}
			req := httptest.NewRequest("DELETE", url, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, table.expectedStatusCode, w.Code)
			assert.Equal(t, table.expectedResponceBody, w.Body.String())

		})
	}
}

func TestUpdateHandler(t *testing.T) {
	type mockBehavior func(m *mock_service.MockProduct, id int, input domain.Product)

	testTables := []struct {
		name                 string
		mockBehavior         mockBehavior
		id                   int
		expectedStatusCode   int
		expectedResponceBody string
		RequestBody          string
		input                domain.Product
		invalidId            string
	}{
		{
			name: "Valid",
			mockBehavior: func(m *mock_service.MockProduct, id int, input domain.Product) {
				m.EXPECT().UpdateProduct(input).Return(nil)
			},
			input: domain.Product{
				Id:    1,
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			id:                   1,
			expectedStatusCode:   200,
			expectedResponceBody: `{"message":"Product updated"}`,
			RequestBody:          `{"name": "name", "price":"50.00"}`,
		},
		{
			name: "Invalid Product ID",
			mockBehavior: func(m *mock_service.MockProduct, id int, input domain.Product) {

			},
			expectedStatusCode:   400,
			invalidId:            "sdsf",
			expectedResponceBody: `{"error":"Invalid product ID"}`,
			RequestBody:          `{"name": "name", "price":"50.00"}`,
		},
		{
			name: "Invalid input Product or empty values",
			mockBehavior: func(m *mock_service.MockProduct, id int, input domain.Product) {

			},
			id:                   1,
			expectedStatusCode:   400,
			expectedResponceBody: `{"error":"Bad Request"}`,
			RequestBody:          `{"name": "", "price":"50.00"}`,
		},
		{
			name: "Product Not Found",
			mockBehavior: func(m *mock_service.MockProduct, id int, input domain.Product) {
				m.EXPECT().UpdateProduct(input).Return(domain.ErrProductNotFound)
			},
			input: domain.Product{
				Id:    1,
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			id:                   1,
			expectedStatusCode:   404,
			expectedResponceBody: `{"error":"Product Not Found"}`,
			RequestBody:          `{"name": "name", "price":"50.00"}`,
		},
		{
			name: "Error of service",
			mockBehavior: func(m *mock_service.MockProduct, id int, input domain.Product) {
				m.EXPECT().UpdateProduct(input).Return(errors.New("some error"))
			},
			input: domain.Product{
				Id:    1,
				Name:  "name",
				Price: decimal.RequireFromString("50.00"),
			},
			id:                   1,
			expectedStatusCode:   500,
			expectedResponceBody: `{"error":"Internal server error"}`,
			RequestBody:          `{"name": "name", "price":"50.00"}`,
		},
	}

	for _, table := range testTables {
		t.Run(table.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockProduct(c)
			table.mockBehavior(mockService, table.id, table.input)

			service := service.Service{
				Product: mockService,
			}
			logger := logger.GetLogger()
			handler := ProductForHandlers{
				service: &service,
				logger:  logger,
			}

			r := gin.New()
			r.PUT("/products/:id", handler.UpdateHandler)

			w := httptest.NewRecorder()

			if table.name == "Invalid Product ID" {
				req := httptest.NewRequest("PUT", "/products/"+table.invalidId, bytes.NewBufferString(table.RequestBody))
				r.ServeHTTP(w, req)
			} else {
				req := httptest.NewRequest("PUT", "/products/"+strconv.Itoa(table.id), bytes.NewBufferString(table.RequestBody))
				r.ServeHTTP(w, req)
			}

			assert.Equal(t, table.expectedStatusCode, w.Code)
			assert.Equal(t, table.expectedResponceBody, w.Body.String())
		})
	}
}
