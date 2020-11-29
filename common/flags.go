package common

import (
	"github.com/micro/cli"
)

var CustomFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "dbhost",
		Value: "127.0.0.1:3307",
		Usage: "database address",
	},
	cli.StringFlag{
		Name:  "mqhost",
		Value: "amqp://guest:guest@localhost:5672/",
		Usage: "rabbitmq address",
	},
	cli.StringFlag{
		Name:  "redishost",
		Value: "127.0.0.1:6379",
		Usage: "redis address",
	},
	cli.StringFlag{
		Name:  "cephhost",
		Value: "127.0.0.1",
		Usage: "ceph address",
	},
}
