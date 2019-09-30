package api

import (
	"serve/common/handler"
	"serve/middleware"

	"github.com/gin-gonic/gin"
)

func GetOrganizationAPI(engine *gin.Engine) {
	v1 := engine.Group("/api/v1alpha1")
	v1.Use(middleware.JWTAuth())

	v1.GET("/organization", handler.ListOrg)
	v1.GET("/organization/:_id", handler.GetOrg)
	v1.GET("/organizationname/:name", handler.ListNameOrg)
	v1.POST("/organization", handler.CreateOrg)
	v1.PUT("/organization/:_id", handler.UpdateOrg)
	v1.DELETE("/organization/:_id", handler.DeleteOrg)
}
