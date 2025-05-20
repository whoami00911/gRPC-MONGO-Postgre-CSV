package handlers

import "github.com/gin-gonic/gin"

func (p *ProductForHandlers) InitRoutes() *gin.Engine {
	router := gin.Default()

	products := router.Group("/products")
	{
		products.POST("/", p.CreateHandler)
		products.PUT("/:id", p.UpdateHandler)
		products.GET("/", p.GetAllHandler)
		products.DELETE("/:id", p.DeleteHandler)
	}
	return router
}
