package middleware

import (
	"github.com/gin-gonic/gin"
	"mocknet/setting"
	"net/http"
	"net/http/httputil"
)

func ReverseProxy(context *gin.Context, extDirector func(req *http.Request)) {

	director := func(req *http.Request) {

		req.Host = setting.GetProxyHost()
		req.URL.Scheme = setting.GetProxySchema()
		req.URL.Host = req.Host
		extDirector(req)
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(context.Writer, context.Request)
}
