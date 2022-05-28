package middleware

import (
	"github.com/gin-gonic/gin"
	"mock_net/config"
	"net/http"
	"net/http/httputil"
	"strings"
)

func NoFundHandle(mockApiList []config.ApiInfo, projectConfig config.ProjectConfig) gin.HandlerFunc {
	return func(context *gin.Context) {
		if checkPath(context.Request.URL.Path, context.Request.Method, mockApiList) &&
			len(projectConfig.ProxyHost) != 0 &&
			len(projectConfig.ProxyScheme) != 0 {

			director := func(req *http.Request) {
				req.Host = projectConfig.ProxyHost
				req.URL.Scheme = projectConfig.ProxyScheme
				req.URL.Host = req.Host
			}
			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(context.Writer, context.Request)

			context.Abort()
		}

	}
}

func checkPath(requestPath string, method string, mockApiList []config.ApiInfo) bool{
	if len(config.PConfig.VideoPath) != 0 && strings.Contains(requestPath, "/videos/") {
		return false
	}
	var hasFond = false
	for _, apiDetail := range mockApiList {
		if apiDetail.Path == requestPath && method == strings.ToUpper(apiDetail.Method) {
			hasFond = true
		}
	}

	return hasFond
}
