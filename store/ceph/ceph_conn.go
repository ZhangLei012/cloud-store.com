package ceph

import (
	"cloud-store.com/config"
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

var cephConn *s3.S3

func GetCephConnection() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}

	//1.初始化ceph的一些信息
	auth := aws.Auth{
		AccessKey: config.CephAccessKey,
		SecretKey: config.CephSecretKey,
	}

	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          config.CephGWEndpoint,
		S3Endpoint:           config.CephGWEndpoint,
		S3BucketEndpoint:     "",
		S3LowercaseBucket:    false,
		S3LocationConstraint: false,
		Sign:                 aws.SignV2,
	}

	return s3.New(auth, curRegion)
}

//GetCephBucket 获取指定的桶
func GetCephBucket(bucket string) *s3.Bucket {
	return GetCephConnection().Bucket(bucket)
}

//PutObject 将文件储存到ceph集群
func PutObject(bucket string, path string, data []byte) error {
	return GetCephBucket(bucket).Put(path, data, "octet-stream", s3.PublicRead)
}
