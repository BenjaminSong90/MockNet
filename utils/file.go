package utils

import (
	"encoding/json"
	"io/ioutil"
)

//load local json
func LoadFileJson(filePath string, v interface{}) (err error) {

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, v)

	return
}
