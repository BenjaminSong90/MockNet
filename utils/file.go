package utils

import (
	"encoding/json"
	"os"
)

// LoadFileJson load local json
func LoadFileJson(filePath string, v interface{}) (err error) {

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, v)

	return
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
