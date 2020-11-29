package main

import (
	"log"
	"time"

	"cloud-store.com/common"
	dbproxy "cloud-store.com/service/dbproxy/client"
	"cloud-store.com/service/download/config"
	downloadproto "cloud-store.com/service/download/proto"
	"cloud-store.com/service/download/route"
	"cloud-store.com/service/download/rpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func startAPIService() {
	router := route.Router()
	if err := router.Run(config.DownloadServiceHost); err != nil {
		log.Fatalf("Error: running router, err:%v", err)
	}
	return
}

func startRPCService() {
	consulReg := consul.NewRegistry( //新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
		registry.Addrs("127.0.0.1:8500"),
	)
	service := micro.NewService(
		micro.Name("go.micro.service.download"),
		micro.RegisterTTL(10*time.Second),
		micro.RegisterInterval(5*time.Second),
		micro.Flags(common.CustomFlags...),
		micro.Registry(consulReg),
	)

	service.Init()

	dbproxy.Init(service)

	downloadproto.RegisterDownloadServiceHandler(service.Server(), new(rpc.Download))
	if err := service.Run(); err != nil {
		log.Fatalf("Error: starting download rpc service, err:%v", err)
	}
}

func main() {
	go startAPIService()

	startRPCService()
}
