# EasyNotification

[English](README.md) | [中文文档](README_CN.md)

## 摘要
 
轻量级的通知平台,用于为特定项目搭建一个通知系统.  

### 账户

用户可以在esn中添加账户,并授予`account`(账户操作) `push`(推送通知) `pull`(拉取通知) 权限.  
`root`账户具有所有权限,可以添加/删除账户,推送/拉取通知.  

### 通知

每个通知具有`ID(编号)` `Title(标题)` `Content(内容)` `Time(时间)` `Target(接收者)` `Source(来源)` 等属性  
`Target`可以指定接收通知的用户,将`Target`设置为`_global_`以向所有用户发送通知  
程序可以通过API来拉取发送给某一账户的通知

### 交互

所有操作均由程序调用API连接到`esn-daemon`以进行,类似于MySQL的操作模式  
亦可以通过我们发行的`终端`来进行用户交互


## 项目结构

| 单元 | 功能 | 仓库地址 |  
| :----- | :----- | :----- |
| esn-daemon | 处理由其他程序通过API发起的连接 | <https://github.com/EasyNotification/esn-daemon> |
| esn-api-golang | 为Golang构建的用于连接daemon的函数库 | <https://github.com/EasyNotification/esn-api-golang>
| esn-api-java | 为Java构建的用于连接daemon的类库 | <https://github.com/EasyNotification/esn-api-java> |
| esn-api-python | 为Python构建的用于连接daemon的库 | <https://github.com/EasyNotification/esn-api-python> |
| esn-terminal-swing | 使用esn-api-java构建的用于管理daemon的账户和消息的终端 | <https://github.com/EasyNotification/esn-terminal-swing> |
| esn-terminal-android | 运行于Android平台上的终端 | <https://github.com/Soulter/esn-terminal-android> |

## (可选) 构建esn-daemon 

如果你想自行编译esn-daemon:

1. 在你的设备上安装golang(版本大于11)
2. 通过命令 `export GO111MODULE=on` (Linux/MacOS) 或者 `set GO111MODULE=on` (Windows) 启用GO111MODULE
3. 克隆本仓库或者从release中下载源代码(推荐)
4. 进入到 `go.mod` 所在目录,运行命令 `go mod tidy` 来解决依赖库.
5. 运行命令 `go build -o bin/esnd-linux src/main/main.go` (Linux/MacOS) 或者 `go build -o bin\esnd-windows src\main\main.go` (Windows).
6. 你可以在 `bin` 目录中找到可执行文件.

## 安装

以下步骤告诉您如何为你自己的项目配置esn-daemon

### 设备需求

一个拥有可访问IP地址的主机，作为服务器运行esn-daemon

### 准备

1. 在服务器上安装MySQL  
2. 在MySQL中添加一个账户，并为其授予访问某一数据库的所有权限  
3. 记住账户名,账户密码及数据库名

### 配置esn-daemon(esnd)

1. 创建esnd的工作文件夹并从[这里](https://github.com/EasyNotification/esn-daemon/releases/latest)下载预构建的可执行文件
2. 在工程目录运行可执行文件  
3. esnd将会创建一个配置文件并自动退出  
4. 进入目录`config/` 并编辑文件 `esnd.conf`  
5. 设置 `service.port` 为esnd将用于接受连接的端口  
6. 设置 `db.user` `db.addr` `db.pass` `db.database` 为上一步中配置的 `MySQL用户名` `MySQL数据库地址:端口` `MySQL用户密码` `数据库名称`    
7. 设置 `root.mask` 为希望esnd中root账户的密码,root账户具有esnd系统中的一切权限
8. 如果你想开启调试模式并且在esnd运行时查看更多输出信息,请在 `esnd.conf` 中追加 `debug.enable=true`
9. 请确保esnd所用于接受连接的端口已经被系统防火墙放行,然后启动esn-daemon.你现在可以使用esn的终端来检测esnd的配置


## 用法

将API导入到您的项目中,并使用API提供的方法连接到esn-daemon以推送/拉取通知 或者 添加/删除 账户.  
请参阅API的README.
