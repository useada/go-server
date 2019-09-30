package api

import (
	"serve/common/handler"
	"serve/middleware"

	"github.com/gin-gonic/gin"
)

func GetCustomerAPI(engine *gin.Engine) {
	v1 := engine.Group("/api/v1alpha1")
	v1.Use(middleware.JWTAuth())

	v1.GET("/customer", handler.ListCustomer)
	v1.GET("/customer/:_id", handler.GetCustomer)
	v1.GET("/customername/:name", handler.ListNameCustomer)
	v1.POST("/customer", handler.CreateCustomer)
	v1.PUT("/customer/:_id", handler.UpdateCustomer)
	v1.DELETE("/customer/:_id", handler.DeleteCustomer)

	v1.GET("/organization/:_id/customer", handler.ListOrgCustomers)
	v1.GET("/product/:_id/customer", handler.ListProductCustomers)
}
