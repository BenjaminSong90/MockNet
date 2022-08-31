package setting

type ConfigParams interface {
	ApiData | MockData
}

type ConfigHandler[T ConfigParams] interface {
	HandleExt() string

	ParserData(fullPath string) (error, T)
}

var configHandlerMap = map[string]interface{}{}
var apiConfigHandler = ApiDataHandler{}
var mockDataConfigHandler = MockDataHandler{}

func init() {
	configHandlerMap[apiConfigHandler.HandleExt()] = apiConfigHandler
	configHandlerMap[mockDataConfigHandler.HandleExt()] = mockDataConfigHandler
}
