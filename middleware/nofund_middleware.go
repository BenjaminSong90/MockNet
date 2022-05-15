package middleware

import (
	"github.com/gin-gonic/gin"
	"mock_net/model"
	"net/http"
	"net/http/httputil"
)

func NoFundHandle(mockConfig model.MockConfig, mockApiInfo model.MockApiInfo) gin.HandlerFunc{
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		var hasFond = false
		for _, apiDetail := range mockApiInfo.ApiInfo {
			if apiDetail.Path == path{
				hasFond = true
			}
		}

		if !hasFond && len(mockConfig.ProxyHost) != 0 && len(mockConfig.Scheme) != 0{
			director := func(req *http.Request) {
				req.Host = mockConfig.ProxyHost
				req.URL.Scheme = mockConfig.Scheme
				req.URL.Host = req.Host
			}
			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(context.Writer, context.Request)

			context.Abort()
		}

	}
}
