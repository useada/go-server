package api

import (
	"serve/common/handler"

	"github.com/gin-gonic/gin"
)

func GetCaptchaAPI(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")

	v1.GET("/captcha/get", handler.GetCaptcha)
	v1.POST("/captcha/verify", handler.VerifyCaptcha)
}
