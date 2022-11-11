# eth tweet

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