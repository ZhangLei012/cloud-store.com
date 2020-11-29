package config

import "fmt"

var (
	//MySQLSource 要连接的数据库源
	//其中root:123456是用户名密码
	//127.0.0.1:3306是用户名端口号
	//cloudstore是数据库名
	//charset=utf8指定了数据以utf8字符编码进行传输
	MySQLSource = "root:123456@tcp(127.0.0.1:3307)/cloudstore?charset=utf8"
)

//UpdateDBHost 更新db的ip与端口等
func UpdateDBHost(host string) {
	MySQLSource = fmt.Sprintf("root:123456@tcp(%s)/cloudstore?charset=utf8", host)
}
