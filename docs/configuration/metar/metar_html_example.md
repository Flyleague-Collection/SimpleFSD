# Metar报文源示例(html)

* `url` 查询地址
* `return_type` 固定为html: 返回一个网页
* `reverse` 数据排列方式, 此值为false时, 服务器取第一行作为metar报文, 反之取最后一行
* `selector` html css 选择器, 用于提取metar报文
* `multiline` 数据分行方式  
  当此行为空时, 如果选择器只匹配到一条记录, 则服务器认为返回值里面仅有一条metar报文  
  不为空时, 服务器按照此配置作为分隔符切分数据并按照`reverse`字段配置返回报文

举个例子  
网址如下：`https://example.com/api/metar?icao=%s&style=html`  
其中`%s`会被替换为查询机场的ICAO码

* 如果网址返回值为:

```html
<!-- 省略部分无关代码 -->
<div class="result">
    <pre>METAR ...</pre>
</div>
```

则配置如下

```json
{
  "url": "https://example.com/api/metar?icao=%s&style=html",
  "return_type": "html",
  "reverse": false,
  "selector": "div.result > pre",
  "multiline": ""
}
```

* 如果网址返回值为:

```html
<!-- 省略部分无关代码 -->
<div class="result">
  <pre>
    METAR ...
    METAR ...
  </pre>
</div>
```

且最新数据为第一行, 则配置如下

```json
{
  "url": "https://example.com/api/metar?icao=%s&style=html",
  "return_type": "html",
  "reverse": false,
  "selector": "div.result > pre",
  "multiline": "\n"
}
```

* 如果网址返回值为:

```html
<!-- 省略部分无关代码 -->
<div class="result">
  <pre>
    METAR ...
    METAR ...
  </pre>
</div>
```

且最新数据为最后一行, 则配置如下

```json
{
  "url": "https://example.com/api/metar?icao=%s&style=html",
  "return_type": "html",
  "reverse": true,
  "selector": "div.result > pre",
  "multiline": "\n"
}
```

* 如果网址返回值为:

```html
<!-- 省略部分无关代码 -->
<div class="result">
    <pre>METAR ...</pre>
    <pre>METAR ...</pre>
</div>
```

且最新数据为第一行, 则配置如下

```json
{
  "url": "https://example.com/api/metar?icao=%s&style=html",
  "return_type": "html",
  "reverse": false,
  "selector": "div.result > pre",
  "multiline": ""
}
```

* 如果网址返回值为:

```html
<!-- 省略部分无关代码 -->
<div class="result">
    <pre>METAR ...</pre>
    <pre>METAR ...</pre>
</div>
```

且最新数据为最后一行, 则配置如下

```json
{
  "url": "https://example.com/api/metar?icao=%s&style=html",
  "return_type": "html",
  "reverse": false,
  "selector": "div.result > pre",
  "multiline": ""
}
```