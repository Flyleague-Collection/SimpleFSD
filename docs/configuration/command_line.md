# 命令行参数

## 参数总览表

| 参数名                                                             | 类型   | 默认值          | 作用                                           |
| :----------------------------------------------------------------- | :----- | :-------------- | :--------------------------------------------- |
| -help                                                              | ×      | ×               | 显示命令帮助                                   |
| [-debug](#debug)                                                   | bool   | false           | 开启调试模式                                   |
| [-config](#config)                                                 | string | "./config.json" | 配置文件路径                                   |
| [-skip_email_verification](#skip_email_verification)               | bool   | false           | 跳过邮箱验证                                   |
| [-update_config](#update_config)                                   | bool   | false           | 迁移配置文件, 迁移旧版本配置文件               |
| [-no_logs](#no_logs)                                               | bool   | false           | 禁用日志输出到文件                             |
| [-message_queue_channel_size](#message_queue_channel_size)         | int    | 128             | 内置消息队列大小                               |
| [-download_prefix](#download_prefix)                               | str    | 本仓库raw地址   | 下载前缀, 用于无法连接到github或其他情况下使用 |
| [-metar_cache_clean_interval](#metar_cache_clean_interval)         | str    | 30m             | 过期metar报文清理间隔                          |
| [-metar_query_thread](#metar_query_thread)                         | int    | 32              | metar报文查询线程数                            |
| [-fsd_record_filter](#fsd_record_filter)                           | int    | 10              | fsd连线记录数值过滤                            |
| [-vatsim](#vatsim)                                                 | bool   | false           | 对管制员登录启用VATSIM协议                     |
| [-vatsim_full](#vatsim_full)                                       | bool   | false           | 对飞行员登录启用VATSIM协议                     |
| [-mutil_thread](#mutil_thread)                                     | bool   | false           | 使用多线程处理客户端连接                       |
| [-visual_pilot](#visual_pilot)                                     | bool   | false           | 启用虚拟坐标                                   |
| [-websocket_heartbeat_interval](#websocket_heartbeat_interval)     | str    | 30s             | websocket心跳间隔                              |
| [-websocket_timeout](#websocket_timeout)                           | str    | 60s             | websocket超时时间                              |
| [-websocket_message_channel_size](#websocket_message_channel_size) | int    | 128             | websocket消息频道大小                          |

## debug

[环境变量#DEBUG_MODE](/configuration/environment.md#debug_mode)

启用FSD的调试模式  
此时FSD会输出大量调试信息方便调试  
注意此选项会影响FSD性能和产生大量日志  
此选项默认关闭

!> 不要在生产环境或非必要情况下打开

## config

[环境变量#CONFIG_FILE_PATH](/configuration/environment.md#config_file_path)

覆盖FSD默认的配置文件路径  
默认路径为`./config.json`

## skip_email_verification

[环境变量#SKIP_EMAIL_VERIFICATION](/configuration/environment.md#skip_email_verification)

让API接口跳过邮箱验证  
用于想快速测试API可用性  
但不想配置邮箱配置的用户  
此选项默认关闭

!> 不要在生产环境或非必要情况下打开

## update_config

[环境变量#UPDATE_CONFIG](/configuration/environment.md#update_config)

开启此选项后  
在配置文件版本不匹配的情况下, FSD不会直接退出  
而是会读取已经存在的配置文件并尝试进行配置文件迁移  
此选项默认关闭

!> 这个功能是实验性支持，请在迁移前备份一份配置文件

当配置文件仅出现一些配置选项的增改而不涉及已有选项的搬移的时候  
本功能是相对可靠的  
但如果配置文件出现一些已有选项的重构或者移动的时候  
本功能很可能会导致先前的配置丢失
但仅会丢失发生变化的那一块的配置, 其余配置保持不变

## no_logs

[环境变量#NO_LOGS](/configuration/environment.md#no_logs)

禁用日志输出到文件  
一般用于单元测试的时候抑制日志输出  
也可以用来防止双重记录日志  
此选项默认关闭

## message_queue_channel_size

[环境变量#MESSAGE_QUEUE_CHANNEL_SIZE](/configuration/environment.md#message_queue_channel_size)

内置消息队列缓冲区大小  
如果你有大量访问API的需求  
可以适当调大这个值  
默认大小为128

## download_prefix

[环境变量#DOWNLOAD_PREFIX](/configuration/environment.md#download_prefix)

资源文件下载前缀  
当FSD发现缺失运行文件时  
会通过此选项拼接文件路径下载文件  
默认路径为`https://raw.githubusercontent.com/Flyleague-Collection/SimpleFSD/refs/heads/main`  
如果您运行FSD的网络环境无法连接GITHUB或者访问速度过慢  
您可以设置此选项为
`https://gh-proxy.com/https://raw.githubusercontent.com/Flyleague-Collection/SimpleFSD/refs/heads/main`  
友链：[ghproxy](https://gh-proxy.com/)

## metar_cache_clean_interval

[环境变量#METAR_CACHE_CLEAN_INTERVAL](/configuration/environment.md#metar_cache_clean_interval)

METAR过期报文清理间隔  
输入值应当是一个Duration字符串  
比如: 30m(30分钟), 10s(10秒), 1h(1小时)  
默认值为30m

## metar_query_thread

[环境变量#METAR_QUERY_THREAD](/configuration/environment.md#metar_query_thread)

METAR报文查询线程数  
当客户端一次要求多个METAR报文时  
服务端并发查询数量  
默认值为32

## fsd_record_filter

[环境变量#FSD_RECORD_FILTER](/configuration/environment.md#fsd_record_filter)

FSD连线记录过滤  
本选项单位为秒  
仅当连线时长高于本选项设置的数值时  
此次联飞时长才会被记录在案  
默认值为10

## vatsim

[环境变量#VATSIM](/configuration/environment.md#vatsim)

是否对管制员登录启用VATSIM协议支持  
若开启，则当管制员登录服务器时  
必须使用VATSIM协议  
详情请看[VATSIM协议指南]()
此选项默认关闭

## vatsim_full

[环境变量#VATSIM_FULL](/configuration/environment.md#vatsim_full)

是否对飞行员登录启用VATSIM协议支持  
当[VATSIM](#VATSIM)为`false`时本项不得为`true`  
若开启，则当飞行员登录服务器时  
必须使用VATSIM协议  
详情请看[VATSIM协议指南]()
此选项默认关闭

## mutil_thread

[环境变量#MUTAR_THREAD](/configuration/environment.md#mutil_thread)

是否在处理客户端消息时使用并发  
打开此开关可以获得更好的性能  
但可能会遇到包括但不限于：ATIS行错位等问题  
除非真的遇到严重的性能问题，否则不建议开启  
此选项默认关闭

## visual_pilot

[环境变量#VISUAL_PILOT](/configuration/environment.md#visual_pilot)

是否启用虚拟飞行员坐标点支持  
如果打开且客户端支持虚拟飞行员坐标点  
则服务器会开启虚拟飞行员坐标点功能  
客户端`0.2s`上传一次当前位置  
注意：此功能只会影响到飞行员相互之间的刷新间隔，即管制员并非0.2s刷新频率  
此选项默认关闭

## websocket_heartbeat_interval

[环境变量#WEBSOCKET_HEART_INTERVAL](/configuration/environment.md#websocket_heartbeat_interval)

websocket心跳包间隔  
超过此间隔未发送消息  
服务器会主动断开连接  
输入值应当是一个Duration字符串  
比如: 30m(30分钟), 10s(10秒), 1h(1小时)  
默认值为30s

## websocket_timeout

[环境变量#WEBSOCKET_TIMEOUT](/configuration/environment.md#websocket_timeout)

websocket超时时间  
输入值应当是一个Duration字符串  
比如: 30m(30分钟), 10s(10秒), 1h(1小时)  
默认值为30s


## websocket_message_channel_size

[环境变量#WEBSOCKET_MESSAGE_CHANNEL_SIZE](/configuration/environment.md#websocket_message_channel_size)

websocket消息队列缓冲区大小  
如果你对websocket有大量写入或读取需求  
可以适当调大这个值  
默认大小为128
