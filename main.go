package main

import (
	"github.com/gin-gonic/gin"
	"mock_net/middleware"
	"mock_net/router"
	"mock_net/setting"
)

func main() {

	setting.LoadLocalConfig()

	r := gin.Default()
	r.Use(middleware.NoFundHandle(setting.GetApiInfo()))

	if len(setting.GetStaticFilePath()) != 0 {
		r.StaticFS("/file", gin.Dir(setting.GetStaticFilePath(), true))
	}

	router.InitApi(r, setting.GetApiInfo())

	r.Run(setting.GetStartAddress()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
