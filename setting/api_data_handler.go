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

type MockDataHandler struct {
}

func (handler MockDataHandler) Handle(configData *ConfigData, path string) bool {
	configData.Lock()
	defer configData.Unlock()

	mockData := ApiData{}
	err := utils.LoadFileJson(path, &mockData)

	if err != nil || mockData.Close {
		return false
	}

	key := fmt.Sprintf("%s,%s", mockData.Path, mockData.Method)

	if v, ok := configData.MockData[key]; ok {
		v[mockData.Key] = &mockData
	} else {
		mockDataMap := make(map[string]*ApiData)
		configData.MockData[key] = mockDataMap
		mockDataMap[mockData.Key] = &mockData
	}

	return true
}

func (handler MockDataHandler) SupportExt(fi fs.FileInfo) bool {
	return filepath.Ext(fi.Name()) == ".data"
}
