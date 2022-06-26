package router

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"mocknet/logger"
	"mocknet/setting"
	"mocknet/utils"
	"net/http"
)

type MethodHandlerFunc func(detail setting.ApiInfo) gin.HandlerFunc

func GetHandler(detail setting.ApiInfo) gin.HandlerFunc{
	return func(context *gin.Context) {
		logger.DebugLogger("request full path"+context.Request.RequestURI)
		context.JSON(http.StatusOK, detail.Data[context.Request.URL.Path])
	}
}

func DeleteHandler(detail setting.ApiInfo) gin.HandlerFunc{
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, detail.Data[context.Request.RequestURI])
	}
}

func PostHandler(detail setting.ApiInfo) gin.HandlerFunc{
	return func(context *gin.Context) {
		handleBodyRequest(context, detail)
	}
}

func PutHandler(detail setting.ApiInfo) gin.HandlerFunc{
	return func(context *gin.Context) {
		handleBodyRequest(context, detail)
	}
}



//处理有body的quest
func handleBodyRequest(context *gin.Context, detail setting.ApiInfo) {

	if detail.Restful {
		context.JSON(http.StatusOK, detail.Data[context.Request.RequestURI])
	} else {
		switch context.ContentType() {
		case binding.MIMEJSON:
			handleJsonType(context, detail)
		case binding.MIMEXML, binding.MIMEXML2:
			fallthrough
		case binding.MIMEPROTOBUF:
			fallthrough
		case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
			fallthrough
		case binding.MIMEYAML:
			fallthrough
		case binding.MIMEMultipartPOSTForm:
			fallthrough
		default: // case MIMEPOSTForm:
			context.JSON(http.StatusOK, gin.H{"status": "error", "message": "content-type not support"})
		}
	}
}

//处理content type 为 application/json 的类型
func handleJsonType(context *gin.Context, detail setting.ApiInfo) {
	jsonData, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"status": "error", "message": err})
		return
	}
	jsonBody := make(map[string]interface{})
	err = json.Unmarshal(jsonData, &jsonBody)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"status": "error", "message": err})
		return
	}
	kvMap := make(map[string]interface{})
	utils.FlatMap(jsonBody, kvMap)
	functionCode, exist := kvMap[detail.KeyName]
	if exist {
		data, exist := detail.Data[functionCode.(string)]
		if exist {
			context.JSON(http.StatusOK, data)
		} else {
			ReverseProxy(context, func(req *http.Request) {
				req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))
			})
		}
	} else {
		ReverseProxy(context, func(req *http.Request) {
			req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))
		})
	}
}

