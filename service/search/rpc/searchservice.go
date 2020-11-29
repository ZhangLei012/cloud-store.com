package rpc

import (
	"context"
	"log"

	"cloud-store.com/common"
	"cloud-store.com/service/search/esclient"
	searchproto "cloud-store.com/service/search/proto"
)

type SearchService struct {
}

func (service *SearchService) SaveDocument(ctx context.Context, req *searchproto.ReqSaveDocument, resp *searchproto.RespSaveDocument) error {
	index := req.Index

	typ := req.Typ
	userFileMeta := req.UserFileMeta
	indexService := esclient.Client().Index().Index(index).Type(typ).BodyJson(userFileMeta)
	if userFileMeta.FileSha1 != "" {
		indexService.Id(userFileMeta.FileSha1)
	}
	result, err := indexService.Do(context.Background())
	if err != nil {
		log.Printf("SaveDocument: Error: failed to save filemeta:%v, err:%v", userFileMeta, err)
		resp.Code = common.StatusServerError
		resp.Success = false
		resp.Message = err.Error()
		return nil
	}

	log.Printf("SaveDocument: Info: result:%v", result)
	resp.Code = common.StatusOK
	resp.Success = true
	resp.Message = "OK"
	return nil
}
