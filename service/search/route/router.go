package route

import (
	"cloud-store.com/service/search/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "content-type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Type", "Content-Disposition"},
	}))

	router.GET("/user/file/search", api.SearchUserFileHandler)
	return router
}
