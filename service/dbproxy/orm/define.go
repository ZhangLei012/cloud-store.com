package orm

import "database/sql"

//TableFile 文件表结构体
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

//TableUser 用户表结构
type TableUser struct {
	UserName     string
	Email        string
	Phone        string
	SignUpAt     string
	LastActiveAt string
	Status       int
}

//TableUserFile 用户文件表结构
type TableUserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

//ExecResult sql函数执行的结果
type ExecResult struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}
