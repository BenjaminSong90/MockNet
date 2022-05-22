package middleware

import (
	"github.com/gin-gonic/gin"
	"mock_net/config"
	"net/http"
	"net/http/httputil"
	"strings"
)

func NoFundHandle(apiInfoList []config.ApiInfo, projectConfig config.ProjectConfig) gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		var hasFond = false
		for _, apiDetail := range apiInfoList {
			if apiDetail.Path == path && context.Request.Method == strings.ToUpper(apiDetail.Method) {
				hasFond = true
			}
		}

		if !hasFond && len(projectConfig.ProxyHost) != 0 && len(projectConfig.ProxyScheme) != 0 {

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
