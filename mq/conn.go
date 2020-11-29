package mq

import (
	"cloud-store.com/config"
	"log"

	"github.com/streadway/amqp"
)

var channel *amqp.Channel
var conn *amqp.Connection

//如果异常关闭，会接收到通知
var notifyClose chan *amqp.Error

//UpdateRabbitHost 更新mq host
func UpdateRabbitHost(host string) {
	config.RabbitURL = host
}

func initChannel(rabbitHost string) bool {
	if channel != nil {
		return true
	}

	conn, err := amqp.Dial(rabbitHost)
	if err != nil {
		log.Printf("Error: dialing rabbit host:%v, err:%v", rabbitHost, err)
		return false
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Printf("Error: creating a AMQP channel, err:%v", err)
		return false
	}
	return true
}

// Init 初始化MQ连接信息
func Init() {
	// 是否开启异步转移功能，开启时才初始化rabbitmq连接
	if !config.AsyncTransferEnabled {
		return
	}
	if !initChannel(config.RabbitURL) {
		channel.NotifyClose(notifyClose)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf("Info onNotifyChannelClosed:%+v", msg)
				initChannel(config.RabbitURL)
			}
		}
	}()
}
