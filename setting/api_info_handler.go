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
	Path         string   `json:"path"`          //url path
	QueryKey     []string `json:"query_key"`     //请求的关心的query信息
	Method       string   `json:"method"`        //request method e.g:POST/GET
	BodyKey      string   `json:"body_key"`      //请求body中的key
	NeedRedirect bool     `json:"need_redirect"` //是否需求重定向
	Plugin       string   `json:"plugin"`
}

func (apiData *MockApiInfoData) String() string {
	return fmt.Sprintf("{ Path: %s, QueryKey:%s, Method: %s, BodyKey: %s, Data: %s, NeedRedirect: %t}",
		apiData.Path,
		strings.Join(apiData.QueryKey, ","),
		apiData.Method,
		apiData.BodyKey,
		apiData.Plugin,
		apiData.NeedRedirect)
}

// Merge merge path and method equal data to current data, success return true , fail return false
func (apiData *MockApiInfoData) Merge(data MockApiInfoData) bool {
	if apiData.Path != data.Path || apiData.Method != data.Method {
		return false
	}

	if apiData.Plugin != data.Plugin || apiData.BodyKey != data.BodyKey {
		logger.E("api data info is change:\n origin: %s,\n new: %s \n", apiData, &data)
	}

	apiData.Plugin = data.Plugin
	apiData.QueryKey = data.QueryKey
	apiData.BodyKey = data.BodyKey
	apiData.NeedRedirect = data.NeedRedirect

	return true
}

func (apiData *MockApiInfoData) GetMockData(key string) (*MockData, bool) {
	k := fmt.Sprintf("%s,%s", apiData.Path, apiData.Method)

	if v, ok := GlobalConfigData.MockData[k]; ok {
		if md, ok := v[key]; ok {
			return md, true
		}
	}
	return nil, false
}

type ApiDataHandler struct {
}

func (handler ApiDataHandler) Handle(configData *ConfigData, path string) bool {
	configData.Lock()
	defer configData.Unlock()

	data := MockApiInfoData{}
	err := utils.LoadFileJson(path, &data)

	if err != nil {
		logger.E("path: %s, error info: %s", path, err.Error())
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
