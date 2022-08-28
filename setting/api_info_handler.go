package setting

import (
	"fmt"
	"mocknet/logger"
	"strings"
)

type ApiData struct {
	Path         string          `json:"path"`          //url path
	QueryKey     []string        `json:"query_key"`     //请求的关心的query信息
	Method       string          `json:"method"`        //request method e.g:POST/GET
	BodyKey      string          `json:"body_key"`      //请求body中的key
	Data         ApiMockInfoData `json:"data"`          //数据配置信息
	NeedRedirect bool            `json:"need_redirect"` //是否需求重定向
}

func (apiData *ApiData) String() string {
	return fmt.Sprintf("{ Path: %s, QueryKey:%s, Method: %s, BodyKey: %s, Data: %s, NeedRedirect: %t}",
		apiData.Path,
		strings.Join(apiData.QueryKey, ","),
		apiData.Method,
		apiData.BodyKey,
		apiData.Data,
		apiData.NeedRedirect)
}

// Merge merge path and method equal data to current data, success return true , fail return false
func (apiData *ApiData) Merge(data ApiData) bool {
	if apiData.Path != data.Path || apiData.Method != data.Method {
		return false
	}

	if apiData.Data != data.Data || apiData.BodyKey != data.BodyKey {
		logger.E("api data info is change:\n origin: %s,\n new: %s \n", apiData, &data)
	}

	apiData.Data = data.Data
	apiData.QueryKey = data.QueryKey
	apiData.BodyKey = data.BodyKey
	apiData.NeedRedirect = data.NeedRedirect

	return true
}

type ApiMockInfoData struct {
	Plugin     string `json:"plugin"`
	FolderPath string `json:"folder_path"`
}

func (data ApiMockInfoData) String() string {
	return fmt.Sprintf("{ Plugin: %s, FolderPath:%s}", data.Plugin, data.FolderPath)
}
