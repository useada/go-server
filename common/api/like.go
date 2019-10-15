package api

import (
	"serve/common/handler"
	"serve/middleware"

	"github.com/gin-gonic/gin"
)

func GetLikeAPI(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")

	v1.Use(middleware.JWTPrepare())
	v1.Use(middleware.JWTAuth())

	v1.GET("/like/list/get", handler.ListUserLike)
	v1.POST("/like/update", handler.UpdateUserLike)
}
