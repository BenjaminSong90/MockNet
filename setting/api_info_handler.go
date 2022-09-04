package setting

import (
	"fmt"
	"io/fs"
	"mocknet/logger"
	"mocknet/utils"
	"path/filepath"
	"strings"
)

type MockApiInfoData struct {
	Path         string      `json:"path"`          //url path
	QueryKey     []string    `json:"query_key"`     //请求的关心的query信息
	Method       string      `json:"method"`        //request method e.g:POST/GET
	BodyKey      string      `json:"body_key"`      //请求body中的key
	Data         MockApiData `json:"data"`          //数据配置信息
	NeedRedirect bool        `json:"need_redirect"` //是否需求重定向
}

func (apiData *MockApiInfoData) String() string {
	return fmt.Sprintf("{ Path: %s, QueryKey:%s, Method: %s, BodyKey: %s, Data: %s, NeedRedirect: %t}",
		apiData.Path,
		strings.Join(apiData.QueryKey, ","),
		apiData.Method,
		apiData.BodyKey,
		apiData.Data,
		apiData.NeedRedirect)
}

// Merge merge path and method equal data to current data, success return true , fail return false
func (apiData *MockApiInfoData) Merge(data MockApiInfoData) bool {
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

func (apiData *MockApiInfoData) GetMockData(key string) (interface{}, bool) {
	k := fmt.Sprintf("%s,%s", apiData.Path, apiData.Method)

	if v, ok := GlobalConfigData.MockData[k]; ok {
		if md, ok := v[key]; ok {
			return md, true
		}
	}
	return nil, false
}

type MockApiData struct {
	Plugin     string `json:"plugin"`
	FolderPath string `json:"folder_path"`
}

func (data MockApiData) String() string {
	return fmt.Sprintf("{ Plugin: %s, FolderPath:%s}", data.Plugin, data.FolderPath)
}

type ApiDataHandler struct {
}

func (handler ApiDataHandler) Handle(configData *ConfigData, path string) bool {
	configData.Lock()
	defer configData.Unlock()

	data := MockApiInfoData{}
	err := utils.LoadFileJson(path, data)

	if err != nil {
		return false
	}

	key := fmt.Sprintf("%s,%s", data.Path, data.Method)

	if v, ok := configData.MockApi[key]; ok {
		v.Merge(data)
	} else {
		configData.MockApi[key] = &data
	}

	return true
}

func (handler ApiDataHandler) SupportExt(fi fs.FileInfo) bool {
	return filepath.Ext(fi.Name()) == ".api"
}
