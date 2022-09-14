
## mock-net 项目配置文件
项目的config文件是和运行文件放在一起的，名称是config.json

```json
{
  "proxy_host": "www.baidu.com",//代理访问host
  "proxy_scheme": "https",//代理访问使用的schema
  "address": ":8080",//本地mock服务启用的地址
  "mock_api_path": "/Users/songlang/Desktop/api_config",//api文件存放的路径
  "file_path": "/Users/songlang/Desktop/video/",//静态文件放的路径
  "file_watcher": true,//是否开启file_watcher，开启后当文件变化的时候会重新启动服务来更新配置信息
  "file_watcher_config": {
    "valid_ext": ".data, .api",
    "no_reload_ext": ".tpl, .tmpl, .html",
    "ignored_folder" : ""
  }
}
```


## API 配置信息
api入口的配置文件，配置文件的拓展名是 .api, 例如 a接口的配置文件是 a.api
```json
{
	"path":"/name/info/:age",
	"method":"POST",
	"need_redirect":true,//默认为true,如果是false本地没有配置数据会直接返回404
	"body_key":"",
	"query_key":[],
	"plugin":"xxx.js",
	"data":{//这不是必要配置，如果这个api下面只有这个一个结果可以在这里配置
		"path":"/name/info/:age",//这里的path 会被api的path替换掉
		"method":"GET",
		"key":"param+query+funcode",
		"data":{
			"say":"hello, world!"
		}
	}
}
```


## Data 数据配置信息
data数据的配置文件，配置文件的拓展名是 .data, 例如 a接口的配置文件是 a.data
```json
{
	"path":"/name/info/:age",
	"method":"POST",
	"key":"param+query+funcode",
	"close":false,
	"data":{}
}
```


```
｜projectFolder
    |——api.api
        |--data
            |
            |--a_data.data
        
```




#feature
- 上传文件
- 下载文件
- get/delete 更具path参数来读取数据
- 图片

