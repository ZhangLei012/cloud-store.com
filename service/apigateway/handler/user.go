package handler

import (
	"context"
	"log"
	"net/http"

	"cloud-store.com/common"
	"cloud-store.com/config"
	accountproto "cloud-store.com/service/account/proto"
	downloadproto "cloud-store.com/service/download/proto"
	uploadproto "cloud-store.com/service/upload/proto"
	"github.com/gin-gonic/gin"
	ratelimit2 "github.com/juju/ratelimit"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
	"github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit"
)

var (
	userClient     accountproto.UserService
	downloadClient downloadproto.DownloadService
	uploadClient   uploadproto.UploadService
)

func init() {
	consulReg := consul.NewRegistry( //新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
		registry.Addrs("127.0.0.1:8500"),
	)
	//配置容量及qps，主要是通过令牌桶算法进行实现的
	bRate := ratelimit2.NewBucketWithRate(100, 1000) //100为每秒加入桶的令牌的数量，1000为令牌桶的总容量
	service := micro.NewService(
		micro.Name("go.micro.service.apigateway"),
		micro.Flags(common.CustomFlags...),
		micro.WrapClient(ratelimit.NewClientWrapper(bRate, false)), //加入限流功能，false为连接超时不等待
		micro.Registry(consulReg),
		// 熔断定义：在微服务架构中，扇出链路的某个微服务不可用或者响应时间太长时，会进行服务的降级，
		// 进而熔断该节点微服务的调用，快速返回错误的响应信息。
		// 当检测到该节点微服务调用响应正常后，恢复调用链路。（可以联想到我们家里的电表的保险丝，当电压负载过高后，保险丝熔断，确保家里的电器等其他安全。）
		// 如果没有熔断机制可能产生雪崩效应。
		// 雪崩效应: 微服务A调用微服务B和微服务C来完成的，而微服务B又需要调用微服务D，
		// 微服务D又需要调用微服务E。如果在调用的链路上对微服务E的调用，响应时间过长或者服务不可用，
		// 那么对微服务D的调用就会占用越来越多的系统资源，进而引起微服务D的系统崩溃，
		// 微服务D的不可用，又会连锁反应的引起微服务B崩溃，进而微服务A崩溃，最终导致整个应用不可用。这也就是所谓的“雪崩效应”。
		micro.WrapClient(hystrix.NewClientWrapper()), //加入熔断功能
	)

	//初始化，解析命令行参数等
	service.Init()

	cli := service.Client()

	//初始化一个account服务的客户端
	userClient = accountproto.NewUserService("go.micro.service.user", cli)

	//初始化一个upload服务的客户端
	uploadClient = uploadproto.NewUploadService("go.micro.service.upload", cli)

	//初始化一个download服务的客户端
	downloadClient = downloadproto.NewDownloadService("go.micro.service.download", cli)
}

//SignUpHandler 响应注册页面
func SignUpHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

//DoSignUpHandler 响应POST请求
func DoSignUpHandler(c *gin.Context) {
	log.Printf("Req:%v", c.Request)
	//1.解析参数
	userName := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	log.Printf("Info: request to sign up, username:%v", userName)

	//2.rpc调用account服务
	userResp, err := userClient.SignUp(context.Background(), &accountproto.ReqSignUp{
		UserName: userName,
		Password: password,
	})
	log.Printf("Resp:%v, err:%v", userResp, err)
	if err != nil {
		log.Printf("Error: signing up, err:%v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": userResp.Code,
		"msg":  userResp.Message,
	})
}

//SignInhandler 响应登录页面
func SignInhandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

//DoSignInHandler 处理登录请求
func DoSignInHandler(c *gin.Context) {
	userName := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	signInResp, err := userClient.SignIn(context.Background(), &accountproto.ReqSignIn{
		UserName: userName,
		Password: password,
	})

	if err != nil {
		log.Printf("Error: signing in, err:%v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if signInResp.Code != common.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"code": signInResp.Code,
			"msg":  "sign in failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": int(signInResp.Code),
		"msg":  "sign in succeed",
		"data": struct {
			Location      string
			Username      string
			Token         string
			UploadEntry   string
			DownloadEntry string
		}{
			Location:      "/static/view/home.html",
			Username:      userName,
			Token:         signInResp.Token,
			UploadEntry:   config.UploadLBHost,
			DownloadEntry: config.DownloadLBHost,
		},
	})
}

//UserInfoHandler 获取用户信息
func UserInfoHandler(c *gin.Context) {
	userName := c.Request.FormValue("username")
	userInfoResp, err := userClient.UserInfo(context.Background(), &accountproto.ReqUserInfo{
		UserName: userName,
	})
	if err != nil {
		log.Printf("Error: getting user info, err:%v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if userInfoResp.Code != common.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "getting user info failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": int(common.StatusOK),
		"msg":  "",
		"data": gin.H{
			"Username":     userName,
			"SignupAt":     userInfoResp.SignUpAt,
			"LastActiveAt": userInfoResp.LastActiveAt,
		},
	})
}
