package orm

import (
	"log"

	"cloud-store.com/service/dbproxy/conn"
)

//UserSignUp 通过用户名及密码完成用户user表的注册
func UserSignUp(userName string, password string) (res ExecResult) {
	log.Printf("Info: request to sign up, username:%v", userName)
	res.Success = false
	statement, err := conn.DBConn().Prepare("insert ignore tbl_user(`user_name`, `user_pwd`) values(?,?)")
	if err != nil {
		log.Printf("Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()

	result, err := statement.Exec(userName, password)
	if err != nil {
		log.Printf("Error: executing statement:%v, err:%v", statement, err)
		res.Msg = err.Error()
		return
	}

	if rowAffected, err := result.RowsAffected(); nil == err && rowAffected > 0 {
		res.Success = true
		return
	}
	if err != nil {
		log.Printf("Error: inserting err:%v", err)
		res.Msg = err.Error()
	} else {
		log.Printf("Error: inserting into tbl_user, username:%v, no row inserted", userName)
		res.Msg = "No row inserted"
	}
	return
}

//UserSignIn 判断密码是否一致
func UserSignIn(userName, password string) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare("select user_pwd from tbl_user where `user_name`=? limit 1")
	if err != nil {
		log.Printf("Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()

	var encodedpassword string
	err = statement.QueryRow(userName).Scan(&encodedpassword)
	if err != nil {
		log.Printf("Info: no user named:%v has been found, err:%v", userName, err)
		res.Data = "No user found"
		return
	}
	if encodedpassword != password {
		res.Data = "UserName/Password is not right"
		return
	}
	res.Success = true
	res.Data = true
	return
}

//UpdateToken 刷新用户Token
func UpdateToken(userName, token string) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare("replace into tbl_user_token (`user_name`, `user_token`) values(?,?)")
	if err != nil {
		log.Printf("Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()

	_, err = statement.Exec(userName, token)
	if err != nil {
		log.Printf("Error: executing statement:%v, err:%v", statement, err)
		res.Msg = err.Error()
		return
	}
	res.Success = true
	return
}

//GetUserInfo 查询用户信息
func GetUserInfo(userName string) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare("select user_name, signup_at from tbl_user where user_name = ? limit 1")
	if err != nil {
		log.Printf("Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()

	user := TableUser{}

	err = statement.QueryRow(userName).Scan(&user.UserName, &user.SignUpAt)
	if err != nil {
		log.Printf("Error: finding user:%v, err:%v", userName, err)
		return
	}

	res.Success = true
	res.Data = user
	return
}

//UserExist 查询用户是否存在
func UserExist(userName string) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare(
		`select 1 from tbl_user where user_name = ? limit 1`)
	if err != nil {
		log.Printf("Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()

	rows, err := statement.Query(userName)
	if err != nil {
		log.Printf("Error: Counting user:%v, err:%v", userName, err)
		res.Msg = err.Error()
		return
	}
	res.Success = true
	res.Data = map[string]bool{
		"exists": rows.Next(),
	}
	return
}
