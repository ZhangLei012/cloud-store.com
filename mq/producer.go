package mq

import (
	"cloud-store.com/config"
	"log"

	"github.com/streadway/amqp"
)

//Publish 发布消息
func Publish(exchange, routingKey string, msg []byte) bool {
	if !initChannel(config.RabbitURL) {
		log.Printf("Error: initing rabbitmq host:%v", config.RabbitURL)
		return false
	}

	if nil == channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	}) {
		return true
	}
	return false
}
