package model

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
