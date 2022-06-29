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
	"strings"
)

/**
	如果 api 是restful请求：
	字典的key RequestURI+key_name按照顺序进行拼接： /user/1?t=1,key1V,key2V,key3V
	key1V是用key1从header和body找查找的，如果没有会直接跳过例如 上面key2没有查找到的，生成的key是：/user/1?t=1,key1V,key3V
	key的查找顺序是header 和 body,如果header和body中都存在，body中的value会覆盖header中的value

	如果 api 不是 restful 请求：
	字典的key key_name按照顺序进行拼接： key1V,key2V,key3V
	key1V是用key1从header和body找查找的，如果没有会直接跳过例如 上面key2没有查找到的，生成的key是：key1V,key3V
	key的查找顺序是header 和 body,如果header和body中都存在，body中的value会覆盖header中的value
	当没有找到后也会使用restful进行查找
 */

type MethodHandlerFunc func(detail setting.ApiInfo) gin.HandlerFunc

func GetHandler(detail *setting.ApiInfo) gin.HandlerFunc{
	return func(context *gin.Context) {
		logger.DebugLogger("request full path"+context.Request.RequestURI)
		handleRequest(context, detail)
		//context.JSON(http.StatusOK, detail.Data[context.Request.URL.Path])
	}
}

func DeleteHandler(detail *setting.ApiInfo) gin.HandlerFunc{
	return func(context *gin.Context) {
		handleRequest(context, detail)
		//context.JSON(http.StatusOK, detail.Data[context.Request.RequestURI])
	}
}

func PostHandler(detail *setting.ApiInfo) gin.HandlerFunc{
	return func(context *gin.Context) {
		handleBodyRequest(context, detail)
	}
}

func PutHandler(detail *setting.ApiInfo) gin.HandlerFunc{
	return func(context *gin.Context) {
		handleBodyRequest(context, detail)
	}
}



//处理有body的quest
func handleBodyRequest(context *gin.Context, detail *setting.ApiInfo) {

	if detail.Restful {
		handleRequest(context , detail)
		//context.JSON(http.StatusOK, detail.Data[context.Request.RequestURI])
	} else {
		switch context.ContentType() {
		case binding.MIMEJSON:
			//handleJsonType(context, detail)
			handleRequest(context , detail)
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

func handleRequest(context *gin.Context, detail *setting.ApiInfo){
	if detail.Restful {

		ks := collectKey(context, detail)
		v := ""
		if len(ks) == 0{
			v = context.Request.RequestURI
		} else {
			v = context.Request.RequestURI + "," +strings.Join(ks, ",")
		}
		logger.DebugLogger(">>>>>v: "+v)

		if v,ok := detail.Data[v]; ok{
			context.JSON(http.StatusOK, v)
			return
		}

		ReverseProxy(context, func(req *http.Request) {})
	} else {
		ks := collectKey(context, detail)

		v := ""
		if len(ks) != 0 {
			v = strings.Join(ks, ",")
		}

		if v,ok := detail.Data[v]; ok{
			context.JSON(http.StatusOK, v)
			return
		}
		if len(v) == 0{
			v = context.Request.RequestURI
		} else {
			v = context.Request.RequestURI + "," + v
		}
		if v,ok := detail.Data[v]; ok{
			context.JSON(http.StatusOK, v)
			return
		}

		ReverseProxy(context, func(req *http.Request) {})
	}
}

func collectKey(context *gin.Context, detail *setting.ApiInfo)[]string{
	result := make([]string ,0)
	if len(detail.KeyName) == 0{
		return result
	}
	keys := strings.Split(detail.KeyName, ",")
	collector := make(map[string]string)
	collectHeaderInfo(context, keys, collector)
	collectBodyInfo(context, keys, collector)
	for _,k := range keys{
		v, ok := collector[k]
		if ok {
			result = append(result, v)
		}
	}
	return result
}

func collectHeaderInfo(context *gin.Context, keys []string, result map[string]string){
	for _,k := range keys{
		v := context.GetHeader(k)
		if len(v) != 0{
			result[k] = v
		}
	}
}

func collectBodyInfo(context *gin.Context, keys []string,  result map[string]string){

	jsonData, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		return
	}
	jsonBody := make(map[string]interface{})
	err = json.Unmarshal(jsonData, &jsonBody)
	if err != nil {
		return
	}
	kv := make(map[string]interface{})
	utils.FlatMap(jsonBody, kv)
	for _, k := range keys{
		v,ok := kv[k]
		if sv,s := v.(string); s && ok{
			result[k] = sv
		}
	}

	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))
}

