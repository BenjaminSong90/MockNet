package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

var defaultConfigInfo = "{\n  " +
		"\"proxy_host\": \"www.baidu.com\",\n  " +
		"\"proxy_scheme\": \"https\",\n  " +
		"\"address\": \":8080\",\n  " +
		"\"mock_api_path\": [\n    " +
			"\"%s\"\n  " +
		"],\n  " +
		"\"file_path\": \"%s\",\n  " +
		"\"file_watcher\": true,\n  " +
		"\"file_watcher_config\": {\n    " +
			"\"valid_ext\": \".json\",\n    " +
			"\"no_reload_ext\": \".tpl, .tmpl, .html\",\n    " +
			"\"ignored_folder\" : \"\"\n  " +
		"}" +
	"\n}"

var defaultApiInfo = "{\n" +
	"  \"api\":[\n" +
	"    {\n" +
	"      \"path\": \"/hello_wrold\",\n" +
	"      \"method\": \"GET\",\n" +
	"      \"is_restful\": true,\n" +
	"      \"data\": {\n" +
	"        \"/hello_wrold\": {\n" +
	"          \"info\": \"hello world\"\n" +
	"        }\n" +
	"      }\n" +
	"    }\n" +
	"  ]\n" +
	"}"

func CheckModuleOrCreate() {
	modulePath, err := os.Getwd()
	if err != nil {
		return
	}

	configPath := filepath.Join(modulePath, "config.json")
	apiDirPath := filepath.Join(modulePath , "api")
	apiPath := filepath.Join(modulePath , "api", "api.json")
	staticDirPath := filepath.Join(modulePath, "static")

	if Exists(configPath) {
		return
	}

	err = createAndWriter(configPath, fmt.Sprintf(defaultConfigInfo, apiDirPath, staticDirPath))
	if err != nil{
		return
	}

	if !Exists(apiDirPath){
		err = os.Mkdir(apiDirPath, os.ModePerm)
		if err == nil{
			_ = createAndWriter(apiPath, defaultApiInfo)
		}
	}


	if !Exists(staticDirPath){
		_ = os.Mkdir(staticDirPath, os.ModePerm)
	}



}

func createAndWriter(path string, config string) (err error){
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer  file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(config)

	if err != nil{
		return err
	}
	return writer.Flush()
}
