#!/usr/bin/env bash

clear
# 声明环境变量
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

TRG_PKG='main'
BUILD_TIME=$(date +"%Y%m%d.%H%M%S")
CommitHash=N/A
GoVersion=N/A
GitTag=N/A

if [[ $(go version) =~ [0-9]+\.[0-9]+\.[0-9]+ ]];
then
    GoVersion=${BASH_REMATCH[0]}
fi

GV=$(git tag || echo 'N/A')
if [[ $GV =~ [^[:space:]]+ ]];
then
    GitTag=${BASH_REMATCH[0]}
fi

GH=$(git log -1 --pretty=format:%h || echo 'N/A')
if [[ $GH =~ 'fatal' ]];
then
    CommitHash=N/A
else
    CommitHash=$GH
fi

FLAG="-X $TRG_PKG.BuildTime=$BUILD_TIME"
FLAG="$FLAG -X $TRG_PKG.CommitHash=$CommitHash"
FLAG="$FLAG -X $TRG_PKG.GoVersion=$GoVersion"
FLAG="$FLAG -X $TRG_PKG.GitTag=$GitTag"


echo 'go build'
go build -v -ldflags "$FLAG" -o mdns-proxy ./*.go
echo "$FLAG"

# export CGO_ENABLED=0 && export GOOS=linux && export GOARCH=amd64 && go build main.go
