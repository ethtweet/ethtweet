# eth tweet
以太坊推文服务

技术原理，使用以太坊 personal_sign 对推文信息进行签名，然后提交到 ipfs 上保存，前端也可以读取 ipfs 上的信息，验证签名确保信息准确。

因为使用的 personal_sign 签名，而且是明文签名，不会被用于恶意攻击。

本程序提供http api 供前端程序调用，实现数据和显示分离。

本地数据持久化，支持sqlite 和 MySQL ，但是用户可以在不使用客户端的情况，就一个钱包就可以使用。

## mysql配置

修改```tweet.yaml```配置文件，正确配置 MySQL 连接信息。

启动节点
```
./EthTweet -config "./tweet.yaml"
```

## 交叉编译
https://github.com/techknowlogick/xgo
```
docker pull techknowlogick/xgo:latest
#export GOPATH="当前目录"
xgo   -targets=darwin/amd64,windows-6.0/amd64,linux/amd64,linux/arm64 -ldflags="-w -s" .
```

window下编译安卓
下载https://dl.google.com/android/repository/android-ndk-r22b-linux-x86_64.zip

使用cmd执行
```
SET CGO_ENABLED=1
SET GOOS=android
SET GOARCH=arm64
set CC=D:\android\ndk\22.1.7171670\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android21-clang.cmd
set CXX=D:\android\ndk\22.1.7171670\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android21-clang++.cmd
go build --tags "android"  -ldflags="-s -w"  -o ipfs
```