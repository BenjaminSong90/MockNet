package setting

import (
	"fmt"
	"io/fs"
	"mocknet/logger"
	"mocknet/utils"
	"path/filepath"
	"strings"
)

type Api struct {
	Path         string   `json:"path"`          //url path
	QueryKey     []string `json:"query_key"`     //请求的关心的query信息
	Method       string   `json:"method"`        //request method e.g:POST/GET
	BodyKey      string   `json:"body_key"`      //请求body中的key
	NeedRedirect bool     `json:"need_redirect"` //是否需求重定向
	Plugin       string   `json:"plugin"`
	Data         ApiData  `json:"data"` //api对应的数据
}

func (apiData *Api) String() string {
	return fmt.Sprintf("{ Path: %s, QueryKey:%s, Method: %s, BodyKey: %s, Data: %s, NeedRedirect: %t}",
		apiData.Path,
		strings.Join(apiData.QueryKey, ","),
		apiData.Method,
		apiData.BodyKey,
		apiData.Plugin,
		apiData.NeedRedirect)
}

// Merge merge path and method equal data to current data, success return true , fail return false
func (apiData *Api) Merge(data *Api) bool {
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

func (apiData *Api) GetMockData(key string) (*ApiData, bool) {
	k := fmt.Sprintf("%s,%s", apiData.Path, apiData.Method)

	if v, ok := GlobalConfigData.Data[k]; ok {
		if md, ok := v[key]; ok {
			return md, true
		}
	}
	return nil, false
}

type ApiHandler struct {
}

func (handler ApiHandler) Handle(path string) bool {

	data := Api{}
	err := utils.LoadFileJson(path, &data)

	if err != nil {
		logger.E("path: %s, error info: %s", path, err.Error())
		return false
	}

	key := fmt.Sprintf("%s,%s", data.Path, data.Method)

	AppendApi(key, &data)
	if !data.Data.IsEmpty() {
		AppendApiData(data.Data.GenerateSaveKey(), &data.Data)
	}

	return true
}

func (handler ApiHandler) SupportExt(fi fs.FileInfo) bool {
	return filepath.Ext(fi.Name()) == ".api"
}

var _ ConfigHandler = ApiHandler{}
