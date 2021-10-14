#!/usr/bin/env bash
gcc -c nacdef.c
ar -rc libnacdef.a nacdef.o
echo "编译依赖静态库成功"
go build -buildmode=c-shared -o libnacgo.so main.go
echo "编译so成功"
rm -rf *.a *.o
echo "清理结束"