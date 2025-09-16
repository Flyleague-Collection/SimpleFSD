## 协议简介

FSD9 协议是目前模拟飞行平台使用最广泛的协议

1. 优点是协议简单, 管制员与飞行员连线不需要额外配置
2. 缺点是密码为明文传输

## 协议规范

### 传输协议

FSD9使用TCP作为传输协议, 默认端口号为6809, 流式传输  
FSD9使用明文传输所有命令, 命令与命令之间使用`\r\n`作为分隔符

### 命令规范

如下是一条典型的FSD9命令
`#AA2352_OBS:SERVER:2352:2352:123456:1:9:1:0:29.86379:119.49287:100`  
命令大致可以分为命令头部与命令载荷两部分  
比如上面的命令就可以分为命令头部`#AA2352_OBS:SERVER`与命令载荷`2352:2352:123456:1:9:1:0:29.86379:119.49287:100`  
命令头部又可以继续细分为命令码、发送者与接受者  
在这个命令中, 命令码是`#AA`, 发送者是`2352_OBS`, 接受者是`SERVER`  
翻译过来就是, `2352_OBS`这个管制员向服务器(`SERVER`)发送了登陆请求(`#AA`)  
FSD9中几乎所有命令都可以像这样解释  
当然也有某些信息就是单纯的广播消息  
这种消息就只有发送者没有接受者

### 命令码

如下为命令码一览(可能会有遗漏), 来自[Swift]源码

| 命令码 | Swift 消息类型                           | 推测类型        |
|:----|:-------------------------------------|:------------|
| #AA | MessageType::AddAtc                  | 管制员上线       |
| #AP | MessageType::AddPilot                | 飞行员上线       |
| %   | MessageType::AtcDataUpdate           | 管制员主视程点更新   |
| $CQ | MessageType::ClientQuery             | 客户端查询       |
| $CR | MessageType::ClientResponse          | 客户端查询回报     |
| #DA | MessageType::DeleteATC               | 管制员下线       |
| #DP | MessageType::DeletePilot             | 飞行员下线       |
| $FP | MessageType::FlightPlan              | 飞行计划提交      |
| #PC | MessageType::ProController           | ?           |
| $!! | MessageType::KillRequest             | 踢出请求        |
| @   | MessageType::PilotDataUpdate         | 飞行员数据更新     |
| ^   | MessageType::VisualPilotDataUpdate   | 虚拟飞行员数据更新   |
| #SL | MessageType::VisualPilotDataPeriodic | 虚拟飞行员周期数据更新 |
| #ST | MessageType::VisualPilotDataStopped  | 虚拟飞行员停止数据   |
| $SF | MessageType::VisualPilotDataToggle   | 切换虚拟飞行员数据上报 |
| $PI | MessageType::Ping                    | Ping        |
| $PO | MessageType::Pong                    | Pong        |
| $ER | MessageType::ServerError             | 服务器错误       |
| #DL | MessageType::ServerHeartbeat         | 服务器心跳包      |
| #TM | MessageType::TextMessage             | 文本消息        |
| #SB | MessageType::PilotClientCom          | 飞行员客户端交流    |
| $XX | MessageType::Rehost                  | ?           |
| #MU | MessageType::Mute                    | ?           |

***下面为抓包获取***

| 命令码 | 推测类型      |
|:----|:----------|
| '   | 管制员副视程点更新 |

#### \#AA

管制员注册   
#AAZSHA_CTR:SERVER:2352:2352:123456:5:9:1:0:29.86379:119.49287:100

| 内容        | 含义       | 备注              |
|:----------|:---------|:----------------|
| #AA       | 命令码      |                 |
| ZSHA_CTR  | 发送方      |                 |
| SERVER    | 接收方      |                 |
| 2352      | RealName |                 |
| 2352      | Cid      |                 |
| 123456    | 密码       |                 |
| 5         | 请求的权限登记  |                 |
| 9         | FSD协议版本  | 对于FSD9来说这个值固定为9 |
| 1         | 未知       | 推测为一固定值         |
| 0         | 未知       | 推测与后面的纬度和经度为一体  |
| 29.86379  | 纬度       |                 |
| 119.49287 | 经度       |                 |
| 100       | 固定值      |                 |

