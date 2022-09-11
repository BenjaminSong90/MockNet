package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"mocknet/logger"
	"mocknet/setting"
	"net/http"
)

type MethodHandlerFunc func(detail setting.Api) gin.HandlerFunc

func MethodHandler(detail *setting.Api) gin.HandlerFunc {
	return func(context *gin.Context) {
		handleRequest(context, detail)
	}
}

var jsonHandler = &MimeJsonHandler{}

// 处理有body的quest
func getRequestHandler(context *gin.Context) MimeParamHandler {

	if context.Request.Method == http.MethodGet {
		return jsonHandler
	}

	switch context.ContentType() {
	case binding.MIMEJSON:
		return jsonHandler
	case binding.MIMEXML, binding.MIMEXML2:
		return nil
	case binding.MIMEPROTOBUF:
		return nil
	case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
		return nil
	case binding.MIMEYAML:
		return nil
	case binding.MIMEMultipartPOSTForm:
		return nil
	default: // case MIMEPOSTForm:
		return nil
	}

}

func handleRequest(context *gin.Context, detail *setting.Api) {
	logger.D("request full path: " + context.FullPath())

	handler := getRequestHandler(context)
	data := handler.CollectParam(context, detail.BodyKey, detail.QueryKey)
	key := data.GenerateKey()
	logger.D("key: %s", key)
	logger.D("Data: %s", fmt.Sprint(detail.GetMockData(key)))
	if v, ok := detail.GetMockData(key); ok {
		context.JSON(http.StatusOK, v.Data)
		return
	}

	context.JSON(http.StatusNotFound, "")
}
