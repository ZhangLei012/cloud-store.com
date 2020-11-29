package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud-store.com/common"
	"cloud-store.com/config"
	"cloud-store.com/mq"
	dbClient "cloud-store.com/service/dbproxy/client"
	"cloud-store.com/service/dbproxy/orm"
	go_micro_service_search "cloud-store.com/service/search/proto"
	"cloud-store.com/store/ceph"
	"cloud-store.com/store/oss"
	"cloud-store.com/util"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
)

var (
	esClient go_micro_service_search.SearchService
)

func Init(service micro.Service) {
	esClient = go_micro_service_search.NewSearchService("go.micro.service.search", service.Client())
}

//DoUploadHandler 处理文件上传
func DoUploadHandler(c *gin.Context) {
	errCode := 0
	log.Printf("Handling upload...")
	defer func() {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		log.Printf("ErrCode:%v", errCode)
		if errCode < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "上传失败",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "上传成功",
			})
		}
	}()

	//1.从form表单中获得内容句柄
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Error: getting file handler from http.request, err:%v", err)
		errCode = -1
		return
	}

	//2.把文件内容转化为[]byte
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Printf("Error: reading from file:%v into buffer, err:%v", file, err)
		errCode = -2
		return
	}

	//3.构建文件元信息
	fileMeta := dbClient.FileMeta{
		FileName: head.Filename,
		FileSha1: util.Sha1(buf.Bytes()),
		FileSize: int64(len(buf.Bytes())),
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	//4.将文件写入临时存储
	fileMeta.Location = config.TempLocalRootDir + fileMeta.FileSha1
	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		log.Printf("Error: creating local file：%v, err:%v", fileMeta.Location, err)
		errCode = -3
		return
	}
	defer newFile.Close()

	nBytes, err := newFile.Write(buf.Bytes())
	if err != nil || int64(nBytes) != fileMeta.FileSize {
		log.Printf("Error: saving file into local file, err:%v", err)
		errCode = -4
		return
	}

	//5.同步或者异步将文件转移到Ceph/OSS, 使用协程，减少http请求响应时间
	go func() {
		newFile.Seek(0, 0)
		//存储在ceph集群
		if config.CurrentStoreType == common.StoreCeph {
			cephPath := config.CephRootDir + fileMeta.FileSha1
			err := ceph.PutObject(config.CephBucket, cephPath, buf.Bytes())
			if err != nil {
				log.Printf("Error: putting:%v in ceph, err:%v", cephPath, err)
				return
			}
			fileMeta.Location = cephPath
		} else if config.CurrentStoreType == common.StoreOSS {
			ossPath := config.OSSRootDir + fileMeta.FileSha1
			if !config.AsyncTransferEnabled {
				err := oss.Bucket().PutObject(ossPath, bytes.NewReader(buf.Bytes()))
				if err != nil {
					log.Printf("Error: putting object:%s, err:%v", ossPath, err)
					return
				}
				fileMeta.Location = ossPath
			} else {
				msg := mq.TransferData{
					FileHash:      fileMeta.FileSha1,
					FileName:      fileMeta.FileName,
					CurLocation:   fileMeta.Location,
					DestLocation:  ossPath,
					DestStoreType: common.StoreOSS,
				}
				msgBytes, _ := json.Marshal(msg)
				success := mq.Publish(config.TransExchangeName, config.TransOSSRoutingKey, msgBytes)
				if !success {
					log.Printf("Error: publishing msg:%v, err:%v", string(msgBytes), err)
					//TODO 转移到失败队列，稍后再试
					return
				}
			}
		}
	}()

	//6.更新文件表
	_, err = dbClient.OnFileUploadFinished(fileMeta)
	if err != nil {
		log.Printf("Error: updating file table, filemeta:%v, err:%v", fileMeta, err)
		errCode = -5
		return
	}

	//7.更新用户文件表
	userName := c.Request.FormValue("username")
	updateResult, err := dbClient.OnUserFileUploadFinished(userName, fileMeta)
	if err != nil || !updateResult.Success {
		log.Printf("Error: updating user&file table, username:%v, filemeta:%v, msg:%v, err:%v", userName, fileMeta, updateResult.Msg, err)
		errCode = -6
		return
	}

	//8.更新elasticsearch文档记录
	saveResp, err := esClient.SaveDocument(context.Background(), &go_micro_service_search.ReqSaveDocument{
		Index: "cloudstore",
		Typ:   "user_file",
		UserFileMeta: &go_micro_service_search.UserFileMeta{
			UserName: userName,
			FileSha1: fileMeta.FileSha1,
			FileName: fileMeta.FileName,
			FileSize: fileMeta.FileSize,
			Location: fileMeta.Location,
			UploadAt: fileMeta.UploadAt,
		},
	})
	if err != nil {
		log.Printf("DoUploadHandler Error: failed to save document to es, err:%v", err)
		errCode = -7
		return
	}
	if !saveResp.Success {
		log.Printf("DoUploadHandler Error: failed to save document, msg:%v", saveResp.Message)
		errCode = -8
		return
	}
}

func TryFastUploadHandler(c *gin.Context) {
	fileHash := c.Request.FormValue("filehash")
	result, err := dbClient.GetFileMeta(fileHash)
	if err != nil {
		log.Printf("Error: selecting filehash:%v's file meta, msg:%v, err:%v", fileHash, result.Msg, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if !result.Success {
		log.Printf("Error: selecting filehash:%v's file meta, msg:%v, err:%v", fileHash, result.Msg, err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "上传失败",
		})
		return
	}

	userName := c.Request.FormValue("username")
	fileMeta := dbClient.TableFileToFileMeta(result.Data.(orm.TableFile))
	upDateResult, err := dbClient.OnUserFileUploadFinished(userName, fileMeta)
	if err != nil {
		log.Printf("Error: updating user&file table, err:%v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if !upDateResult.Success {
		log.Printf("Error: updating user&file table, err_msg:%v", upDateResult.Msg)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "上传失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "上传成功",
	})
}