#### \#AP

添加飞行员  
#APB2352:SERVER:2352:123456:1:9:16:2352 ZGHA

| 内容        | 含义       | 备注              |
|:----------|:---------|:----------------|
| #AP       | 命令码      |                 |
| B2352     | 发送方      |                 |
| SERVER    | 接收方      |                 |
| 2352      | CID      |                 |
| 123456    | 密码       |                 |
| 1         | 请求权限登记   |                 |
| 9         | FSD协议版本  | 对于FSD9来说这个值固定为9 |
| 16        | 模拟器类型    |                 |
| 2352 ZGHA | RealName |                 |

***模拟器类型一览***

| 定义         | 描述                           | 值  | 
|:-----------|:-----------------------------|:---|
| Unknown    | Unknown simulator type       | 0  |
| MSFS95     | MS Flight Simulator 95       | 1  |
| MSFS98     | MS Flight Simulator 98       | 2  |
| MSCFS      | MS Combat Flight Simulator   | 3  |
| MSFS2000   | MS Flight Simulator 2000     | 4  |
| MSCFS2     | MS Combat Flight Simulator 2 | 5  |
| MSFS2002   | MS Flight Simulator 2002     | 6  |
| MSCFS3     | MS Combat Flight Simulator 3 | 7  |
| MSFS2004   | MS Flight Simulator 2004     | 8  |
| MSFSX      | MS Flight Simulator X        | 9  |
| MSFS       | MS Flight Simulator 2020     | 10 |
| MSFS2024   | MS Flight Simulator 2024     | 11 |
| XPLANE8    | X-Plane 8                    | 12 |
| XPLANE9    | X-Plane 9                    | 13 |
| XPLANE10   | X-Plane 10                   | 14 |
| XPLANE11   | X-Plane 11                   | 15 |
| XPLANE12   | X-Plane 12                   | 16 |
| P3Dv1      | Prepar3D V1                  | 17 |
| P3Dv2      | Prepar3D V2                  | 18 |
| P3Dv3      | Prepar3D V3                  | 19 |
| P3Dv4      | Prepar3D V4                  | 20 |
| P3Dv5      | Prepar3D V5                  | 21 |
| FlightGear | Flight Gear                  | 22 |

#### %

管制员主视程点更新  
%ZSHA_CTR:24550:6:600:5:29.86379:119.49287:0

| 内容        | 含义   | 备注                     |
|:----------|:-----|:-----------------------|
| %         | 命令码  |                        |
| ZSHA_CTR  | 发送方  |                        |
| 24550     | 管制频率 | (实际管制频率为该数值+100000)kHz |
| 6         | 席位代码 |                        |
| 600       | 视程范围 |                        |
| 5         | 权限代码 |                        |
| 29.86379  | 纬度   |                        |
| 119.49287 | 精度   |                        |

#### '

管制员副视程点更新  
'ZSHA_CTR:0:36.67349:120.45621

| 内容        | 含义     | 备注                               |
|:----------|:-------|:---------------------------------|
| '         | 命令码    |                                  |
| ZSHA_CTR  | 发送方    |                                  |
| 0         | 副视程点索引 | 管制有一个主视程点和三个副视程点, 所以0对应的是.vis2命令 |
| 36.67349  | 纬度     |                                  |
| 120.45621 | 精度     |                                  |

#### $CQ

#### $CR

#### \#DA

#### \#DP

#### $FP

#### \#PC

#### $!!

#### @

#### ^

#### \#SL

#### \#ST

#### $SF

#### $PI

#### $PO

#### $ER

#### \#DL

#### \#TM

#### \#SB

#### $XX

#### \#MU

[Swift]: https://github.com/swift-project/pilotclient/blob/main/src/core/fsd/fsdclient.cpp#L1049