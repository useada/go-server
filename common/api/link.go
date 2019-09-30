package api

import (
	"serve/common/handler"
	"serve/middleware"

	"github.com/gin-gonic/gin"
)

func GetLinkAPI(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")

	v1.GET("/links", handler.GetLinks)
	v1.GET("/link/:_id", handler.GetLink)

	v1.Use(middleware.JWTAuth())

	//v1.GET("/links", handler.ListLinks)
	v1.PUT("/link", handler.CreateLink)
	v1.POST("/link/:_id", handler.UpdateLink)
	v1.DELETE("/link/:_id", handler.DeleteLink)
}
