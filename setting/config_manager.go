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

func (configData *MockConfigData) AppendApiData(key string, mockData *ApiData) {
	configData.Lock()
	defer configData.Unlock()

	if v, ok := configData.Data[key]; ok {
		v[mockData.Key] = mockData
	} else {
		mockDataMap := make(map[string]*ApiData)
		configData.Data[key] = mockDataMap
		mockDataMap[mockData.Key] = mockData
	}
}

// ClearConfigData 清除api 相关的数据缓存
func (configData *MockConfigData) ClearConfigData() {
	configData.Lock()
	defer configData.Unlock()
	configData.Data = make(map[string]map[string]*ApiData)
	configData.Api = make(map[string]*Api)
}

func (configData *MockConfigData) AppendApi(key string, api *Api) {
	configData.Lock()
	defer configData.Unlock()
	if v, ok := configData.Api[key]; ok {
		v.Merge(api)
	} else {
		configData.Api[key] = api
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
	GlobalConfigData.ClearConfigData()
	for _, p := range filePathList {
		_ = filepath.Walk(p, func(jsonPath string, info fs.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				HandleConfigFile(jsonPath, info)
			}
			return nil
		})
	}
}
