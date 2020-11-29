package handler

import (
	"context"
	"encoding/json"

	"cloud-store.com/common"
	proto "cloud-store.com/service/account/proto"
	dbClient "cloud-store.com/service/dbproxy/client"
)

//UserFiles 获取用户文件表
func (u *User) UserFiles(ctx context.Context, req *proto.ReqUserFile, resp *proto.RespUserFile) error {
	dbResp, err := dbClient.QueryUserFileMetas(req.UserName, int(req.Limit))
	if err != nil || !dbResp.Success {
		resp.Code = common.StatusServerError
		return nil
	}
	userFiles := dbClient.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		resp.Code = common.StatusServerError
		return nil
	}
	resp.FileData = data
	resp.Code = common.StatusOK
	return nil
}

//UserFileRename 用户文件重命名
func (u *User) UserFileRename(ctx context.Context, req *proto.ReqUserFileRename, resp *proto.RespUserFileRename) error {
	dbResp, err := dbClient.RenameFileName(req.UserName, req.FileHash, req.NewFileName)
	if err != nil || !dbResp.Success {
		resp.Code = common.StatusServerError
		return nil
	}

	userFile := dbClient.ToTableUserFile(dbResp.Data)
	data, _ := json.Marshal(userFile)
	resp.FileData = data
	resp.Code = common.StatusOK
	return nil
}
