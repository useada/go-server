package api

import (
	"serve/common/handler"
	"serve/middleware"

	"github.com/gin-gonic/gin"
)

func GetCommentAPI(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")

	v1.Use(middleware.JWTPrepare())

	v1.GET("/comments", handler.GetComments)

	v1.Use(middleware.JWTAuth())
	v1.PUT("/comment", handler.CreateComment)

	v1.Use(middleware.JWTAdmin())
	v1.DELETE("/comment/:_id", handler.DeleteComment)
}
