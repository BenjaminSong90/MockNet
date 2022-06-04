package setting

type Setting map[string]interface{}

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

type ApiInfo struct {
	Method        string
	Path          string
	KeyName       string
	Restful       bool
	LocalFilePath string
	Data          map[string]map[string]interface{}
}

func (apiInfo ApiInfo) isDataApi() bool {
	if len(apiInfo.LocalFilePath) == 0 {
		return true
	}
	return false
}

var globalSetting Setting = map[string]interface{}{
	"proxy_host":          "",
	"proxy_schema":        "https",
	"address":             ":8080",
	"local_api_info_path": []string{"."},
	"static_path":         ".",
	"api_info":            []ApiInfo{},
}

func GetProxyHost() string {
	return globalSetting.getString("proxy_host", "")
}
func setProxyHost(v string) {
	globalSetting.setKVOrDefault("proxy_host", v)
}

func GetProxySchema() string {
	return globalSetting.getString("proxy_schema", "https")
}

func setProxySchema(v string) {
	globalSetting.setKVOrDefault("proxy_schema", v)
}

func GetStartAddress() string {
	return globalSetting.getString("address", ":8080")
}

func setStartAddress(v string) {
	globalSetting.setKVOrDefault("address", v)
}

func GetLocalApiInfoPath() []string {
	v := globalSetting["local_api_info_path"]
	localApiInfoPath, ok := v.([]string)
	if ok {
		return localApiInfoPath
	} else {
		return []string{"."}
	}
}

func setLocalApiInfoPath(v []string) {
	globalSetting.setKVOrDefault("local_api_info_path", v)
}

func GetStaticFilePath() string {
	return globalSetting.getString("static_path", ".")
}

func setStaticFilePath(v string) {
	globalSetting.setKVOrDefault("static_path", v)
}

func GetApiInfo() *[]ApiInfo {
	ai := globalSetting["api_info"]
	apiInfo, ok := ai.([]ApiInfo)
	if ok {
		return &apiInfo
	} else {
		return &[]ApiInfo{}
	}
}

func setApiInfo(v *[]ApiInfo) {
	if v != nil && len(*v) != 0{
		globalSetting["api_info"] = *v
	} else {
		globalSetting["api_info"] = []ApiInfo{}
	}
}

func CheckProxyInfo() bool {
	return len(GetProxyHost()) != 0 &&
		len(GetProxySchema()) != 0
}

func LoadProjectConfig() {
	loadProjectConfig()

}

func LoadApiInfo()  {
	loadApiInfo(GetLocalApiInfoPath())
}
