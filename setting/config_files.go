package setting

import (
	"fmt"
	"io/fs"
	"mocknet/logger"
	"mocknet/utils"
	"path/filepath"
)

type projectConfig struct {
	ProxyHost         string            `json:"proxy_host"`          //请求代理的host
	ProxyScheme       string            `json:"proxy_scheme"`        //请求代理的host
	Address           string            `json:"address"`             //服务端启动的address
	MockApiPath       []string          `json:"mock_api_path"`       //加载 mock api 信息的地址
	StaticFilePath    string            `json:"file_path"`           //视频文件地址
	FileWatcher       bool              `json:"file_watcher"`        //是否开启文件更新刷新server
	FileWatcherConfig map[string]string `json:"file_watcher_config"` //文件变化通知配置信息
}

func loadProjectConfig() {
	config := projectConfig{}
	err := utils.LoadFileJson("config.json", &config)
	if err != nil {

		panic(logger.FormatPanicString(err, fmt.Sprintf("config.json parse is fail")))
	}
	setProxyHost(config.ProxyHost)
	setProxySchema(config.ProxyScheme)
	setStartAddress(config.Address)
	setLocalApiInfoPath(config.MockApiPath)
	setStaticFilePath(config.StaticFilePath)
	setFileWatcherOpen(config.FileWatcher)
	validExt, ok := config.FileWatcherConfig["valid_ext"]
	if ok {
		setFileWatcherValidExt(validExt)
	}

	noReloadExt, ok := config.FileWatcherConfig["no_reload_ext"]
	if ok {
		setFileWatcherNoReloadExt(noReloadExt)
	}

	ignoredFolder, ok := config.FileWatcherConfig["ignored_folder"]
	if ok {
		setFileWatcherIgnoredFolder(ignoredFolder)
	}

}

type mockApiInfo struct {
	ApiInfo []mockApi `json:"api"`
	Close   bool      `json:"close"`
}

type mockApi struct {
	Path      string                            `json:"path"`       //请求的路径
	IsRestful bool                              `json:"is_restful"` //是否restful请求
	Data      map[string]map[string]interface{} `json:"data"`       //mock返回的数据
	Method    string                            `json:"method"`     //请求方法
	KeyName   string                            `json:"key_name"`   //在非restful 的情况下用于判断mock数据
}

// file路径集合 加载json文件 信息
func loadApiInfo(filePathList []string) {

	var apiInfoList []mockApi
	for _, p := range filePathList {
		_ = filepath.Walk(p, func(jsonPath string, info fs.FileInfo, err error) error {
			if err == nil && !info.IsDir() && filepath.Ext(info.Name()) == ".json" {
				info := mockApiInfo{}

				e := utils.LoadFileJson(jsonPath, &info)
				if e == nil {
					if !info.Close {
						apiInfoList = append(apiInfoList, info.ApiInfo...)
					}
				} else {
					logger.E("json format error %v\n", e)
				}
			}
			return nil
		})
	}

	//concat api info
	concatMap := make(map[string]mockApi)
	for _, info := range apiInfoList {
		if v, ok := concatMap[info.Path+info.Method]; ok && !info.IsRestful {
			for nk, nv := range info.Data {
				v.Data[nk] = nv
			}
		} else {
			concatMap[info.Path+info.Method] = info
		}

	}

	var result []*ApiInfo
	for _, v := range concatMap {
		result = append(result, v.toApiInfo())
	}

	if len(result) == 0 {
		panic(fmt.Errorf(" mock api is empty!"))
	}

	setApiInfo(&result)
}

func (mockApi mockApi) toApiInfo() *ApiInfo {
	apiInfo := ApiInfo{}

	apiInfo.Method = mockApi.Method
	apiInfo.KeyName = mockApi.KeyName
	apiInfo.Data = mockApi.Data
	apiInfo.Path = mockApi.Path
	apiInfo.Restful = mockApi.IsRestful

	return &apiInfo
}
