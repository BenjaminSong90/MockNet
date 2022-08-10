package middleware

import (
	"mocknet/logger"
	"mocknet/server/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoFundHandle() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
		if context.Writer.Status() == http.StatusNotFound ||
			context.Writer.Status() == http.StatusMethodNotAllowed {
			logger.D(">>>>>>%s<<<<<< is proxy", context.Request.URL.Path)
			router.ReverseProxy(context, func(req *http.Request) {})
		}
	}
}
