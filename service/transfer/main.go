package main

import (
	"log"
	"time"

	"cloud-store.com/common"
	"cloud-store.com/config"
	"cloud-store.com/mq"
	dbClient "cloud-store.com/service/dbproxy/client"
	"cloud-store.com/service/transfer/process"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func StartRPCService() {
	consulReg := consul.NewRegistry( //新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
		registry.Addrs("127.0.0.1:8500"),
	)
	service := micro.NewService(
		micro.Name("go.micro.service.transfer"),
		micro.RegisterTTL(10*time.Second),
		micro.RegisterInterval(5*time.Second),
		micro.Registry(consulReg),
		micro.Flags(common.CustomFlags...),
	)

	service.Init(micro.Action(func(c *cli.Context) {
		mqHost := c.String("mqhost")
		if mqHost != "" {
			log.Printf("Info: current mqhost:%v", mqHost)
			mq.UpdateRabbitHost(mqHost)
		}
	}))

	//初始化dbproxy client
	dbClient.Init(service)

	log.Println("Transfer rpc service running...")
	if err := service.Run(); err != nil {
		panic(err)
	}
}

func StartTransferService() {
	if !config.AsyncTransferEnabled {
		log.Println("Error: async transfer not enabled")
	}

	log.Println("Start file async transfer service, working...")

	//初始化mq
	mq.Init()

	mq.StartConsume(config.TransOSSQueueName, "transfer_oss_customer", process.Transfer)

}

func main() {
	//文件转移服务
	go StartTransferService()

	//rpc服务
	StartRPCService()
}
