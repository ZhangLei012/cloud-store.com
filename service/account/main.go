package main

import (
	"time"

	"cloud-store.com/common"
	"cloud-store.com/service/account/handler"
	proto "cloud-store.com/service/account/proto"
	dbproxy "cloud-store.com/service/dbproxy/client"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func main() {
	consulReg := consul.NewRegistry( //新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
		registry.Addrs("127.0.0.1:8500"),
	)
	service := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(consulReg),
		micro.Flags(common.CustomFlags...),
	)
	//获取命令行参数，初始化服务
	service.Init()

	//初始化dbproxy client
	dbproxy.Init(service)

	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		panic(err)
	}
}
