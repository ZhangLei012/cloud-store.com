package main

import (
	"log"
	"time"

	"cloud-store.com/common"
	"cloud-store.com/service/dbproxy/config"
	"cloud-store.com/service/dbproxy/conn"
	dbproto "cloud-store.com/service/dbproxy/proto"
	"cloud-store.com/service/dbproxy/rpc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func startService() {
	consulReg := consul.NewRegistry( //新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
		registry.Addrs("127.0.0.1:8500"),
	)
	service := micro.NewService(
		micro.Name("go.micro.service.dbproxy"), //服务名称
		micro.RegisterTTL(time.Second*10),      //设置超时时间
		micro.RegisterInterval(time.Second*5),  //每5秒注册一次，超过10秒没有注册，就会被移除
		micro.Registry(consulReg),
		micro.Flags(common.CustomFlags...),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			//检查是否指定dbhost
			dbhost := c.String("dbhost")
			if len(dbhost) > 0 {
				log.Printf("custom db address:%v", dbhost)
				config.UpdateDBHost(dbhost)
			}
		}),
	)

	conn.InitDBConn()

	dbproto.RegisterDBProxyServiceHandler(service.Server(), new(rpc.DBProxy))
	if err := service.Run(); err != nil {
		log.Printf("Error: runing dbproxy service, err:%v", err)
		panic(err)
	}
}

func main() {
	startService()
}
