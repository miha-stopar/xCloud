#!/bin/bash
#docker build -t xworker ../../docker
#docker run -d xworker
#docker run -d xworker
#docker run -d xworker
go build client0.go
go build client1.go
go build client2.go
./client0 -ip=192.168.1.11 -workerId=0
#./client1 -ip=192.168.1.11 -workerId=1
#./client2 -ip=192.168.1.11 -workerId=2

