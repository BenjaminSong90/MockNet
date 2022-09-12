package setting

import (
	"io/fs"
	"mocknet/logger"
	"path/filepath"
	"sync"
)

type MockConfigData struct {
	Data map[string]map[string]*ApiData // map[path+method]map[param+query+funcCode]ApiData
	Api  map[string]*Api                // map[path+method]MockApiInfoData
	sync.Mutex
}

var GlobalConfigData = MockConfigData{
	Data: make(map[string]map[string]*ApiData),
	Api:  make(map[string]*Api),
}

func AppendApiData(key string, mockData *ApiData) {
	GlobalConfigData.Lock()
	defer GlobalConfigData.Unlock()

	if v, ok := GlobalConfigData.Data[key]; ok {
		v[mockData.Key] = mockData
	} else {
		mockDataMap := make(map[string]*ApiData)
		GlobalConfigData.Data[key] = mockDataMap
		mockDataMap[mockData.Key] = mockData
	}
}

// ClearConfigData 清除api 相关的数据缓存
func ClearConfigData() {
	GlobalConfigData.Lock()
	defer GlobalConfigData.Unlock()
	GlobalConfigData.Data = make(map[string]map[string]*ApiData)
	GlobalConfigData.Api = make(map[string]*Api)
}

func AppendApi(key string, api *Api) {
	GlobalConfigData.Lock()
	defer GlobalConfigData.Unlock()
	if v, ok := GlobalConfigData.Api[key]; ok {
		v.Merge(api)
	} else {
		GlobalConfigData.Api[key] = api
	}
}

type ConfigHandler interface {
	SupportExt(fi fs.FileInfo) bool

	Handle(path string) bool
}

// EmptyConfigHandler help print not support file ext info
type EmptyConfigHandler struct {
}

func (handler EmptyConfigHandler) Handle(path string) bool {
	logger.W("file: %s is not handle", path)
	return true
}

func (handler EmptyConfigHandler) SupportExt(_ fs.FileInfo) bool {
	return true
}

var configHandlers = []ConfigHandler{ApiHandler{}, ApiDataHandler{}, EmptyConfigHandler{}}

// HandleConfigFile handle config file
func HandleConfigFile(path string, fi fs.FileInfo) {
	for _, handler := range configHandlers {
		if handler.SupportExt(fi) && handler.Handle(path) {
			break
		}
	}
}

// file路径集合 加载json文件 信息
func loadApiInfo(filePathList []string) {
	ClearConfigData()
	for _, p := range filePathList {
		_ = filepath.Walk(p, func(jsonPath string, info fs.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				HandleConfigFile(jsonPath, info)
			}
			return nil
		})
	}
}
