package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"mock_net/middleware"
	"mock_net/model"
	"mock_net/utils"
	"net/http"
	"net/http/httputil"
	"strings"
)

var mockConfig = model.MockConfig{}

func main() {

	err := utils.LoadFileJson("config/mock_config.json", &mockConfig)
	if err != nil {
		panic(err)
	}

	mockApiInfo := model.MockApiInfo{}

	err = utils.LoadFileJson("config/api.json", &mockApiInfo)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(middleware.NoFundHandle(mockConfig, mockApiInfo))

	for _, apiDetail := range mockApiInfo.ApiInfo {

		switch strings.ToUpper(apiDetail.Method) {
		case "GET":
			addGetApi(r, apiDetail)
		case "POST":
			addPostApi(r, apiDetail)
		case "DELETE":
			addDeleteApi(r, apiDetail)
		case "PUT":
			addPutApi(r, apiDetail)
		default:
			addGetApi(r, apiDetail)
		}
	}
	address := "8080"
	if len(mockConfig.Address) != 0 {
		address = mockConfig.Address
	}

	r.Run(address) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func addGetApi(r *gin.Engine, detail model.ApiInfoDetail) {
	r.GET(detail.Path, func(context *gin.Context) {
		context.JSON(http.StatusOK, detail.Data[context.Request.URL.Path])
	})
}

func addPostApi(r *gin.Engine, detail model.ApiInfoDetail) {
	r.POST(detail.Path, func(context *gin.Context) {
		handleBodyRequest(context, detail)
	})
}

func addDeleteApi(r *gin.Engine, detail model.ApiInfoDetail) {
	r.DELETE(detail.Path, func(context *gin.Context) {
		context.JSON(http.StatusOK, detail.Data[context.Request.URL.Path])
	})
}

func addPutApi(r *gin.Engine, detail model.ApiInfoDetail) {
	r.PUT(detail.Path, func(context *gin.Context) {
		handleBodyRequest(context, detail)
	})
}

//处理有body的quest
func handleBodyRequest(context *gin.Context, detail model.ApiInfoDetail) {

	if detail.IsRestful {
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
func handleJsonType(context *gin.Context, detail model.ApiInfoDetail) {
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
			reverseProxy(context, jsonData)
		}
	} else {
		reverseProxy(context, jsonData)
	}
}

//转发已经解析过body的请求
func reverseProxy(context *gin.Context, jsonData []byte) {
	director := func(req *http.Request) {
		req.Host = mockConfig.ProxyHost
		req.URL.Scheme = mockConfig.ProxyScheme
		req.URL.Host = req.Host
		req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))
		fmt.Println("proxy api: " + req.URL.String())
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(context.Writer, context.Request)
}
