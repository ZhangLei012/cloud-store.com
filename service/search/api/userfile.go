package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud-store.com/common"
	"cloud-store.com/service/dbproxy/orm"
	"cloud-store.com/service/search/esclient"
	go_micro_service_search "cloud-store.com/service/search/proto"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
)

func SearchUserFileHandler(c *gin.Context) {
	userName := c.Request.FormValue("username")
	query := c.Request.FormValue("query")
	log.Printf("userName:%v, query:%v", userName, query)
	esQuery := elastic.NewBoolQuery()
	esQuery.Must(elastic.NewMatchPhraseQuery("fileName", query), elastic.NewMatchPhraseQuery("userName", userName))
	searchResult, err := esclient.Client().Search().Index("cloudstore").Type("user_file").Query(esQuery).Do(context.Background())
	if err != nil {
		log.Printf("SearchUserFileHandler Error: failed to search file name, %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	if searchResult.TotalHits() <= 0 {
		log.Printf("SearchUserFileHandler Info: hit zero result")
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusOK,
			"msg":  "OK",
			"data": []byte("[]"),
		})
		return
	}
	var userFiles []orm.TableUserFile
	log.Printf("SearchUserFile info: total hits:%v", searchResult.TotalHits())
	for _, hit := range searchResult.Hits.Hits {
		log.Printf("SearchUserFile: hit:%v", string(*hit.Source))
		var userFileMeta go_micro_service_search.UserFileMeta
		err := json.Unmarshal([]byte(*hit.Source), &userFileMeta)
		if err != nil {
			log.Printf("SearchUserFile Error: failed to unmarshal, %v", err)
			c.JSON(http.StatusOK, gin.H{
				"code": common.StatusServerError,
			})
		}
		userFiles = append(userFiles, orm.TableUserFile{
			UserName: userFileMeta.UserName,
			FileHash: userFileMeta.FileSha1,
			FileName: userFileMeta.FileName,
			FileSize: userFileMeta.FileSize,
			UploadAt: userFileMeta.UploadAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": common.StatusOK,
		"msg":  "OK",
		"data": userFiles,
	})
}
