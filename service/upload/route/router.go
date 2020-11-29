package route

import (
	"cloud-store.com/service/upload/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	//处理静态请求
	router.Static("/static", "./static/")

	//使用gin插件，支持跨域请求
	//跨域请求： 在 HTML 中，<a>, <form>, <img>, <script>, <iframe>, <link> 等标签以及 Ajax 都可以指向一个资源地址，而所谓的跨域请求就是指：当前发起请求的域与该请求指向的资源所在的域不一样。这里的域指的是这样的一个概念：我们认为若协议 + 域名 + 端口号均相同，那么就是同域s。
	//aaa.com向aaa.com/user/info发出请求那么这个请求是同域的，因为资源路径的协议、域名以及端口号与当前域一致（例子中协议名默认为http，端口号默认为80）。
	//但是，如果发起一个资源路径为bbb.com/pay/purchase的 Ajax 请求，那么这个请求就是跨域请求，因为域不一致，与此同时由于安全问题，这种请求会受到同源策略限制。
	//跨源资源共享 Cross-Origin Resource Sharing(CORS)
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"}, // []string{"http://localhost:8080"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "content-Type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition"},
		// AllowCredentials: true,
	}))
	//文件上传接口
	router.POST("/file/upload", api.DoUploadHandler)
	//秒传接口
	router.POST("/file/tryfast", api.TryFastUploadHandler)

	return router
}
