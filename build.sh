#!/bin/bash
tag="0.2.9"

xgo -out EthTweet_ -go go-1.19.x -targets="windows-6.0/amd64,linux/amd64" -ldflags="-w -s" .


upx EthTweet-linux-amd64
mv EthTweet-linux-amd64 EthTweet
zip -m EthTweet-${tag}-linux-amd64.zip EthTweet

upx EthTweet-windows-6.0-amd64.exe
mv EthTweet-windows-6.0-amd64.exe EthTweet.exe
zip -m EthTweet-${tag}-windows-6.0-amd64.zip EthTweet.exe