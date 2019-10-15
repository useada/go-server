package api

import (
	"github.com/gin-gonic/gin"
)

func RunHTTPServer(engine *gin.Engine) {
	GetUserAPI(engine)
	GetLinkAPI(engine)
	GetCommentAPI(engine)
	GetCaptchaAPI(engine)
	GetLikeAPI(engine)

	//GetOrganizationAPI(engine)
	//GetCustomerAPI(engine)
	//GetProductAPI(engine)
	//GetMappingAPI(engine)
	//GetDeviceAPI(engine)
}
