package router

import (
	"github.com/gin-gonic/gin"
	"mocknet/setting"
	"strings"
)

func InitApi(router *gin.Engine, apiInfoList *[]setting.ApiInfo) {
	for _, apiDetail := range *apiInfoList {

		switch strings.ToUpper(apiDetail.Method) {
		case "GET":
			router.GET(apiDetail.Path, GetHandler(apiDetail))
		case "POST":
			router.POST(apiDetail.Path, PostHandler(apiDetail))
		case "DELETE":
			router.DELETE(apiDetail.Path, DeleteHandler(apiDetail))
		case "PUT":
			router.PUT(apiDetail.Path, PutHandler(apiDetail))
		default:
			router.GET(apiDetail.Path, GetHandler(apiDetail))
		}
	}
}

