package setting

import (
	"fmt"
	"io/fs"
	"mock_net/utils"
	"path/filepath"
)

type projectConfig struct {
	ProxyHost      string   `json:"proxy_host"`    //请求代理的host
	ProxyScheme    string   `json:"proxy_scheme"`  //请求代理的host
	Address        string   `json:"address"`       //服务端启动的address
	MockApiPath    []string `json:"mock_api_path"` //加载 mock api 信息的地址
	StaticFilePath string   `json:"file_path"`    //视频文件地址
}

var configJsonFormatInfo = `
	{
	  "proxy_host": "host",
	  "proxy_scheme": "scheme",
	  "address": ":8080",
	  "mock_api_path":[""]
	}
`

func loadProjectConfig() {
	config := projectConfig{}
	err := utils.LoadFileJson("config.json", &config)
	if err != nil {

		panic(utils.FormatPanicString(err, fmt.Sprintf("please creat config.json file in project root dir:\n%s", configJsonFormatInfo)))
	}
	setProxyHost(config.ProxyHost)
	setProxySchema(config.ProxyScheme)
	setStartAddress(config.Address)
	setLocalApiInfoPath(config.MockApiPath)
	setStaticFilePath(config.StaticFilePath)
}

type mockApiInfo struct {
	ApiInfo []mockApi `json:"api"`
}

type mockApi struct {
	Path      string                            `json:"path"`       //请求的路径
	IsRestful bool                              `json:"is_restful"` //是否restful请求
	Data      map[string]map[string]interface{} `json:"data"`       //mock返回的数据
	Method    string                            `json:"method"`     //请求方法
	KeyName   string                            `json:"key_name"`   //在非restful 的情况下用于判断mock数据
}

//file路径集合 加载json文件 信息
func loadApiInfo(filePathList []string) {

	var apiInfoList []mockApi
	for _, p := range filePathList {
		_ = filepath.Walk(p, func(jsonPath string, info fs.FileInfo, err error) error {
			if err == nil && !info.IsDir() && filepath.Ext(info.Name()) == ".json" {
				info := mockApiInfo{}

				e := utils.LoadFileJson(jsonPath, &info)
				if e == nil {
					apiInfoList = append(apiInfoList, info.ApiInfo...)
				} else {
					utils.ErrorLogger("json format error %v\n", e)
				}
			}
			return nil
		})
	}

	//根据 path 和 method来过滤 Api 信息
	filterMap := make(map[string]mockApi)
	for _, info := range apiInfoList {
		filterMap[info.Path+info.Method] = info
	}

	var result []ApiInfo
	for _, v := range filterMap {
		result = append(result, v.toApiInfo())
	}

	if len(result) == 0 {
		panic(fmt.Errorf(" mock api is empty!"))
	}

	setApiInfo(&result)
}

func (mockApi mockApi) toApiInfo() ApiInfo {
	apiInfo := ApiInfo{}

	apiInfo.Method = mockApi.Method
	apiInfo.KeyName = mockApi.KeyName
	apiInfo.Data = mockApi.Data
	apiInfo.Path = mockApi.Path
	apiInfo.Restful = mockApi.IsRestful

	return apiInfo
}
