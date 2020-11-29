package common

//ErrCode 错误码
type ErrCode int32

const (
	_ int32 = iota + 9999
	//StatusOK 10000正常
	StatusOK
	//StatusParamInvalid 10001 参数无效
	StatusParamInvalid
	//StatusServerError 10002 服务出错
	StatusServerError
	//StatusRegisterFailed 10003 注册失败
	StatusRegisterFailed
	//StatusLoginFailed 10004 登录失败
	StatusLoginFailed
	//StatusTokenInvalid 10005 token无效
	StatusTokenInvalid
	//StatusUserNotExists 10006 用户不存在
	StatusUserNotExists
)
