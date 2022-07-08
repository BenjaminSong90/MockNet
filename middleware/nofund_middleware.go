package middleware

import (
	"github.com/gin-gonic/gin"
	"mocknet/logger"
	"mocknet/server/router"
	"net/http"
)

func NoFundHandle() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
		if context.Writer.Status() == http.StatusNotFound ||
			context.Writer.Status() == http.StatusMethodNotAllowed {
			logger.DebugLogger(">>>>>>%s<<<<<< is proxy", context.Request.URL.Path)
			router.ReverseProxy(context, func(req *http.Request) {})
		}
	}
}
