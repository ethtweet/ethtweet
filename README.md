# eth tweet
以太坊推文服务

技术原理，使用以太坊 personal_sign 对推文信息进行签名，然后提交到 ipfs 上保存，前端也可以读取 ipfs 上的信息，验证签名确保信息准确。

因为使用的 personal_sign 签名，而且是明文签名，不会被用于恶意攻击。

本程序提供http api 供前端程序调用，实现数据和显示分离。

本地数据持久化，支持sqlite 和 MySQL ，但是用户可以在不使用客户端的情况，就一个钱包就可以使用。

在线体验  https://app.ethtweet.io/#/home

http api文档  [api](api.md)

## 原理

使用钱包对推文信息进行签名，然后提交到ipfs，然后广播给其它节点，其它节点保存并索引，再提交到ipfs。

每个地址的推文有一个自增nonce值，0开始，用于标记推文发布顺序，节点拉数据的时候也有顺序性。

可以使用web前端，然后调用rpc接口。

## 配置文件

默认读取 `tweet.yaml` 文件，如果不存在就读取命令行参数。

## mysql配置

修改```tweet.yaml```配置文件，正确配置 MySQL 连接信息。

window环境下，如果在程序运行目录下的有一个MySQL 8.0 ，会自动启动 MySQL ，需要的文件及其路径
```
mysql\bin\mysqld.exe
mysql\bin\libprotobuf-lite.dll
```

启动节点
```
./EthTweet -config "./tweet.yaml"
```

## 交叉编译
https://github.com/techknowlogick/xgo
```
docker pull techknowlogick/xgo:latest
#export GOPATH="当前目录"
xgo -out EthTweet -targets="darwin/amd64,windows-6.0/amd64,linux/amd64,linux/arm64" -ldflags="-w -s" .
```

window下编译安卓
下载https://dl.google.com/android/repository/android-ndk-r25b-windows.zip

使用cmd执行
```
SET CGO_ENABLED=1
SET GOOS=android
SET GOARCH=arm64
set CC=D:\sdk\ndk\25.1.8937393\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android21-clang.cmd
set CXX=D:\sdk\ndk\25.1.8937393\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android21-clang++.cmd
go build --tags "android"  -ldflags="-s -w"  -o ethtweet
```
## docker-compose 运行

docker-compose目录下，集成一个带MySQL的环境，可一键启动。

## todo 

增加节点统计，记录最长在线时间，每次启动的时候连接

加密私信功能，可以进行发生加密私信，地址持有人才可以解密阅读。

app端

pc版web

增加附件支持：可以上传类似zip的文件，然后下载

增加附件播种：节点进行再次播种
