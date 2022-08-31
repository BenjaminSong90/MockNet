package setting

import (
	"fmt"
	"mocknet/utils"
)

type MockData struct {
	Path string      `json:"path"` //数据对应path
	Key  string      `json:"key"`  //param+query+funcode生成的key
	Data interface{} `json:"data"` //mock 返回的数据
}

func (m MockData) String() string {
	return fmt.Sprintf("{Path: %s, Key: %s, Data: %v}", m.Path, m.Key, m.Data)
}

type MockDataHandler struct {
}

func (collector MockDataHandler) HandleExt() string {
	return ".data"
}

func (collector MockDataHandler) ParserData(fullPath string) (error, MockData) {
	data := MockData{}
	err := utils.LoadFileJson(fullPath, data)
	return err, data
}

var _ ConfigHandler[MockData] = MockDataHandler{}
