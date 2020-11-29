package oss

import (
	"log"

	"cloud-store.com/config"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossCli *oss.Client

func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}
	ossCli, err := oss.New(config.OSSEndpoint, config.OSSAccessKeyID, config.OSSAccessKeySecret)
	if err != nil {
		log.Printf("Error: creating oss client, err:%v", err)
		return nil
	}
	return ossCli
}

func Bucket() *oss.Bucket {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(config.OSSBucket)
		if err != nil {
			log.Printf("Error: getting bucket instance named:%v, err:%v", config.OSSBucket, err)
			return nil
		}
		return bucket
	}
	return nil
}

//DownloadURL 获取临时授权下载URL
func DownloadURL(objectKey string) string {
	signedURL, err := Bucket().SignURL(objectKey, oss.HTTPGet, 3600)
	if err != nil {
		log.Printf("Error: getting signedURL, err:%v", err)
		return ""
	}
	return signedURL
}

//GenerateFileMeta 构建文件元信息
func GenerateFileMeta(metas map[string]string) []oss.Option {
	options := []oss.Option{}
	for k, v := range metas {
		options = append(options, oss.Meta(k, v))
	}
	return options
}
