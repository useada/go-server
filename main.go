package main

import (
	"serve/common/api"
	"serve/db"
	"serve/middleware"

	"github.com/gin-gonic/gin"
)

const Port = "9002"

func main() {
	db.Connect()
	router := gin.Default()
	router.Use(middleware.Connect)
	//router.Use(middleware.JWTAuth())
	router.Use(middleware.Cors())
	api.RunHTTPServer(router)
	// router.GET("/test", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "hello world",
	// 	})
	// })

	router.Run(":" + Port)
}
