package mq

import (
	"cloud-store.com/common"
)

// TransferData 将要写到rabbitmq的数据的结构体
type TransferData struct {
	FileHash      string
	FileName      string
	CurLocation   string
	DestLocation  string
	DestStoreType common.StoreType
}
