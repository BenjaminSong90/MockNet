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


如果 api 是restful请求：
    字典的key RequestURI+key_name按照顺序进行拼接： /user/1?t=1,key1V,key2V,key3V
    key1V是用key1从header和body找查找的，如果没有会直接跳过例如 上面key2没有查找到的，生成的key是：/user/1?t=1,key1V,key3V
	key的查找顺序是header 和 body,如果header和body中都存在，body中的value会覆盖header中的value

如果 api 不是 restful 请求：
	字典的key key_name按照顺序进行拼接： key1V,key2V,key3V
	key1V是用key1从header和body找查找的，如果没有会直接跳过例如 上面key2没有查找到的，生成的key是：key1V,key3V
	key的查找顺序是header 和 body,如果header和body中都存在，body中的value会覆盖header中的value
	当没有找到后也会使用restful进行查找
	
	
**api.json** 是mock api返回数据
```json
{
  "api":[
    {
          "path": "/test5",
          "method": "GET",
          "key_name": "h",
          "data": {
            "/test5": {
              "a": "aa",
              "nb": "real nb"
            }
          }
        }
  ]
}
```
上面的请求也是有可能用restful方式查找的，应为GET的body是null,如果header中没有 h 这个key,那么会用restful方式进行查找，
因为values是空，最后查找的key就是uri了,如果header中有key h那么h对应的value就是/test5，不然会用代理请求




#feature
- 上传文件
- 下载文件
- get/delete 更具path参数来读取数据
- 图片

