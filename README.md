# EasyNotification

[English](README.md) | [中文文档](README_CN.md)

## Abstract

A light-weight,easy-to-use,cross-platfrom api for programmers to build up a notification system for their own projects.  

### Account
Users can add account through API and grant `account`(account operation) `push`(push notifications) `pull`(pull notifications) permissions.  
`root` account has all permissions,and can add/remove accounts or push/pull notifications.

### Notification

Each notification has its own `ID` `Title` `Content` `Time` `Target` `Source` fileds.  
Set `Target` to select specific user to receive this notification,or set `Target` to `_global_` to send this notification to all users.  
Programs can pull notifications sent to specific account through APIs.

### Interacting

Programs connect to `esn-daemon` through APIs to interact with `esn-daemon`,like the style of MySQL's interacting.  
Users can also use `Terminal` we published to interact with esn-daemon.


## Project Structure

| unit | function | repo |  
| :----- | :----- | :----- |
| esn-daemon | Handle connections from api | <https://github.com/EasyNotification/esn-daemon> |
| esn-api-golang | API for Golang projects to connect to esn-daemon | <https://github.com/EasyNotification/esn-api-golang>
| esn-api-java | API for Java projects to connect to esn-daemon | <https://github.com/EasyNotification/esn-api-java> |
| esn-api-python | API for Python projects to connect to esn-daemon | <https://github.com/EasyNotification/esn-api-python> |
| esn-terminal-swing | Made with esn-api-java,for users to manage notifications or accounts | <https://github.com/EasyNotification/esn-terminal-swing> |
| esn-terminal-android | A terminal for Android platform | <https://github.com/Soulter/esn-terminal-android> |


## (Optional) Build esn-daemon

If you don't want to use pre-build files,you can install golang on your device and build it by yourself.

1. Install golang(version above 11) on your device.
2. Enable GO111MODULE by execute `export GO111MODULE=on` on Linux/MacOS or `set GO111MODULE=on` on Windows.
3. Clone this repo or download source code from release(recommended).
4. Change directory to where the file `go.mod` in,execute `go mod tidy` to solve requirements.
5. Execute `go build -o bin/esnd-linux src/main/main.go` on Linux/MacOS or `go build -o bin\esnd-windows src\main\main.go` on Windows.
6. You can find executable file in `bin` directory.

## Install

Here are the steps to guide you to configure esn-daemon for your own projects.

### Requirements

A host with reachable IP address as server to run esn-daemon.

### Get Ready

1. Install MySQL on server.  
2. Add an account in MySQL and grant all permissions to access a specific database.  
3. Remember the user name,password and database name.

### Configure esn-daemon(esnd)

1. Make a work directory for esnd and download pre-build executable file from [here](https://github.com/EasyNotification/esn-daemon/releases/latest).  

2. Use `./esnd-linux-x64` on Linux or `.\esnd-windows-x64.exe` on Windows to launch esnd first time.  
3. Esnd will auto generate config files and exit.
4. Change directory to `config/` and edit `esnd.conf`.
5. Set `service.port` to the port specific for esnd.  
6. Set `db.user` `db.addr` `db.pass` `db.database` to the `MySQL user name` `MySQL address and port` `MySQL user password` `database name` as you configured in previous steps.  
7. Set `root.mask` to the password for root account in esnd you expected,root account has all permissions in esnd system.
8. If you want to enable debug mode and see more output during esnd runtime,append `debug.enable=true` in the `esnd.conf`.
9. Make sure the serivce port is open by firewall and then launch esnd.You can test esnd configuration in esn-terminal now.

## Usage

Import APIs to your projects,and use provided ways to connect to daemon to push/pull notifications or add/remove accounts.  
Please refer to READMEs of API repos.