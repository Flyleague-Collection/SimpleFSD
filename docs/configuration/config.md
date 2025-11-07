# 配置文件介绍

## config_version(配置文件版本)

配置文件版本，通常情况下主版本号与FSD版本号一致  
即如果FSD的版本为`x.y.z`  
那么配置文件的版本号为`x.y.a`  
如果FSD启动时发现配置文件版本不匹配  
那么FSD会在报错后直接退出  
你也可以通过添加`-update_config`命令行参数来让FSD自动迁移配置文件  
详情参见[命令行参数#update_config](/configuration/command_line.md#update_config)

## server(FSD服务配置)

此配置块控制FSD的所有服务内容  
包括FSD本体, 配套的HTTP服务器与GRPC服务器

---

### general(通用配置项)

#### simulator_server(模拟机服务器)

控制服务器是否为模拟机服务器  
由于需要检查网页提交计划与实际连线计划是否一致  
所以飞行计划存储使用用户cid进行标识  
但模拟机所有的模拟机都是一个用户cid, 此时就会出问题  
即模拟机计划错误或者无法获取到计划  
此时需要将此配置设置为`true`  
这样服务器就会使用呼号作为标识  
但是与此同时就失去了呼号匹配检查的功能  
网页提交计划仍然可用, 只是没有检查功能  
此选项默认关闭

#### bcrypt_cost(密码加密轮数)

bcrypt密码加密轮数  
该选项只影响FSD加密密码时进行的轮数  
默认与推荐值为`12`  
取值范围: `[4, 32)`

---

### fsd_server(FSD服务器配置)

#### fsd_name(FSD名称)

FSD名称, 会被发送到连接到服务器的客户端作为motd消息  
motd消息请见[首行motd格式](#first_motd_line首行motd格式)

#### host(监听地址)

FSD服务器监听地址  
一般填写0.0.0.0即可

#### port(监听端口)

FSD服务器监听端口  
默认的监听端口为`6809`  
由于EuroScope不支持更改端口  
所以不建议修改此选项

#### airport_data_file(机场数据路径)

机场数据路径, 如果不存在会自动下载  
默认值为`data/airport.json`

#### pos_update_points(历史飞行路径记录间隔)

飞行员历史飞行路径记录间隔  
间隔指的是：当客户端每发过来N个包就记录一次位置  
比如：如果这个值为1, 则服务器每个包都记录  
取值范围`(1, ∞)`

#### heartbeat_interval(心跳间隔)

FSD服务器心跳间隔, 超过此时间FSD会认为客户端已断开连接  
由于EuroScope通常情况下两个包最大间隔为25s, 加之服务器处理也需要时间  
故推荐此数值大于30s, 但必须大于25s  
默认值为`40s`  
取值范围`(25s, ∞)`

#### whazzup_cache_time(whazzup缓存时间)

whazzup缓存时间  
你可以将其设为0来禁用缓存, 此举可能会导致性能问题  
默认值为`15s`  
取值范围`[0s, ∞)`

#### session_clean_time(会话过期时间)

FSD服务器会话过期时间  
即在客户端断开后服务器不会立刻销毁会话信息  
而是会等待过期时间  
在过期时间内重连, 则会话不会被销毁, 服务器会自动匹配断开时的session  
反之则会彻底销毁会话信息  
默认值为`40s`

#### max_workers(最大工作线程数)

FSD最大工作线程数, 也可以理解为最大同时连接的客户端数目  
默认值为`128`

#### max_broadcast_workers(广播最大线程数)

最大广播线程数, 用于广播消息的最大线程数  
推荐与[最大工作线程数](#max_workers最大工作线程数)保持一致  
默认值为`128`

#### first_motd_line(首行motd格式)

首行发送到客户端的motd格式  
默认格式为`Welcome to use %[1]s v%[2]s`  
其中`%[1]s`会被替换为[FSD名称](#fsd_namefsd名称)的值  
`%[2]s`会被替换为FSD版本

#### motd(motd消息)

要发送到客户端的motd消息  
第一行为[首行motd格式](#first_motd_line首行motd格式)的值  
后续为本配置项设置的内容

#### range_limit(视程范围限制)

EuroScope视程范围限制

| 配置项              | 默认值   | 说明               |
|:-----------------|:------|:-----------------|
| refuse_out_range | false | 是否断开超出视程范围限制的客户端 |
| observer         | 300   | 观察员视程范围限制        |
| delivery         | 20    | 放行视程范围限制         |
| ground           | 20    | 地面视程范围限制         |
| tower            | 50    | 塔台视程范围限制         |
| approach         | 150   | 进近视程范围限制         |
| center           | 600   | 区域视程范围限制         |
| apron            | 20    | 机坪视程范围限制         |
| supervisor       | 300   | 监管者视程范围限制        |
| administrator    | 300   | 管理员视程范围限制        |
| fss              | 1500  | 飞服视程范围限制         |

---

### http_server(Http服务器配置)

#### enabled(启用Http服务器)

是否启用Http服务器  
默认不启用Http服务器

#### server_address(访问地址)

Http服务器的访问地址, 用于拼接传统Whazzup接口  
需要是完整的url访问路径  
比如：`http://127.0.0.1:6810` 这种是合法的  
类似 `127.0.0.1:6810` 是非法的  
本字段用于[littlenavmap]的在线航班显示  
如果你不需要[littlenavmap]的在线航班显示, 你可以忽略此字段  
关于littlenavmap如何配置在线显示, 请见[配置LittleNavMap的在线航班显示]()

#### host(监听地址)

Http服务器监听地址  
一般填写0.0.0.0即可

#### port(监听端口)

Http服务器监听端口  
默认的监听端口为`6810`  
此端口可以任意修改

#### proxy_type(代理类型)

本字段表明位于Http服务器前方的代理服务器如何向本服务器传递客户端真实IP

| 数值 | 含义                             |
|:---|:-------------------------------|
| 0  | 直连无代理服务器                       |
| 1  | 代理服务器使用Http头部`X-Forwarded-For` |
| 2  | 代理服务器使用Http头部`X-Real-Ip`       |

#### trusted_ip_range(信任的代理服务器地址)

本字段表示信任的代理服务器地址, 当[代理类型](#proxy_type代理类型)为`1`时此配置项起效  
详情请见: [https://echo.labstack.com/docs/ip-address](https://echo.labstack.com/docs/ip-address)  
简单来说这里填写的是前方代理服务器的IP地址, 否则服务器无法确定这个IP地址是用户的还是代理服务器的

?> 注意：内网地址和回环地址默认被信任, 即`127.0.0.1`或者`192.168.1.1`此类地址是默认被信任不用手动添加的

如果在服务器前方配置有CDN, 那么这里需要填写CDN节点的<span style="font-size: 1.1rem;font-weight: 600;">所有</span>
可能节点IP  
这通常可以在CDN提供商文档处查到, 具体请看[CDN配置指引](../advance_configuration/cdn.md)  
此处要求的IP地址格式为CIDR, 例如: `101.71.100.0/24`

#### body_limit(请求体大小限制)

POST请求的请求体大小限制  
将本选项设置为空字符串可以禁用大小限制  
默认值为`10MB`

---

#### store(存储引擎配置)

Http服务器文件存储引擎配置  
此处配置Http服务器接受到上传文件保存配置

##### store_type(存储引擎类型)

Http服务器存储引擎, 可选值如下

| 数值 | 含义       |
|:---|:---------|
| 0  | 本地存储     |
| 1  | 阿里云OSS存储 |
| 2  | 腾讯云COS存储 |

以*开头的字段仅当[存储引擎类型](#store_type存储引擎类型)不为`0`时此有效  
即仅当存储类型不是本地存储时有效

- *`region` 储存桶地域
- *`bucket` 储存桶名称
- *`access_id` 访问ID
- *`access_key` 访问秘钥
- *`cdn_domain` CDN访问加速域名  
  如果此配置项不为空, 则在存储引擎最后返回访问链接的时候  
  使用此配置项作为基础域名去拼接最终访问链接  
  比如：  
  一个文件的路径是`xxxxxx.png`  
  本配置项为`https://cdn.example.com`  
  那么最终的访问路径就是`https://cdn.example.com/xxxxxx.png`
- *`use_internal_url` 使用内网上传文件  
  使用内网地址上传文件, 仅阿里云OSS存储此字段有效
- `local_store_path` 本地文件保存路径  
  文件保存的本地路径
- *`remote_store_path` 远程文件保存路径  
  文件保存的远程路径

##### file_limit(文件限制)

- `image_limit` 图片文件限制
    - `max_file_size` 允许的最大文件大小, 单位是B
    - `allowed_file_ext` 允许的文件后缀名列表
    - `store_prefix`  
      存储路径前缀, 会拼接在[本地文件保存路径](#local_store_path本地文件保存路径)
      和[远程文件保存路径](#remote_store_path远程文件保存路径)后面  
      如果此项为`xxx`, [本地文件保存路径](#local_store_path本地文件保存路径)为
      `aaa`, [远程文件保存路径](#remote_store_path远程文件保存路径)为`bbb`, 文件名为`ccc.png`
      那么最终的文件路径就是
        - 本地路径`aaa/xxx/ccc.png`
        - 远程路径`bbb/xxx/ccc.png`
    - `store_in_server`  
      是否在本地也保存一份, 当[存储引擎类型](#store_type存储引擎类型)为`0`时此字段必须为true
- `file_limit` 文本文件限制  
  配置项同`image_limit(图片文件限制)`

#### rate_limit&rate_limit_window(API访问速率限制)

每个IP对每个接口单独计算限制  
采用滑动窗口限速  
即访问限制为 `rate_limit次每rate_limit_window`  
`rate_limit`默认值为`15`  
`rate_limit_window`默认值为`1m`  
所以默认的访问速率为`15次每分钟`

#### email(邮箱配置)

- `host` SMTP服务器地址
- `port` SMTP服务器端口
- `username` 发信账号
- `password` 发信密码
- `verify_expired_time` 邮箱验证码过期时间  
  默认值为`5m`
- `send_interval` 验证码发送间隔  
  两次验证码发送间隔  
  默认值为`1m`

##### template(邮件模板)

!> `verify_code_email`邮件无法被关闭  
如果真的想关闭  
请使用[命令行参数#skip_email_verification](/configuration/command_line.md#skip_email_verification)  
注意这个参数会关闭整个邮件服务

- `verify_code_email` 验证码邮件模板
    - `file_path` 模板文件路径
    - `email_title` 邮件标题
    - `enable` 是否启用该邮件
- `atc_rating_change_email` 管制员管制权限变更邮件模板
    - 配置项同验证码邮件模板
- `permission_change_email` 管理权限变更邮件模板
    - 配置项同验证码邮件模板
- `kicked_from_server_email` 踢出服务器通知邮件模板
    - 配置项同验证码邮件模板
- `password_change_email` 密码修改通知邮件模板
    - 配置项同验证码邮件模板
- `application_passed_email` 管制员申请通过通知邮件模板
    - 配置项同验证码邮件模板
- `application_rejected_email` 管制员申请拒绝通知邮件模板
    - 配置项同验证码邮件模板
- `application_processing_email` 管制员申请进度通知邮件模板
    - 配置项同验证码邮件模板
- `ticket_reply_email` 工单回复通知邮件模板
    - 配置项同验证码邮件模板

#### jwt(JWT配置)

#### secret(加密秘钥)

JWT对称加密秘钥  
请<span style="font-size: 1.25rem;font-weight: 600;color: red;">一定</span>要保护好这个秘钥, 并确保不被任何不信任的人知道

!> 如果该秘钥泄露, <span style="font-size: 1.5rem;font-weight: 600;color: red;">任何人</span>都可以伪造管理员用户

?> 更安全的做法是将本字段置空, 这样每次服务器重启都会使之前签发的所有秘钥全部失效  
~~&nbsp;&nbsp;只要连我自己都不知道秘钥是什么, 那就没人知道秘钥是什么&nbsp;&nbsp;~~  (bushi)

#### expires_time(主密钥过期时间)

JWT主密钥过期时间  
默认值为`15m`

?> 过期时间建议不要大于1小时, JWT秘钥是无状态的, 如果主密钥过期时间太长可能会导致安全问题

#### refresh_time(刷新秘钥过期时间)

JWT刷新秘钥过期时间  
默认值为`24h`
该时间是在[主密钥过期时间](#expires_time主密钥过期时间)之后的时间  
比如两者都是`1h`, 那么刷新秘钥的过期时间就是`2h`  
因为不可能你刷新秘钥比主密钥过期还早 :(

#### ssl(SSL配置)

- `enable` 是否启用SSL
- `enable_hsts` 是否启用HSTS
- `hsts_expired_time` HSTS过期时间(s)
- `include_domain` HSTS是否包括子域名

  !> 警告：如果你的其他子域名没有全部部署SSL证书  
  打开此开关可能导致没有SSL证书的域名无法访问  
  如果不懂请不要打开此开关

- `cert_file` SSL证书文件路径
- `key_file` SSL私钥文件路径

#### navigraph(Navigraph秘钥)

- `enable` 是否启用
- `token` 刷新秘钥

用于对外提供航图查询代理  
详情请看[Navigraph航图代理](../advance_configuration/navigraph.md)

### voice_server(语音服务器配置)

- `enabled` 是否启用语音服务器
- `tcp_host` 语音服务器TCP监听地址
- `tcp_port` 语音服务器TCP监听端口
- `udp_host` 语音服务器UDP监听地址
- `udp_port` 语音服务器UDP监听端口
- `timeout_interval` 语音服务器心跳包超时时间
- `max_data_size` 语音包最大大小
- `broadcast_limit` 广播限制
- `udp_packet_limit` UDP包数量限制
- `tcp_packet_limit` TCP包数量限制

### grpc_server(GRPC服务器配置)

?> 暂未开发, 配置文件无参考性

## metar_source(Metar报文源)

本配置项为列表, 列表项所有可能的配置项如下表

| 配置项         | 说明                              |
|:------------|:--------------------------------|
| url         | 查询地址, %s会被替换为机场ICAO码            |
| return_type | url返回类型, 返回值分为raw, html, json三种 |
| reverse     | 数据排列方式                          |
| multiline   | 数据分行方式                          |
| selector    | 数据选择器, 仅返回类型为html或json时生效       |

详细配置介绍和配置示例请看[Metar报文源配置](/configuration/metar/metar_source.md)

## database(数据库配置)

### type(数据库类型)

支持的数据库类型: `mysql`, `postgres`, `sqlite3`  
默认数据库类型为`sqlite3`

### database

当数据库类型为`sqlite3`的时候, 这里是数据库存放路径和文件名  
反之则为要使用的数据库名称

### host(数据库地址)

### port(数据库端口)

### username(数据库用户名)

### password(数据库密码)

### enable_ssl(是否启用SSL)

### connect_idle_timeout(连接超时时间)

数据库连接池连接超时时间

### connect_timeout(查询超时时间)

数据库查询超时时间

### server_max_connections(最大连接数)

数据库最大连接数  
这个数字请求改为你实际的数据库配置

## rating(权限配置)

你可以通过配置文件覆写管制权限与管制席位对照表  
注意!!! 这个字段会`覆盖`默认的对照表  
所以在明确的知道你在做什么之前, 不要修改这个配置  
配置文件字段为`rating`

```json5
{
  // 特殊权限配置
  "rating": {
    // 键为想要修改的权限识别名的权限值
    // 比如我想让Normal也可以上OBS席位, 也就是普通飞行员也可以以OBS身份连线
    // Normal的权限值是0, 那我的键就是0
    // 值为想要许可连线的席位的席位编码之和
    // 比如我想让飞行员可以正常连线, 也可以以OBS连线
    // 那么值就是 128 + 1 = 129
    // 如果我想他还能上个飞服(请勿模仿)
    // 那么值就是 128 + 1 + 2 = 131
    // 其他权限的对照表保持为默认
    // 你也可以将某个权限的值写为0来禁止使用该权限登录fsd
    "0": 3
  }
}
```

## facility(席位配置)

你可以通过配置文件覆写呼号后缀与管制席位对照表  
注意!!! 这个字段会`覆盖`默认的对照表  
默认的席位对照表如下

| 席位后缀名 | 允许的席位 |
|:------|:------|
| ADM   | ADM   |
| SUP   | SUP   |
| OBS   | OBS   |
| DEL   | DEL   |
| RMP   | RMP   |
| GND   | GND   |
| TWR   | TWR   |
| APP   | APP   |
| CTR   | CTR   |
| FSS   | FSS   |
| ATIS  | TWR   |

在明确的知道你在做什么之前, 不要修改这个配置  
配置文件字段为`facility`

## 配置文件示例

```json
{
  "config_version": "0.8.2",
  "server": {
    "general": {
      "simulator_server": false,
      "bcrypt_cost": 12
    },
    "fsd_server": {
      "fsd_name": "Simple-Fsd",
      "host": "0.0.0.0",
      "port": 6809,
      "airport_data_file": "data/airport.json",
      "pos_update_points": 1,
      "heartbeat_interval": "40s",
      "whazzup_cache_time": "15s",
      "session_clean_time": "40s",
      "max_workers": 128,
      "max_broadcast_workers": 128,
      "first_motd_line": "Welcome to use %[1]s v%[2]s",
      "range_limit": {
        "refuse_out_range": false,
        "observer": 300,
        "delivery": 20,
        "ground": 20,
        "tower": 50,
        "approach": 150,
        "center": 600,
        "apron": 20,
        "supervisor": 300,
        "administrator": 300,
        "fss": 1500
      },
      "motd": [
        "This is my test fsd server"
      ]
    },
    "http_server": {
      "enabled": false,
      "server_address": "http://127.0.0.1:6810",
      "host": "0.0.0.0",
      "port": 6810,
      "proxy_type": 0,
      "trusted_ip_range": [],
      "body_limit": "10MB",
      "store": {
        "store_type": 0,
        "region": "",
        "bucket": "",
        "access_id": "",
        "access_key": "",
        "cdn_domain": "",
        "use_internal_url": false,
        "local_store_path": "uploads",
        "remote_store_path": "fsd",
        "file_limit": {
          "image_limit": {
            "max_file_size": 5242880,
            "allowed_file_ext": [
              ".jpg",
              ".png",
              ".bmp",
              ".jpeg"
            ],
            "store_prefix": "images",
            "store_in_server": true
          },
          "file_limit": {
            "max_file_size": 10485760,
            "allowed_file_ext": [
              ".md",
              ".txt",
              ".pdf",
              ".doc",
              ".docx"
            ],
            "store_prefix": "files",
            "store_in_server": false
          }
        }
      },
      "limits": {
        "rate_limit": 15,
        "rate_limit_window": "1m"
      },
      "email": {
        "host": "smtp.example.com",
        "port": 465,
        "username": "noreply@example.cn",
        "password": "123456",
        "verify_expired_time": "5m",
        "send_interval": "1m",
        "template": {
          "verify_code_email": {
            "file_path": "template/email_verify.template",
            "email_title": "邮箱验证码",
            "enable": true
          },
          "atc_rating_change_email": {
            "file_path": "template/atc_rating_change.template",
            "email_title": "管制权限变更通知",
            "enable": true
          },
          "permission_change_email": {
            "file_path": "template/permission_change.template",
            "email_title": "管理权限变更通知",
            "enable": true
          },
          "kicked_from_server_email": {
            "file_path": "template/kicked_from_server.template",
            "email_title": "踢出服务器通知",
            "enable": true
          },
          "password_change_email": {
            "file_path": "template/password_change.template",
            "email_title": "飞控密码更改通知",
            "enable": true
          },
          "application_passed_email": {
            "file_path": "template/application_passed.template",
            "email_title": "管制员申请通过",
            "enable": true
          },
          "application_rejected_email": {
            "file_path": "template/application_rejected.template",
            "email_title": "管制员申请被拒",
            "enable": true
          },
          "application_processing_email": {
            "file_path": "template/application_processing.template",
            "email_title": "管制员申请进度通知",
            "enable": true
          },
          "ticket_reply_email": {
            "file_path": "template/ticket_reply.template",
            "email_title": "工单回复通知",
            "enable": true
          }
        }
      },
      "jwt": {
        "secret": "123456",
        "expires_time": "15m",
        "refresh_time": "24h"
      },
      "ssl": {
        "enable": false,
        "enable_hsts": false,
        "hsts_expired_time": 5184000,
        "include_domain": false,
        "cert_file": "",
        "key_file": ""
      }
    },
    "voice_server": {
      "enabled": true,
      "tcp_host": "0.0.0.0",
      "tcp_port": 6808,
      "udp_host": "0.0.0.0",
      "udp_port": 6807,
      "timeout_interval": "30s",
      "max_data_size": 1048576,
      "broadcast_limit": 128,
      "udp_packet_limit": 8192,
      "tcp_packet_limit": 32
    },
    "grpc_server": {
      "enabled": false,
      "host": "0.0.0.0",
      "port": 6811,
      "whazzup_cache_time": "15s"
    }
  },
  "metar_source": [
    {
      "url": "https://aviationweather.gov/api/data/metar?ids=%s",
      "return_type": "raw",
      "reverse": false,
      "multiline": ""
    },
    {
      "url": "https://example.com/api/metar?icao=%s&style=html",
      "return_type": "html",
      "reverse": false,
      "selector": "body > div > pre",
      "multiline": "\n"
    },
    {
      "url": "https://example.com/api/metar?icao=%s&style=json",
      "return_type": "json",
      "selector": "$.data.metar"
    }
  ],
  "database": {
    "type": "mysql",
    "database": "go-fsd",
    "host": "localhost",
    "port": 3306,
    "username": "root",
    "password": "123456",
    "enable_ssl": false,
    "connect_idle_timeout": "1h",
    "connect_timeout": "5s",
    "server_max_connections": 32
  },
  "rating": {},
  "facility": {}
}
```

[littlenavmap]: https://albar965.github.io/littlenavmap.html