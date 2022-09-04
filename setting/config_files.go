package setting

import (
	"fmt"
	"mocknet/logger"
	"mocknet/utils"
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
