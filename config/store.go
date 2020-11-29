package config

import (
	"cloud-store.com/common"
)

const (
	//TempLocalRootDir 本地临时存储的路径
	//TempLocalRootDir = "/data/cloudstore/"
	TempLocalRootDir = "./data/clounstore"
	//TempPartRootDir 分块文件本地临时存储路径
	TempPartRootDir = "./data/partstore/"
	//CephRootDir Ceph的存储路径Prefix
	CephRootDir = "/ceph"
	//OSSRootDir OSS的存储路径prefix
	OSSRootDir = "oss/"
	//CurrentStoreType 当前的存储类型
	CurrentStoreType = common.StoreOSS
)
