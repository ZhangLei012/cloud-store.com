package mq

import "log"

var done chan bool

//StartConsume 	消费消息
func StartConsume(qName, cName string, callbock func([]byte) bool) {
	msgs, err := channel.Consume(qName, cName, true, false, false, false, nil)
	if err != nil {
		log.Printf("Error: consuming message, err:%v", err)
		return
	}

	done = make(chan bool)

	go func() {
		for msg := range msgs {
			success := callbock(msg.Body)
			if !success {
				log.Printf("Error: processing msg:%s, err:%v", msg.Body, err)
				//TODO 将任务写入错误队列，待后续处理
			}
		}
	}()

	<-done

	channel.Close()
}

//StopConsume 停止监听队列
func StopConsume() {
	done <- true
}
