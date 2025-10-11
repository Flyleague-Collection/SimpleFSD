# Metar报文源示例(json)

* `url` 查询地址
* `return_type` 固定为json: 返回json格式的数据
* `reverse` 数据排列方式, 此值为false时, 服务器取第一行作为metar报文, 反之取最后一行
* `selector` jsonpath字符串, 用于提取json中的metar报文
* `multiline` 数据分行方式

在此模式下, 服务器行为会有一点不一样

* 如果选择到的值为字符串:
    * 当`multiline`为空时, 则服务器认为返回值里面仅有一条metar报文
    * 当`multiline`不为空时, 服务器按照此配置作为分隔符切分数据并按照`reverse`字段配置返回报文
* 如果选择到的值为字符串数组:
    * 按照`reverse`字段配置返回报文
* 如果都不是则报错

举个例子  
网址如下：`https://example.com/api/metar?icao=%s&style=json`  
其中`%s`会被替换为查询机场的ICAO码

* 如果服务器返回的值为

```json
{
  "data": {
    "metar": "METAR..."
  }
}
```

则配置如下

```json
{
  "url": "https://example.com/api/metar?icao=%s&style=json",
  "return_type": "json",
  "selector": "$.data.metar",
  "reverse": false,
  "multiline": ""
}
```

* 如果服务器返回的值为

```json
{
  "data": {
    "metar": "METAR...\nMETAR..."
  }
}
```

且最新数据为第一行, 则配置如下

```json
{
  "url": "https://example.com/api/metar?icao=%s&style=json",
  "return_type": "json",
  "selector": "$.data.metar",
  "reverse": false,
  "multiline": "\n"
}
```

* 如果服务器返回的值为

```json
{
  "data": {
    "metar": "METAR...\nMETAR..."
  }
}
```

且最新数据为最后一行, 则配置如下

```json
{
  "url": "https://example.com/api/metar?icao=%s&style=json",
  "return_type": "json",
  "selector": "$.data.metar",
  "reverse": true,
  "multiline": "\n"
}
```

* 如果服务器返回的值为

```json
{
  "data": {
    "metar": [
      "METAR ...",
      "METAR ..."
    ]
  }
}
```

且最新数据为第一行, 则配置如下

```json
{
  "url": "https://example.com/api/metar?icao=%s&style=json",
  "return_type": "json",
  "selector": "$.data.metar",
  "reverse": false,
  "multiline": ""
}
```

* 如果服务器返回的值为

```json
{
  "data": {
    "metar": [
      "METAR ...",
      "METAR ..."
    ]
  }
}
```

且最新数据为最后一行, 则配置如下

```json
{
  "url": "https://example.com/api/metar?icao=%s&style=json",
  "return_type": "json",
  "selector": "$.data.metar",
  "reverse": true,
  "multiline": ""
}
```

* 如果服务器返回的值为

```json
{
  "data": {
    "metar": [
      "METAR ...",
      "METAR ...",
      "METAR ..."
    ]
  }
}
```

且最新数据为第二行, 则配置如下

```json
{
  "url": "https://example.com/api/metar?icao=%s&style=json",
  "return_type": "json",
  "selector": "$.data.metar.[1]",
  "reverse": false,
  "multiline": ""
}
```