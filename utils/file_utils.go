package utils

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"mock_net/model"
	"path/filepath"
)
//file路径集合 加载json文件 信息
func LoadJson(filePathList []string) *[]model.ApiInfo {

	var apiInfoList []model.ApiInfo

	for _, p := range filePathList {
		_ = filepath.Walk(p, func(jsonPath string, info fs.FileInfo, err error) error {
			if err == nil && !info.IsDir() && filepath.Ext(info.Name()) == ".json" {
				info := model.MockApiInfo{}
				e := LoadFileJson(jsonPath, &info)
				if e == nil {
					apiInfoList = append(apiInfoList, info.ApiInfo...)
				} else {
					fmt.Printf("\033[0;40;31m json format err %v\033[0m\n",e)
				}
			}
			return nil
		})
	}
	//根据 path 和 method来过滤 Api 信息
	filterMap := make(map[string]model.ApiInfo)
	for _, info := range apiInfoList {
		filterMap[info.Path + info.Method] = info
	}

	var result []model.ApiInfo
	for _, v := range filterMap {
		result = append(result, v)
	}

	return &result
}

//加载json文件 信息
func LoadFileJson(filePath string, v interface{}) (err error) {

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, v)

	return
}
