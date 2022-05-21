package middleware

import (
	"github.com/gin-gonic/gin"
	"mock_net/model"
	"net/http"
	"net/http/httputil"
	"strings"
)

func NoFundHandle(apiInfoList [] model.ApiInfo) gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		var hasFond = false
		for _, apiDetail := range apiInfoList {
			if apiDetail.Path == path && context.Request.Method == strings.ToUpper(apiDetail.Method) {
				hasFond = true
			}
		}
		var mockConfig model.ProjectConfig
		data, _ := context.Get("project_config")
		mockConfig, hasProjectConfig := data.(model.ProjectConfig)

		if !hasFond && hasProjectConfig && len(mockConfig.ProxyHost) != 0 && len(mockConfig.ProxyScheme) != 0 {

			director := func(req *http.Request) {
				req.Host = mockConfig.ProxyHost
				req.URL.Scheme = mockConfig.ProxyScheme
				req.URL.Host = req.Host
			}
			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(context.Writer, context.Request)

			context.Abort()
		}

	}
}
