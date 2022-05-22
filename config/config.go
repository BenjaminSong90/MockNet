package config

import (
	"fmt"
	"io/fs"
	"mock_net/utils"
	"path/filepath"
)

type MockApiInfo struct {
	ApiInfo []ApiInfo `json:"api"`
}

type ApiInfo struct {
	Path      string                            `json:"path"`       //请求的路径
	IsRestful bool                              `json:"is_restful"` //是否restful请求
	Data      map[string]map[string]interface{} `json:"data"`       //mock返回的数据
	Method    string                            `json:"method"`     //请求方法
	KeyName   string                            `json:"key_name"`   //在非restful 的情况下用于判断mock数据
}

type ProjectConfig struct {
	ProxyHost   string   `json:"proxy_host"`    //请求代理的host
	ProxyScheme string   `json:"proxy_scheme"`  //请求代理的host
	Address     string   `json:"address"`       //服务端启动的address
	MockApiPath []string `json:"mock_api_path"` //加载 mock api 信息的地址

}

var configJsonFormatInfo = `
	{
	  "proxy_host": "host",
	  "proxy_scheme": "scheme",
	  "address": ":8080",
	  "mock_api_path":[""]
	}
`

var PConfig = ProjectConfig{}

func LoadProjectConfig() {
	err := utils.LoadFileJson("config.json", &PConfig)
	if err != nil {
		panic(fmt.Errorf("\033[0;40;31m please creat config.json file in project root dir:\n%s \033[0m\n", configJsonFormatInfo))
	}
}

//file路径集合 加载json文件 信息
func LoadConfigJson(filePathList []string) *[]ApiInfo {

	var apiInfoList []ApiInfo
	for _, p := range filePathList {
		_ = filepath.Walk(p, func(jsonPath string, info fs.FileInfo, err error) error {
			if err == nil && !info.IsDir() && filepath.Ext(info.Name()) == ".json" {
				info := MockApiInfo{}

				e := utils.LoadFileJson(jsonPath, &info)
				if e == nil {
					apiInfoList = append(apiInfoList, info.ApiInfo...)
				} else {
					fmt.Printf("\033[0;40;31m json format error %v\033[0m\n", e)
				}
			}
			return nil
		})
	}

	//根据 path 和 method来过滤 Api 信息
	filterMap := make(map[string]ApiInfo)
	for _, info := range apiInfoList {
		filterMap[info.Path+info.Method] = info
	}

	var result []ApiInfo
	for _, v := range filterMap {
		result = append(result, v)
	}

	return &result
}
