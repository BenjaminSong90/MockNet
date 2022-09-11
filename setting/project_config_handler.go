package setting

import (
	"fmt"
	"mocknet/logger"
	"mocknet/utils"
)

type Setting map[string]interface{}

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

func (setting Setting) setKVOrDefault(k string, v interface{}) {
	switch v.(type) {
	case string:
		value, ok := v.(string)
		if ok && len(value) != 0 {
			setting[k] = value
		}
	case []string:
		value, ok := v.([]string)
		if ok && len(value) != 0 {
			setting[k] = value
		}
	default:
		setting[k] = v
	}
}

func (setting Setting) getString(k string, defaultValue string) string {
	v := setting[k]
	str, ok := v.(string)
	if ok {
		return str
	} else {
		return defaultValue
	}
}

// 项目设置的模版
var projectSetting Setting = map[string]interface{}{
	"proxy_host":                  "",
	"proxy_schema":                "https",
	"address":                     ":8080",
	"local_api_info_path":         []string{"."},
	"static_path":                 ".",
	"file_watcher_open":           false,
	"file_watcher_valid_ext":      ".json",
	"file_watcher_no_reload_ext":  ".tpl, .tmpl, .html",
	"file_watcher_ignored_folder": "",
}

func GetProxyHost() string {
	return projectSetting.getString("proxy_host", "")
}
func setProxyHost(v string) {
	projectSetting.setKVOrDefault("proxy_host", v)
}

func GetProxySchema() string {
	return projectSetting.getString("proxy_schema", "https")
}

func setProxySchema(v string) {
	projectSetting.setKVOrDefault("proxy_schema", v)
}

func GetStartAddress() string {
	return projectSetting.getString("address", ":8080")
}

func setStartAddress(v string) {
	projectSetting.setKVOrDefault("address", v)
}

func GetLocalApiInfoPath() []string {
	v := projectSetting["local_api_info_path"]
	localApiInfoPath, ok := v.([]string)
	if ok {
		return localApiInfoPath
	} else {
		return []string{"."}
	}
}

func setLocalApiInfoPath(v []string) {
	projectSetting.setKVOrDefault("local_api_info_path", v)
}

func GetStaticFilePath() string {
	return projectSetting.getString("static_path", ".")
}

func setStaticFilePath(v string) {
	projectSetting.setKVOrDefault("static_path", v)
}

func setFileWatcherOpen(isOpen bool) {
	projectSetting["file_watcher_open"] = isOpen
}

func IsFileWatcherOpen() bool {
	info := projectSetting["file_watcher_open"]
	isOpen, ok := info.(bool)
	if ok {
		return isOpen
	} else {
		return false
	}
}

func setFileWatcherValidExt(v string) {
	projectSetting.setKVOrDefault("file_watcher_valid_ext", v)
}

func GetFileWatcherValidExt() string {
	return projectSetting.getString("file_watcher_valid_ext", ".json")
}

func setFileWatcherNoReloadExt(v string) {
	projectSetting.setKVOrDefault("file_watcher_no_reload_ext", v)
}

func GetFileWatcherNoReloadExt() string {
	return projectSetting.getString("file_watcher_no_reload_ext", ".tpl, .tmpl, .html")
}

func setFileWatcherIgnoredFolder(v string) {
	projectSetting.setKVOrDefault("file_watcher_ignored_folder", v)
}

func GetFileWatcherIgnoredFolder() string {
	return projectSetting.getString("file_watcher_ignored_folder", "")
}

func CheckProxyInfo() bool {
	return len(GetProxyHost()) != 0 &&
		len(GetProxySchema()) != 0
}

func LoadProjectConfig() {
	loadProjectConfig()

}

func LoadApiInfo() {
	loadApiInfo(GetLocalApiInfoPath())
}
