# 环境变量

!> 注意：环境变量设置会覆盖对应的[命令行参数](/configuration/command_line.md)

## 环境变量总览表

| 参数名                                                       | 类型     | 默认值             | 作用                          |
|:----------------------------------------------------------|:-------|:----------------|:----------------------------|
| [DEBUG_MODE](#DEBUG_MODE)                                 | bool   | false           | 开启调试模式                      |
| [CONFIG_FILE_PATH](#CONFIG_FILE_PATH)                     | string | "./config.json" | 配置文件路径                      |
| [SKIP_EMAIL_VERIFICATION](#SKIP_EMAIL_VERIFICATION)       | bool   | false           | 跳过邮箱验证                      |
| [UPDATE_CONFIG](#UPDATE_CONFIG)                           | bool   | false           | 迁移配置文件, 迁移旧版本配置文件           |
| [NO_LOGS](#NO_LOGS)                                       | bool   | false           | 禁用日志输出到文件                   |
| [MESSAGE_QUEUE_CHANNEL_SIZE](#MESSAGE_QUEUE_CHANNEL_SIZE) | int    | 128             | 内置消息队列大小                    |
| [DOWNLOAD_PREFIX](#DOWNLOAD_PREFIX)                       | str    | 本仓库raw地址        | 下载前缀, 用于无法连接到github或其他情况下使用 |
| [MESSAGE_QUEUE_CHANNEL_SIZE](#MESSAGE_QUEUE_CHANNEL_SIZE) | int    | 128             | 内置消息队列大小                    |
| [METAR_CACHE_CLEAN_INTERVAL](#METAR_CACHE_CLEAN_INTERVAL) | str    | 30s             | 过期metar报文清理间隔               |
| [METAR_QUERY_THREAD](#METAR_QUERY_THREAD)                 | int    | 32              | metar报文查询线程数                |
| [FSD_RECORD_FILTER](#FSD_RECORD_FILTER)                   | int    | 10              | fsd连线记录数值过滤                 |

## DEBUG_MODE

注意本环境变量会覆盖[命令行参数#debug_mode](/configuration/command_line.md#debug)

启用FSD的调试模式  
此时FSD会输出大量调试信息方便调试  
注意此选项会影响FSD性能和产生大量日志  
此选项默认关闭

!> 不要在生产环境或非必要情况下打开

## CONFIG_FILE_PATH

注意本环境变量会覆盖[命令行参数#config_file_path](/configuration/command_line.md#config)

覆盖FSD默认的配置文件路径  
默认路径为`./config.json`

## SKIP_EMAIL_VERIFICATION

注意本环境变量会覆盖[命令行参数#skip_email_verification](/configuration/command_line.md#skip_email_verification)

让API接口跳过邮箱验证  
用于想快速测试API可用性  
但不想配置邮箱配置的用户  
此选项默认关闭

!> 不要在生产环境或非必要情况下打开

## UPDATE_CONFIG

注意本环境变量会覆盖[命令行参数#update_config](/configuration/command_line.md#update_config)

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

## NO_LOGS

注意本环境变量会覆盖[命令行参数#no_logs](/configuration/command_line.md#no_logs)

禁用日志输出到文件  
一般用于单元测试的时候抑制日志输出  
也可以用来防止双重记录日志  
此选项默认关闭

## MESSAGE_QUEUE_CHANNEL_SIZE

注意本环境变量会覆盖[命令行参数#message_queue_channel_size](/configuration/command_line.md#message_queue_channel_size)

内置消息队列缓冲区大小  
如果你有大量访问API的需求  
可以适当调大这个值  
默认大小为128

## DOWNLOAD_PREFIX

注意本环境变量会覆盖[命令行参数#download_prefix](/configuration/command_line.md#download_prefix)

资源文件下载前缀  
当FSD发现缺失运行文件时  
会通过此选项拼接文件路径下载文件  
默认路径为`https://raw.githubusercontent.com/Flyleague-Collection/SimpleFSD/refs/heads/main`  
如果您运行FSD的网络环境无法连接GITHUB或者访问速度过慢  
您可以设置此选项为
`https://gh-proxy.com/https://raw.githubusercontent.com/Flyleague-Collection/SimpleFSD/refs/heads/main`  
友链：[ghproxy](https://gh-proxy.com/)

## METAR_CACHE_CLEAN_INTERVAL

注意本环境变量会覆盖[命令行参数#metar_cache_clean_interval](/configuration/command_line.md#metar_cache_clean_interval)

METAR过期报文清理间隔  
输入值应当是一个Duration字符串  
比如: 30m(30分钟), 10s(10秒), 1h(1小时)  
默认值为30m

## METAR_QUERY_THREAD

注意本环境变量会覆盖[命令行参数#metar_query_thread](/configuration/command_line.md#metar_query_thread)

METAR报文查询线程数  
当客户端一次要求多个METAR报文时  
服务端并发查询数量  
默认值为32

## FSD_RECORD_FILTER

注意本环境变量会覆盖[命令行参数#fsd_record_filter](/configuration/command_line.md#fsd_record_filter)

FSD连线记录过滤  
本选项单位为秒  
仅当连线时长高于本选项设置的数值时  
此次联飞时长才会被记录在案  
默认值为10  
