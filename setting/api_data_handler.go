package setting

import (
	"fmt"
	"io/fs"
	"mocknet/utils"
	"path/filepath"
)

type ApiData struct {
	Path   string      `json:"path"`   //数据对应path
	Method string      `json:"method"` // 请求方法
	Key    string      `json:"key"`    //param+query+funcode生成的key
	Close  bool        `json:"close"`  //是否使用当前数据
	Data   interface{} `json:"data"`   //mock 返回的数据
}

func (m ApiData) String() string {
	return fmt.Sprintf("{Path: %s, Method: %s, Key: %s, Data: %v}", m.Path, m.Method, m.Key, m.Data)
}

func (m ApiData) IsEmpty() bool {
	if m.Path == "" || m.Method == "" {
		return true
	}

	return false
}

func (m ApiData) GenerateSaveKey() string {
	return fmt.Sprintf("%s,%s", m.Path, m.Method)
}

type ApiDataHandler struct {
}

func (handler ApiDataHandler) Handle(path string) bool {

	mockData := ApiData{}
	err := utils.LoadFileJson(path, &mockData)

	if err != nil || mockData.Close {
		return false
	}

	GlobalConfigData.AppendApiData(mockData.GenerateSaveKey(), &mockData)

	return true
}

func (handler ApiDataHandler) SupportExt(fi fs.FileInfo) bool {
	return filepath.Ext(fi.Name()) == ".data"
}

var _ ConfigHandler = ApiDataHandler{}
