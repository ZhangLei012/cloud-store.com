package config

const (
	//IP地址总结：每台计算机可以用三个ip地址表示，它的回环地址127.0.0.1，它在局域网内的局域网地址如192.168.1.10，它的公网ip如101.200.168.140
	//server侦听127.0.0.1:3307，本机的client可以连上127.0.0.1:3307，但连不上192.168.1.10:3307，同局域网内的其他机器、外网连不上server
	//server侦听192.168.1.10:3303，本机client可以连上192.168.1.10:3303，但连不上127.0.0.1:3307，同局域网的其他机器可以连上192.168.1.10:3303，外网无法访问server
	//server侦听0.0.0.0:3307，本机client可以连上127.0.0.1:3307，192.168.1.10:3303，同局域网机器可以连上192.168.1.10:3303，外网可以连上101.200.168.140:3307
	//UploadServiceHost 上传服务侦听地址
	UploadServiceHost = "0.0.0.0:8080"
	//UploadLBHost 上传服务LB（load balance）地址
	//UploadLBHost = "http://upload.cloud.store.com"
	//UploadLBHost = "101.200.168.140:28080"
	UploadLBHost = "http://101.200.168.140:28080"
	//DownloadLBHost 下载服务LB（load balance）地址
	//DownloadLBHost = "http://download.cloud.store.com"
	//DownloadLBHost = "101.200.168.140:38080"
	DownloadLBHost = "http://101.200.168.140:38080"
)
