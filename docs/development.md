## 二次开发指引

如果需要更多指引或者对代码有任何疑问, 请联系我
邮箱: Half_nothing@163.com或者halfnothingno@gmail.com

### 如果想修改输出的whazzup文件格式

1. 首先前往文件[internal/interfaces/fsd/whazzup.go](../internal/interfaces/fsd/whazzup.go), 这里存放whazzup文件的格式,
   如果您想输出文字类型的whazzup可以跳过这一步
2. 前往文件[internal/fsd_server/packet/client_manager.go](../internal/fsd_server/packet/client_manager.go), 找到函数
   `generateWhazzupFile`, 此函数负责生成Whazzup文件, 详细函数介绍请看代码注释
