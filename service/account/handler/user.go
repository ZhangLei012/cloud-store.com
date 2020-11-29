package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud-store.com/common"
	"cloud-store.com/config"
	proto "cloud-store.com/service/account/proto"
	dbClient "cloud-store.com/service/dbproxy/client"
	"cloud-store.com/util"
)

type User struct{}

//GenerateToken 生成访问token
func GenerateToken(userName string) string {
	//40位字符：md5(userName+timestamp+_tokensalt)+timestamp[:8]
	timestamp := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(userName + timestamp + "_tokensalt"))
	return tokenPrefix + timestamp[:8]
}

//SignUp 处理用户注册请求
func (u *User) SignUp(ctx context.Context, req *proto.ReqSignUp, resp *proto.RespSignUp) error {
	userName := req.UserName
	password := req.Password
	log.Printf("Info, signing up request, username:%v, password:%v", userName, password)
	if len(userName) < 3 || len(password) < 6 {
		resp.Code = common.StatusParamInvalid
		resp.Message = "Register parameters invalid"
		return nil
	}

	//对密码进行加盐及取Sha1加密
	encPassword := util.MD5([]byte(password + config.PasswordSalt))
	dbResp, err := dbClient.UserSignUp(userName, encPassword)
	if err == nil && dbResp.Success {
		resp.Code = common.StatusOK
		resp.Message = "Register succeed"
	} else {
		resp.Code = common.StatusServerError
		resp.Message = "Internal error"
	}
	return nil
}

//SignIn 处理用户登录请求
func (u *User) SignIn(ctx context.Context, req *proto.ReqSignIn, resp *proto.RespSignIn) error {
	userName := req.UserName
	password := req.Password
	encPassword := util.MD5([]byte(password + config.PasswordSalt))

	//1.检验用户名及密码
	dbResp, err := dbClient.UserSignIn(userName, encPassword)
	if err != nil || !dbResp.Success {
		resp.Code = common.StatusRegisterFailed
		return nil
	}

	//2.生成访问凭证（token）
	token := GenerateToken(userName)
	updateRes, err := dbClient.UpdateToken(userName, token)
	if err != nil || !updateRes.Success {
		resp.Code = common.StatusServerError
		return nil
	}

	//3.用户登录成功，返回凭证
	resp.Code = common.StatusOK
	resp.Username = userName
	resp.Token = token
	return nil
}

//UserInfo 查询用户信息
func (u *User) UserInfo(ctx context.Context, req *proto.ReqUserInfo, resp *proto.RespUserInfo) error {
	//1.查询用户信息
	userName := req.UserName
	dbResp, err := dbClient.GetUserInfo(userName)
	if err != nil {
		resp.Code = common.StatusServerError
		resp.Message = "Service internal error"
		return nil
	}
	//2.用户不存在
	if !dbResp.Success {
		resp.Code = common.StatusUserNotExists
		resp.Message = "User does not exist"
		return nil
	}

	//3.组装并响应用户数据
	user := dbClient.ToTableUser(dbResp.Data)
	resp.Code = common.StatusOK
	resp.UserName = user.UserName
	resp.LastActiveAt = user.LastActiveAt
	resp.SignUpAt = user.SignUpAt
	resp.Email = user.Email
	resp.Phone = user.Phone
	resp.Status = int32(user.Status)
	return nil
}
