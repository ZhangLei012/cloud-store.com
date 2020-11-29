package route

import (
	"net/http"
	"strings"

	"cloud-store.com/asset"
	"cloud-store.com/service/apigateway/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	assetfs "github.com/moxiaomomo/go-bindata-assetfs"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix, filePath string) bool {
	if p := strings.TrimPrefix(filePath, prefix); len(p) < len(filePath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: asset.AssetInfo,
		Prefix:    root,
	}
	return &binaryFileSystem{
		fs: fs,
	}
}

//Router: 网关api路由
func Router() *gin.Engine {
	router := gin.Default()

	//将静态文件打包到bin文件
	router.Use(static.Serve("/static/", BinaryFileSystem("static")))
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"}, // []string{"http://localhost:8080"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "content-Type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition"},
		// AllowCredentials: true,
	}))
	//注册
	router.GET("/user/signup", handler.SignUpHandler)
	router.POST("/user/signup", handler.DoSignUpHandler)
	//登录
	router.GET("/user/signin", handler.SignInhandler)
	router.POST("/user/signin", handler.DoSignInHandler)

	//用户鉴权
	router.Use(handler.Authorize())

	//获取用户信息
	router.POST("/user/info", handler.UserInfoHandler)

	//查询用户的文件
	router.POST("/file/query", handler.FileQueryHandler)

	//重命名用户文件
	router.POST("/file/update", handler.FileMetaUpdateHandler)

	return router
}
