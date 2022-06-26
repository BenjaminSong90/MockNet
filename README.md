# MockNet
MockNet 是对服务端的mock工具，帮助前端开发初期解决数据mock的问题，调用的api不存在的时候,会根据本地的proxy配置信息来调用远端的api


**mock_config.json**工具的配置工具
```json
	{
	  "proxy_host": "www.xxx.com",
	  "proxy_scheme": "https",
	  "address": ":8080",
	  "mock_api_path": [
		"/xx/xx/xx/api_folder"
	  ],
	  "file_path": "/xx/xx/xx/static_file_folder",
	  "file_watcher": true,
	  "file_watcher_config": {
		"valid_ext": ".json",
		"no_reload_ext": ".tpl, .tmpl, .html",
		"ignored_folder" : ""
	  }
	}
```

**api.json** 是mock api返回数据
```json
{
  "api":[
    {
      "path": "/test",
      "method": "GET",
      "data": {
        "/test": {
          "a": "aa",
          "nb": "real nb"
        }
      }
    },
    {
      "path": "/test1",
      "method": "POST",
      "is_restful": false,
      "key_name": "h",
      "data": {
        "test2": {
          "a": "aa",
          "nb": "real nb"
        }
      }
    },
    {
      "path": "/test2",
      "method": "POST",
      "is_restful": true,
      "data": {
        "/test2?t=1": {
          "a": "aa",
          "nb": "real nb"
        }
      }
    }
  ]
}
```
当api的配置信息为 restful的时候，会根据api的uri信息直接返回配置的数据，不存在的时候会返回空。
如果配置信息 不是restful的时候，会读取body中的key_name设置的值，如果key_name为空的时候，会代理请求远端服务，如果body中没有key_name设置的值，也会代理请求远端服务



#feature
- 上传文件
- 下载文件
- get/delete 更具path参数来读取数据
- 图片

