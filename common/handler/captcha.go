package handler

import (
	"net/http"
	"serve/models"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

func GetCaptcha(c *gin.Context) {
	//数字验证码配置
	var configDigit = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}

	captchaId, captcaInterfaceInstance := base64Captcha.GenerateCaptcha("", configDigit)
	base64blob := base64Captcha.CaptchaWriteToBase64Encoding(captcaInterfaceInstance)

	var data struct {
		CaptchaID  string `json:"captchaID" bson:"captchaID"`
		Base64blob string `json:"base64blob" bson:"base64blob"`
	}

	data.CaptchaID = captchaId
	data.Base64blob = base64blob

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "success",
		"data":   data,
	})

	//or you can just write the captcha content to the httpResponseWriter.
	//before you put the captchaId into the response COOKIE.
	//captcaInterfaceInstance.WriteTo(w)

	//set json response
	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//body := map[string]interface{}{"code": 1, "data": base64blob, "captchaId": captchaId, "msg": "success"}
	//json.NewEncoder(w).Encode(body)
}

func VerifyCaptcha(c *gin.Context) {
	var req models.CaptchaData

	if c.BindJSON(&req) != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "Json parse error!",
		})
	}

	verifyResult := base64Captcha.VerifyCaptcha(req.CaptchaID, req.VerifyValue)
	if !verifyResult {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "验证失败",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "success",
	})
}
