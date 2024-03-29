package router

import (
	"github.com/gin-gonic/gin"
	"mocknet/logger"
	"mocknet/server/router/handler"
	"mocknet/setting"
	"strings"
)

func InitApi(router *gin.Engine, apiInfoList map[string]*setting.Api) {
	for _, apiDetail := range apiInfoList {

		switch strings.ToUpper(apiDetail.Method) {
		case "GET":
			router.GET(apiDetail.Path, handler.MethodHandler(apiDetail))
		case "POST":
			router.POST(apiDetail.Path, handler.MethodHandler(apiDetail))
		case "DELETE":
			router.DELETE(apiDetail.Path, handler.MethodHandler(apiDetail))
		case "PUT":
			router.PUT(apiDetail.Path, handler.MethodHandler(apiDetail))
		default:
			logger.E("this method not support %s", apiDetail.Method)
		}
	}
}
