package rpc

import (
	"context"

	"cloud-store.com/service/download/config"
	proto "cloud-store.com/service/download/proto"
)

type Download struct{}

func (d *Download) DownloadEntry(ctx context.Context, in *proto.ReqEntry, out *proto.RespEntry) error {
	out.Entry = config.DownloadEntry
	return nil
}
