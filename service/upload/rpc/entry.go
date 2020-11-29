package rpc

import (
	"context"

	"cloud-store.com/service/upload/config"
	upproto "cloud-store.com/service/upload/proto"
)

//Upload upload结构体
type Upload struct{}

//UploadEntry 获取上传接口
func (u *Upload) UploadEntry(ctx context.Context, req *upproto.ReqEntry, resp *upproto.RespEntry) error {
	resp.Entry = config.UploadEntry
	return nil
}
