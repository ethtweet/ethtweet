#!/bin/bash
tag=$1

xgo -out EthTweet -go go-1.20.x -targets="windows-6.0/amd64,linux/amd64" -ldflags="-w -s" .


mv EthTweet-linux-amd64 EthTweet
chmod 0777 EthTweet
zip -r EthTweet-${tag}-linux-amd64.zip EthTweet templates
sha512sum EthTweet-${tag}-linux-amd64.zip -t > EthTweet-${tag}-linux-amd64.zip.sha512

mv EthTweet-windows-6.0-amd64.exe EthTweet.exe
zip -r EthTweet-${tag}-windows-amd64.zip EthTweet.exe templates
sha512sum EthTweet-${tag}-windows-amd64.zip -t > EthTweet-${tag}-windows-amd64.zip.sha512

docker build -t chenjia404/ethtweet --build-arg VERSION="${tag}"  .
docker image push  chenjia404/ethtweet

