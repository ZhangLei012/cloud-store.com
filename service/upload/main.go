package main

import (
	"os"
	"time"

	"cloud-store.com/config"
	"cloud-store.com/mq"
	dbproxy "cloud-store.com/service/dbproxy/client"
	"cloud-store.com/service/upload/api"
	uploadconfig "cloud-store.com/service/upload/config"
	upproto "cloud-store.com/service/upload/proto"
	"cloud-store.com/service/upload/route"
	"cloud-store.com/service/upload/rpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func startRPCService() {
	consulReg := consul.NewRegistry( //新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
		registry.Addrs("127.0.0.1:8500"),
	)
	service := micro.NewService(
		micro.Name("go.micro.service.upload"),
		micro.RegisterTTL(10*time.Second),
		micro.RegisterInterval(5*time.Second),
		micro.Registry(consulReg),
		//micro.Flags(common.CustomFlags...),
	)

	service.Init(
	/*
		micro.Action(func(c *cli.Context) {
			mqHost := c.String("mqhost")
			if mqHost != "" {
				mq.UpdateRabbitHost(mqHost)
			}
		}),
	*/
	)

	//初始化dbproxy client
	dbproxy.Init(service)

	//初始化mq
	mq.Init()
	api.Init(service)
	//注册handler
	upproto.RegisterUploadServiceHandler(service.Server(), new(rpc.Upload))
	if err := service.Run(); err != nil {
		panic(err)
	}
}

func startAPIService() {
	router := route.Router()
	router.Run(uploadconfig.UploadServiceHost)
}

func main() {
	//0777的含义：在golang中，数字以0开头表示是一个八进制数，以0x/0X开头表示是一个十六进制数
	//7表示111，分别对应读、写、执行位，每个7依次代表文件所有者、同group人员、其他人员的权限，所以0777表示所有人都可读可写可执行
	os.MkdirAll(config.TempLocalRootDir, 0777)
	os.MkdirAll(config.TempPartRootDir, 0777)

	//
	go startAPIService()

	startRPCService()
}
