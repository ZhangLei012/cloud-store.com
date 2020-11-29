package main

import (
	"time"

	searchproto "cloud-store.com/service/search/proto"
	"cloud-store.com/service/search/route"
	"cloud-store.com/service/search/rpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func startAPIService() {
	router := route.Router()
	router.Run(":48080")
}

func main() {
	go startAPIService()

	consulReg := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)
	service := micro.NewService(
		micro.Name("go.micro.service.search"),
		micro.RegisterTTL(10*time.Second),
		micro.RegisterInterval(5*time.Second),
		micro.Registry(consulReg),
	)

	//初始化服务
	service.Init()

	searchproto.RegisterSearchServiceHandler(service.Server(), new(rpc.SearchService))
	if err := service.Run(); err != nil {
		panic(err)
	}
}
