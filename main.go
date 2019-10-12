package main

import (
	"os"
	"serve/common/api"
	"serve/db"
	"serve/middleware"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const Port = "9002"

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{})

	log.SetLevel(log.DebugLevel)

	file, err := os.OpenFile("./logs/funny-link.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("failed to log to file")
	}
	log.WithFields(log.Fields{}).Info("serve running")

	db.Connect()
	router := gin.Default()
	router.Use(middleware.LoggerToFile())
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
