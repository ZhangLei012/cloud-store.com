package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	redisPool "cloud-store.com/cache/redis"
	"cloud-store.com/config"
	dbClient "cloud-store.com/service/dbproxy/client"
	"cloud-store.com/util"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadId   string
	ChunkSize  int
	ChunkCount int
}

func init() {
	os.MkdirAll(config.TempPartRootDir, 0777)
}

func InitMultipartUploadHandler(c *gin.Context) {
	//1.解析用户请求参数
	userName := c.Request.FormValue("username")
	fileHash := c.Request.FormValue("filehash")
	fileSize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		log.Printf("Error: failed to get file size, err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "failed to get parameter file size",
		})
	}

	//2.获取redis连接
	redisConn := redisPool.RedisPool().Get()
	defer redisConn.Close()

	//3.生成分块上传的初始化信息
	uploadInfo := MultipartUploadInfo{
		FileHash:   fileHash,
		FileSize:   fileSize,
		UploadId:   userName + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: (fileSize - 1) / (5 * 1024 * 1024),
	}

	//4.将分块初始化信息写入redis
	redisConn.Do("HSET", "MP_"+uploadInfo.UploadId, "chunkcount", uploadInfo.ChunkCount)
	redisConn.Do("HSET", "MP_"+uploadInfo.UploadId, "filesize", uploadInfo.FileSize)
	redisConn.Do("HSET", "MP"+uploadInfo.UploadId, "chunksize", uploadInfo.ChunkSize)

	//5.将初始化信息返回给客户端
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "OK",
		"data": uploadInfo,
	})
}

// UploadPartHandler 处理文件分块
func UploadPartHandler(c *gin.Context) {
	//1.解析用户请求参数
	uploadId := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("chunkindex")

	//2.获取连接池的一个连接
	redisConn := redisPool.RedisPool().Get()
	defer redisConn.Close()

	//3.获取文件内容
	filePath := config.TempPartRootDir + uploadId + "/" + chunkIndex
	os.MkdirAll(path.Dir(filePath), 0744)
	fd, err := os.Create(filePath)
	if err != nil {
		log.Println("Error: creating file:%v, err:%v", filePath, err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "failed to upload",
		})
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	//4.更新redis缓存状态
	redisConn.Do("HSET", "MP_"+uploadId, "chunkindex_"+chunkIndex, 1)

	//5.返回处理结果给客户端
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "OK",
		"data": nil,
	})
}

// CompleteUploadHandler 通知上传完成
func CompleteUploadHandler(c *gin.Context) {
	//1.请求用户参数
	uploadId := c.Request.FormValue("uploadid")
	userName := c.Request.FormValue("username")
	fileHash := c.Request.FormValue("filehash")
	fileSize := c.Request.FormValue("filesize")
	fileName := c.Request.FormValue("fileName")

	//2.获取redis连接池中的一个连接
	redisConn := redisPool.RedisPool().Get()
	defer redisConn.Close()

	//3.通过uploadid查询redis并判断是否所有分块已上传完成
	data, err := redis.Values(redisConn.Do("HGETALL", "MP_"+uploadId))
	if err != nil {
		log.Printf("Error: get redis values, err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "server internal error",
			"data": nil,
		})
		return
	}

	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i++ {
		k := string(data[i].([]byte))
		v := string(data[i].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chunkindex_") && v == "1" {
			chunkCount++
		}
	}

	if totalCount != chunkCount {
		log.Printf("Error: have not uploaded all parts")
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "have not uploaded all parts",
			"data": nil,
		})
		return
	}

	//4.合并分块
	srcPath := config.TempPartRootDir + uploadId + "/"
	destPath := config.TempLocalRootDir + fileHash
	cmd := fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", srcPath, destPath)
	mergeRes, err := util.ExecLinuxShell(cmd)
	if err != nil {
		log.Printf("Error:failed to merge file, err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "failed to merge file",
			"data": nil,
		})
	}
	log.Printf("Info: %v", mergeRes)

	//5.更新唯一文件表以及用户文件表
	fSize, err := strconv.Atoi(fileSize)
	if err != nil {
		log.Printf("Error: failed to get parameter file size, err:%v", fSize)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "failed to get file size",
			"data": nil,
		})
	}

	fileMeta := dbClient.FileMeta{
		FileSha1: fileHash,
		FileSize: int64(fSize),
		Location: destPath,
		FileName: fileName,
	}
	res, err := dbClient.OnFileUploadFinished(fileMeta)
	if err != nil {
		log.Printf("Error: failed to update file table, err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "internal server error",
			"data": nil,
		})
		return
	}
	if !res.Success {
		log.Printf("Error: failed to update file table, msg:%v", res.Msg)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("failed to upload, msg:%v", res.Msg),
		})
		return
	}

	res, err = dbClient.OnUserFileUploadFinished(userName, fileMeta)
	if err != nil {
		log.Printf("Error: failed to update user&file table, err:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "internal server error",
			"data": nil,
		})
		return
	}
	if !res.Success {
		log.Printf("Error: failed to update user&file table, msg:%v", res.Msg)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("failed to upload, msg:%v", res.Msg),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  "OK",
		"data": nil,
	})
}
