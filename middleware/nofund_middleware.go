package middleware

import (
	"github.com/gin-gonic/gin"
	"mock_net/server/router"
	"mock_net/setting"
	"net/http"
	"strings"
)

func NoFundHandle(apiList *[]setting.ApiInfo) gin.HandlerFunc {
	return func(context *gin.Context) {
		if !fundApiPath(context.Request.URL.Path, context.Request.Method, apiList) &&
			setting.CheckProxyInfo(){

			router.ReverseProxy(context, func(req *http.Request) {})

			context.Abort()
		}

	}
}

func fundApiPath(requestPath string, method string, apiList *[]setting.ApiInfo) bool {
	if len(setting.GetStaticFilePath()) != 0 && strings.Contains(requestPath, "/file") {
		return true
	}
	var hasFond = false
	for _, apiDetail := range *apiList {
		if apiDetail.Path == requestPath && method == strings.ToUpper(apiDetail.Method) {
			hasFond = true
		}
	}

	return hasFond
}
