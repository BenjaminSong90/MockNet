package router

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"mock_net/setting"
	"mock_net/utils"
	"net/http"
	"strings"
)

func InitApi(router *gin.Engine, apiInfoList *[]setting.ApiInfo) {
	for _, apiDetail := range *apiInfoList {

		switch strings.ToUpper(apiDetail.Method) {
		case "GET":
			AddGetApi(router, apiDetail)
		case "POST":
			AddPostApi(router, apiDetail)
		case "DELETE":
			AddDeleteApi(router, apiDetail)
		case "PUT":
			AddPutApi(router, apiDetail)
		default:
			AddGetApi(router, apiDetail)
		}
	}
}

func AddGetApi(r *gin.Engine, detail setting.ApiInfo) {
	r.GET(detail.Path, func(context *gin.Context) {
		context.JSON(http.StatusOK, detail.Data[context.Request.URL.Path])
	})
}

func AddPostApi(r *gin.Engine, detail setting.ApiInfo) {
	r.POST(detail.Path, func(context *gin.Context) {
		handleBodyRequest(context, detail)
	})
}

func AddDeleteApi(r *gin.Engine, detail setting.ApiInfo) {
	r.DELETE(detail.Path, func(context *gin.Context) {
		context.JSON(http.StatusOK, detail.Data[context.Request.URL.Path])
	})
}

func AddPutApi(r *gin.Engine, detail setting.ApiInfo) {
	r.PUT(detail.Path, func(context *gin.Context) {
		handleBodyRequest(context, detail)
	})
}

//处理有body的quest
func handleBodyRequest(context *gin.Context, detail setting.ApiInfo) {

	if detail.Restful {
		context.JSON(http.StatusOK, detail.Data[context.Request.URL.Path])
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
		context.JSON(http.StatusOK, gin.H{"status": err})
		return
	}
	jsonBody := make(map[string]interface{})
	err = json.Unmarshal(jsonData, &jsonBody)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"status": err})
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
