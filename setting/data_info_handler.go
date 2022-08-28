package setting

import "fmt"

type MockData struct {
	Path string      `json:"path"` //数据对应path
	Key  string      `json:"key"`  //param+query+funcode生成的key
	Data interface{} `json:"data"` //mock 返回的数据
}

func (m MockData) String() string {
	return fmt.Sprintf("{Path: %s, Key: %s, Data: %v}", m.Path, m.Key, m.Data)
}
