package handlers

import (
	"encoding/csv"
	"strconv"
	"webApp/domain"
	"webApp/internal/service"
	"webApp/pkg/logger"

	"github.com/gin-gonic/gin"
)

type ProductForHandlers struct {
	domain.Product
	logger  *logger.Logger
	service *service.Service
}

func NewProductForHandlers(logger *logger.Logger, service *service.Service) *ProductForHandlers {
	return &ProductForHandlers{
		logger:  logger,
		service: service,
	}
}

func (p *ProductForHandlers) GetAllHandler(c *gin.Context) {
	products, err := p.service.GetAllProducts()
	if len(products) == 0 && err == nil {
		c.JSON(404, gin.H{"Error": "No data found"})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"Internal server error": "error"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/csv")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=products.csv")

	writer := csv.NewWriter(c.Writer)
	writer.Comma = ';'
	defer writer.Flush()

	for _, product := range products {
		record := []string{
			strconv.Itoa(product.Id),
			product.Name,
			product.Price.String(),
		}

		if err := writer.Write(record); err != nil {
			p.logger.Errorf("Error writing CSV record: %s", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
	}
}

func (p *ProductForHandlers) CreateHandler(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		p.logger.Errorf("Invalid input: %s", err.Error())
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	if err := p.service.AddProduct(product); err != nil {
		if err == domain.ErrProductExists {
			p.logger.Errorf("Product exists: %s", err)
			c.JSON(400, gin.H{"error": "Product Exists"})
			return
		}

		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{"message": "Product created"})
}

func (p *ProductForHandlers) DeleteAllHandler(c *gin.Context) {
	if err := p.service.DeleteAllProducts(); err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{"message": "All products deleted"})
}

func (p *ProductForHandlers) DeleteHandler(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.Atoi(id)
	if err != nil || productID == 0 {
		p.logger.Errorf("Invalid product ID: %s", err)
		c.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := p.service.DeleteProduct(productID); err != nil {
		if err == domain.ErrProductNotFound {
			c.JSON(400, gin.H{"error": "Product Not Found"})
			return
		}

		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{"message": "Product deleted"})
}

func (p *ProductForHandlers) UpdateHandler(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.Atoi(id)
	if err != nil || productID == 0 {
		p.logger.Errorf("Invalid product ID: %s", err)
		c.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		p.logger.Errorf("Invalid input: %s", err.Error())
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	product.Id = productID
	if err = p.service.UpdateProduct(product); err != nil {
		if err == domain.ErrProductNotFound {
			c.JSON(404, gin.H{"error": "Product Not Found"})
			return
		}

		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"message": "Product updated"})
}
