package setting

import (
	"io/fs"
	"mocknet/logger"
	"path/filepath"
	"sync"
)

type ConfigData struct {
	MockData map[string]map[string]MockData // map[path+method]map[param+query+funcode]MockData
	MockApi  map[string]*MockApiInfoData    // map[path+method]MockApiInfoData
	sync.Mutex
}

var GlobalConfigData = ConfigData{
	MockData: make(map[string]map[string]MockData),
	MockApi:  make(map[string]*MockApiInfoData),
}

type ConfigHandler interface {
	SupportExt(fi fs.FileInfo) bool

	Handle(data *ConfigData, path string) bool
}

// EmptyConfigHandler help print not support file ext info
type EmptyConfigHandler struct {
}

func (handler EmptyConfigHandler) Handle(_ *ConfigData, path string) bool {
	logger.W("file: %s is not handle", path)
	return true
}

func (handler EmptyConfigHandler) SupportExt(_ fs.FileInfo) bool {
	return true
}

var configS = []ConfigHandler{ApiDataHandler{}, MockDataHandler{}, EmptyConfigHandler{}}

// HandleConfigFile handle config file
func HandleConfigFile(path string, fi fs.FileInfo) {
	for _, handler := range configS {
		if handler.SupportExt(fi) && handler.Handle(&GlobalConfigData, path) {
			break
		}
	}
}

// file路径集合 加载json文件 信息
func loadApiInfo(filePathList []string) {

	for _, p := range filePathList {
		_ = filepath.Walk(p, func(jsonPath string, info fs.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				HandleConfigFile(jsonPath, info)
			}
			return nil
		})
	}
}
