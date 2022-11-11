#!/bin/bash
tag="0.1.9"
xgo -out EthTweet_${tag}  -targets="darwin/amd64,windows-6.0/amd64,linux/amd64,linux/arm64" -ldflags="-w -s" .

upx EthTweet_${tag}-darwin-10.12-amd64
mv EthTweet_${tag}-darwin-10.12-amd64 EthTweet
zip -m EthTweet_${tag}-darwin-10.12-amd64.zip EthTweet

upx EthTweet_${tag}-linux-amd64
mv EthTweet_${tag}-linux-amd64 EthTweet
zip -m EthTweet_${tag}-linux-amd64.zip EthTweet

upx EthTweet_${tag}-linux-arm64
mv EthTweet_${tag}-linux-arm64 EthTweet
zip -m EthTweet_${tag}-linux-arm64.zip EthTweet

upx EthTweet_${tag}-windows-6.0-amd64.exe
mv EthTweet_${tag}-windows-6.0-amd64.exe EthTweet.exe
zip -m EthTweet_${tag}-windows-6.0-amd64.zip EthTweet.exe