package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func LoadFileJson(filePath string, v interface{}) (err error){
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, v)

	return
}
