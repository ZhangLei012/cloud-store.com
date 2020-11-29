package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"cloud-store.com/common"
	accountproto "cloud-store.com/service/account/proto"
	"github.com/gin-gonic/gin"
)

//FileQueryHandler 批量查询文件元信息
func FileQueryHandler(c *gin.Context) {
	limitCnt, _ := strconv.Atoi(c.Request.FormValue("limit"))
	userName := c.Request.FormValue("username")

	fileQueryResp, err := userClient.UserFiles(context.Background(), &accountproto.ReqUserFile{
		UserName: userName,
		Limit:    int32(limitCnt),
	})
	if err != nil {
		log.Printf("Error: query file infos, err:%v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if fileQueryResp.Code != common.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"code": fileQueryResp.Code,
			"msg":  "getting file info failed",
		})
	}

	if len(fileQueryResp.FileData) <= 0 {
		fileQueryResp.FileData = []byte("[]")
	}
	c.Data(http.StatusOK, "application/json", fileQueryResp.FileData)
}

//FileMetaUpdateHandler 更新元信息(文件重命名)
func FileMetaUpdateHandler(c *gin.Context) {
	opType := c.Request.FormValue("op")
	fileHash := c.Request.FormValue("filehash")
	userName := c.Request.FormValue("username")
	newFileName := c.Request.FormValue("newfilename")

	if opType != "0" || len(newFileName) <= 0 {
		c.Status(http.StatusForbidden)
		return
	}
	renameResp, err := userClient.UserFileRename(context.Background(), &accountproto.ReqUserFileRename{
		FileHash:    fileHash,
		UserName:    userName,
		NewFileName: newFileName,
	})
	if err != nil {
		log.Printf("Error: renaming, err:%v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if len(renameResp.FileData) <= 0 {
		renameResp.FileData = []byte("[]")
	}
	c.JSON(http.StatusOK, gin.H{
		"code": int(common.StatusOK),
		"msg":  "OK",
		"data": renameResp.FileData,
	})
}
