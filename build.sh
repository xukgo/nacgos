#!/usr/bin/env bash
#gcc -c nacdef.c
#ar -rc libnacdef.a nacdef.o
#echo "编译依赖静态库成功"
go env -w GOPROXY=https://goproxy.cn,direct
go build -buildmode=c-shared -o libnacgo.so main.go
echo "编译so成功"
cp libnacgo.* ../../C/nacgoDemo