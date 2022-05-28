package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mock_net/config"
	"mock_net/middleware"
	"mock_net/router"
)

func main() {

	config.ParseProjectConfig()

	var apiInfoList = config.LoadApiInfo(config.PConfig.MockApiPath)


	if len(*apiInfoList) == 0 {
		panic(fmt.Errorf(" mock api is empty!"))
	}

	r := gin.Default()
	r.Use(middleware.NoFundHandle(*apiInfoList, config.PConfig))

	if len(config.PConfig.VideoPath) != 0{
		r.StaticFS("/videos",gin.Dir(config.PConfig.VideoPath, true))
	}


	router.InitApi(r, apiInfoList)

	address := "8080"
	if len(config.PConfig.Address) != 0 {
		address = config.PConfig.Address
	}

	r.Run(address) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
