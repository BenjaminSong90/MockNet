package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mock_net/middleware"
	"mock_net/model"
	"mock_net/router"
	"mock_net/utils"
	"strings"
)

var mockConfig = model.ProjectConfig{}

var configJsonFormatInfo = `
	{
	  "proxy_host": "",
	  "proxy_scheme": "",
	  "address": ":",
	  "mock_api_path":[""]
	}
`

func main() {


	err := utils.LoadFileJson("config.json", &mockConfig)
	if err != nil {

		panic(fmt.Errorf("\033[0;40;31m 请在项目文件夹下创建 config.json 文件:\n%s \033[0m\n",configJsonFormatInfo))
	}

	var apiInfoList = utils.LoadJson(mockConfig.MockApiPath)

	r := gin.Default()
	r.Use(middleware.EnhanceMDW(mockConfig))
	r.Use(middleware.NoFundHandle(*apiInfoList))

	for _, apiDetail := range *apiInfoList {

		switch strings.ToUpper(apiDetail.Method) {
		case "GET":
			router.AddGetApi(r, apiDetail)
		case "POST":
			router.AddPostApi(r, apiDetail)
		case "DELETE":
			router.AddDeleteApi(r, apiDetail)
		case "PUT":
			router.AddPutApi(r, apiDetail)
		default:
			router.AddGetApi(r, apiDetail)
		}
	}
	address := "8080"
	if len(mockConfig.Address) != 0 {
		address = mockConfig.Address
	}

	r.Run(address) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

