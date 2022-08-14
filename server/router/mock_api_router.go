package router

import (
	"github.com/gin-gonic/gin"
	"mocknet/setting"
	"strings"
)

func InitApi(router *gin.Engine, apiInfoList *[]*setting.ApiInfo) {
	for _, apiDetail := range *apiInfoList {

		switch strings.ToUpper(apiDetail.Method) {
		case "GET":
			router.GET(apiDetail.Path, MethodHandler(apiDetail))
		case "POST":
			router.POST(apiDetail.Path, MethodHandler(apiDetail))
		case "DELETE":
			router.DELETE(apiDetail.Path, MethodHandler(apiDetail))
		case "PUT":
			router.PUT(apiDetail.Path, MethodHandler(apiDetail))
		default:
			router.GET(apiDetail.Path, MethodHandler(apiDetail))
		}
	}
}
