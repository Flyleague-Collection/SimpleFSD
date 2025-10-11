# Metar报文源配置

| 配置项         | 说明                              |
|:------------|:--------------------------------|
| url         | 查询地址, %s会被替换为机场ICAO码            |
| return_type | url返回类型, 返回值分为raw, html, json三种 |
| reverse     | 数据排列方式                          |
| multiline   | 数据分行方式                          |
| selector    | 数据选择器, 仅返回类型为html或json时生效       |

?> 报文源配置错误可能会导致无法正确获取metar信息

## Metar报文源配置示例

* [RawMetar源示例](/configuration/metar/metar_raw_example.md)
* [HtmlMetar源示例](/configuration/metar/metar_html_example.md)
* [JsonMetar源示例](/configuration/metar/metar_json_example.md)