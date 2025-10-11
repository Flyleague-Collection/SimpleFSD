# Metar报文源示例(raw)

* `url` 查询地址
* `return_type` 固定为raw: 直接返回metar文本
* `reverse` 数据排列方式, 此值为false时, 服务器取第一行作为metar报文, 反之取最后一行
* `multiline` 数据分行方式  
  当此行为空时, 服务器认为返回值里面仅有一条metar报文  
  不为空时, 服务器按照此配置作为分隔符切分数据并按照`reverse`字段配置返回报文

举个例子  
网址如下：`https://aviationweather.gov/api/data/metar?ids=%s`  
其中`%s`会被替换为查询机场的ICAO码

* 如果网址直接返回文本类型的结果: `METAR ...` 则配置如下

```json
{
  "url": "https://aviationweather.gov/api/data/metar?ids=%s",
  "return_type": "raw",
  "reverse": false,
  "multiline": ""
}
```

* 如果网址返回的结果为: `METAR ...\nMETAR ...\nMETAR ...` 且最新数据为第一行, 则配置如下

```json
{
  "url": "https://aviationweather.gov/api/data/metar?ids=%s",
  "return_type": "raw",
  "reverse": false,
  "multiline": "\n"
}
```

* 如果网址返回的结果为: `METAR ...\nMETAR ...\nMETAR ...` 且最新数据为最后一行, 则配置如下

```json
{
  "url": "https://aviationweather.gov/api/data/metar?ids=%s",
  "return_type": "raw",
  "reverse": true,
  "multiline": "\n"
}
```