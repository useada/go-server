package api

import (
	"serve/common/handler"
	"serve/middleware"

	"github.com/gin-gonic/gin"
)

func GetProductAPI(engine *gin.Engine) {
	v1 := engine.Group("/api/v1alpha1")
	v1.Use(middleware.JWTAuth())

	v1.GET("/product", handler.ListProduct)
	v1.GET("/product/:_id", handler.GetProduct)
	v1.GET("/productname/:name", handler.ListNameProduct)
	v1.POST("/product", handler.CreateProduct)
	v1.PUT("/product/:_id", handler.UpdateProduct)
	v1.DELETE("/product/:_id", handler.DeleteProduct)

	v1.GET("/organization/:_id/product", handler.ListOrgProducts)
	v1.GET("/customer/:_id/product", handler.ListCustomerProducts)
}
