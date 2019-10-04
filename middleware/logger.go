package middleware

import (
	//"fmt"
	//"os"
	//"path"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	//logFilePath := config.Log_FILE_PATH
	//logFileName := config.LOG_FILE_NAME

	//日志文件
	//fileName := path.Join(logFilePath, logFileName)

	//写入文件
	//src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	//if err != nil {
	//	fmt.Println("err", err)
	//}

	//实例化
	//logger := logrus.New()

	//设置输出
	//logger.Out = src

	//设置日志级别
	//logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	//logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		log.WithFields(log.Fields{
			"statusCode":  statusCode,
			"latencyTime": latencyTime,
			"clientIP":    clientIP,
			"reqMethod":   reqMethod,
			"reqUri":      reqUri,
		}).Info("coming")
	}
}
