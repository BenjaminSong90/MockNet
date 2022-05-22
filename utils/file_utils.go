package utils

import (
	"encoding/json"
	"io/ioutil"
)

//加载json文件 信息
func LoadFileJson(filePath string, v interface{}) (err error) {

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, v)

	return
}
