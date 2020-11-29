package test

import (
	"context"
	"fmt"
	"testing"

	"cloud-store.com/service/search/api"
	go_micro_service_search "cloud-store.com/service/search/proto"
	"cloud-store.com/service/search/rpc"
	"github.com/gin-gonic/gin"
)

func TestSaveDocument(t *testing.T) {
	searchService := rpc.SearchService{}
	req := &go_micro_service_search.ReqSaveDocument{
		Index: "cloudstore",
		Typ:   "user_file_test",
		UserFileMeta: &go_micro_service_search.UserFileMeta{
			FileSha1: "dafds321",
			FileName: "main.go",
		},
	}
	resp := &go_micro_service_search.RespSaveDocument{}
	err := searchService.SaveDocument(context.Background(), req, resp)
	if err != nil || !resp.Success {
		panic(fmt.Sprintf("Error:%v, msg:%v", err, resp.Message))
	}
}

func TestSearchUserFileHandler(t *testing.T) {
	var c gin.Context
	c.Request.Form.Add("fileSha1", "4550e828266e2939268325bd4f0abf346156235a")
	c.Request.Form.Add("fileName", "index.html")
	api.SearchUserFileHandler(&c)

}
