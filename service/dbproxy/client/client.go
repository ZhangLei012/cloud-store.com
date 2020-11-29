package client

import (
	"context"
	"encoding/json"
	"log"

	"cloud-store.com/service/dbproxy/orm"
	dbproto "cloud-store.com/service/dbproxy/proto"
	"github.com/micro/go-micro"
	"github.com/mitchellh/mapstructure"
)

//FileMeta 文件元信息
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var (
	dbCli dbproto.DBProxyService
)

//Init 初始化服务
func Init(service micro.Service) {
	dbCli = dbproto.NewDBProxyService("go.micro.service.dbproxy", service.Client())
}

//TableFileToFileMeta TableFile转化为FileMeta
func TableFileToFileMeta(tableFile orm.TableFile) FileMeta {
	return FileMeta{
		FileSha1: tableFile.FileHash,
		FileName: tableFile.FileName.String,
		FileSize: tableFile.FileSize.Int64,
		Location: tableFile.FileAddr.String,
	}
}

//execAction 向dbproxy请求执行action
func execAction(funcName string, paramsJSON []byte) (*dbproto.RespExec, error) {
	return dbCli.ExecuteAction(context.TODO(), &dbproto.ReqExec{
		Action: []*dbproto.SingleAction{
			&dbproto.SingleAction{
				Name:   funcName,
				Params: paramsJSON,
			},
		},
	})
}

//parseBody 转化rpc返回的结果
func parseBody(resp *dbproto.RespExec) *orm.ExecResult {
	if resp == nil || resp.Data == nil {
		return nil
	}
	resList := []orm.ExecResult{}
	_ = json.Unmarshal(resp.Data, &resList)
	if len(resList) > 0 {
		return &resList[0]
	}
	return nil
}

func ToTableUser(src interface{}) orm.TableUser {
	user := orm.TableUser{}
	mapstructure.Decode(src, &user)
	return user
}

func ToTableFile(src interface{}) orm.TableFile {
	file := orm.TableFile{}
	mapstructure.Decode(src, &file)
	return file
}

func ToTableFiles(src interface{}) []orm.TableFile {
	files := []orm.TableFile{}
	mapstructure.Decode(src, &files)
	return files
}

func ToTableUserFile(src interface{}) orm.TableUserFile {
	userFile := orm.TableUserFile{}
	mapstructure.Decode(src, &userFile)
	return userFile
}

func ToTableUserFiles(src interface{}) []orm.TableUserFile {
	userFiles := []orm.TableUserFile{}
	mapstructure.Decode(src, &userFiles)
	return userFiles
}

func GetFileMeta(fileHash string) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{fileHash})
	res, err := execAction("GetFileMeta", params)
	return parseBody(res), err
}

func GetFileMetaList(limit string) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{limit})
	res, err := execAction("GetFileMetaList", params)
	return parseBody(res), err
}

//OnFileUploadFinished 文件上传成功后更新文件元信息表
func OnFileUploadFinished(fileMeta FileMeta) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize, fileMeta.Location})
	res, err := execAction("OnFileUploadFinished", params)
	return parseBody(res), err
}

func UpdateFileLocation(fileHash, newLocation string) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{fileHash, newLocation})
	res, err := execAction("UpdateFileLocation", params)
	return parseBody(res), err
}

func UserSignUp(userName, encPassword string) (*orm.ExecResult, error) {
	log.Printf("Info: request to sign up, username:%v", userName)
	params, _ := json.Marshal([]interface{}{userName, encPassword})
	res, err := execAction("UserSignUp", params)
	return parseBody(res), err
}

func UserSignIn(userName, encPassword string) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{userName, encPassword})
	res, err := execAction("UserSignIn", params)
	return parseBody(res), err
}

func GetUserInfo(userName string) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{userName})
	res, err := execAction("GetUserInfo", params)
	return parseBody(res), err
}

func UserExist(userName string) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{userName})
	res, err := execAction("UserExist", params)
	return parseBody(res), err
}

func UpdateToken(userName string, token string) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{userName, token})
	res, err := execAction("UpdateToken", params)
	return parseBody(res), err
}

func QueryUserFileMeta(userName string, fileHash string) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{userName, fileHash})
	res, err := execAction("QueryUserFileMeta", params)
	return parseBody(res), err
}

func QueryUserFileMetas(userName string, limit int) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{userName, limit})
	res, err := execAction("QueryUserFileMetas", params)
	return parseBody(res), err
}

//OnUserFileUploadFinished 新增/更新用户文件元信息表
func OnUserFileUploadFinished(userName string, fileMeta FileMeta) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{userName, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize})
	res, err := execAction("OnUserFileUploadFinished", params)
	return parseBody(res), err
}

func RenameFileName(userName, fileHash, fileName string) (*orm.ExecResult, error) {
	params, _ := json.Marshal([]interface{}{userName, fileHash, fileName})
	res, err := execAction("RenameFileName", params)
	return parseBody(res), err
}
