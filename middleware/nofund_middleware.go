package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mocknet/server/router"
	"net/http"
)

func NoFundHandle() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
		if context.Writer.Status() != http.StatusOK{
			fmt.Printf(">>>>>>%s<<<<<< is proxy\n", context.Request.URL.Path)
			router.ReverseProxy(context, func(req *http.Request) {})
		}
	}
}
