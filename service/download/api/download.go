package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"cloud-store.com/common"
	"cloud-store.com/config"
	dbClient "cloud-store.com/service/dbproxy/client"
	"cloud-store.com/store/ceph"
	"cloud-store.com/store/oss"
	"github.com/gin-gonic/gin"
)

//DownloadURLHandler 生成临时下载链接
func DownloadURLHandler(c *gin.Context) {
	fileHash := c.Request.FormValue("filehash")
	dbResp, err := dbClient.GetFileMeta(fileHash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusServerError,
			"msg":  "server internal error",
		})
	}

	fileMeta := dbClient.ToTableFile(dbResp.Data)
	//判断文件是存储在本地，ceph，还是OSS
	if strings.HasPrefix(fileMeta.FileAddr.String, config.TempLocalRootDir) || strings.HasPrefix(fileMeta.FileAddr.String, config.CephRootDir) {
		userName := c.Request.FormValue("username")
		token := c.Request.FormValue("token")
		tmpURL := fmt.Sprintf("http://%s/file/download?username=%v&token=%v&filehash=%v", c.Request.Host, userName, token, fileHash)
		c.Data(http.StatusOK, "application/octet-stream", []byte(tmpURL))
	} else if strings.HasPrefix(fileMeta.FileAddr.String, config.OSSRootDir) {
		//oss下载url
		signedURL := oss.DownloadURL(config.OSSRootDir + fileMeta.FileHash)
		log.Printf("Info: Getting download url:%v", signedURL)
		c.Data(http.StatusOK, "application/octet-stream", []byte(signedURL))
	}
}

//DownloadHandler
func DownloadHandler(c *gin.Context) {

	fileHash := c.Request.FormValue("filehash")
	userName := c.Request.FormValue("username")
	log.Printf("Info: handling download, filehash:%v, username:%v", fileHash, userName)
	fileResp, err := dbClient.GetFileMeta(fileHash)
	userFileResp, err := dbClient.QueryUserFileMeta(userName, fileHash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusServerError,
			"msg":  "server internal error",
		})
	}

	fileMeta := dbClient.ToTableFile(fileResp.Data)
	userFileMeta := dbClient.ToTableUserFile(userFileResp.Data)
	log.Printf("DownloadHandler: downloading, file name:%v", fileMeta.FileName)
	if strings.HasPrefix(fileMeta.FileAddr.String, config.TempLocalRootDir) {
		//文件存储在本地
		log.Printf("DownloadHandler downloading from local, file name:%v", fileMeta.FileName)
		c.FileAttachment(fileMeta.FileAddr.String, userFileMeta.FileName)
	} else if strings.HasPrefix(fileMeta.FileAddr.String, config.CephRootDir) {
		//文件存储在ceph集群
		data, _ := ceph.GetCephBucket(config.CephBucket).Get(fileMeta.FileName.String)
		c.Header("content-disposition", "attachment; filename=\""+userFileMeta.FileName+"\"")
		c.Data(http.StatusOK, "application/octet-stream", data)
	}
}
