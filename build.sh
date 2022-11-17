#!/bin/bash
tag="0.2.8"

xgo -out EthTweet_${tag}  -targets="windows-6.0/amd64,linux/amd64" -ldflags="-w -s" .


upx EthTweet_${tag}-linux-amd64
mv EthTweet_${tag}-linux-amd64 EthTweet
zip -m EthTweet_${tag}-linux-amd64.zip EthTweet

upx EthTweet_${tag}-windows-6.0-amd64.exe
mv EthTweet_${tag}-windows-6.0-amd64.exe EthTweet.exe
zip -m EthTweet_${tag}-windows-6.0-amd64.zip EthTweet.exe