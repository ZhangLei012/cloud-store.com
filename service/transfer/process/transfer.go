package process

import (
	"encoding/json"
	"log"
	"os"

	"cloud-store.com/mq"
	dbClient "cloud-store.com/service/dbproxy/client"
	"cloud-store.com/store/oss"
	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func Transfer(msg []byte) bool {
	var transferData mq.TransferData
	err := json.Unmarshal(msg, &transferData)
	if err != nil {
		log.Printf("Error: unmarshaling:%s, err:%v", msg, err)
		return false
	}
	log.Printf("Transfer: handling transfer data:%v", transferData)
	file, err := os.Open(transferData.CurLocation)
	if err != nil {
		log.Printf("Error: opening file:%v, err:%v", transferData.CurLocation, err)
		return false
	}
	defer file.Close()

	options := []alioss.Option{
		alioss.ContentDisposition("attachment;filename=\"" + transferData.FileName + "\""),
	}
	err = oss.Bucket().PutObject(transferData.DestLocation, file, options...)
	if err != nil {
		log.Printf("Error: putting object:%v, err:%v", file.Name(), err)
		return false
	}
	updateResult, err := dbClient.UpdateFileLocation(transferData.FileHash, transferData.DestLocation)
	if err != nil || !updateResult.Success {
		log.Printf("Error: updating file location, err:%v", err)
		return false
	}
	return true
}
