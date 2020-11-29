package handler

import (
	"net/http"

	"cloud-store.com/common"
	"cloud-store.com/config"
	"github.com/gin-gonic/gin"
)

func isTokenValid(userName string, token string) bool {
	if len(token) != config.TokenLength {
		return false
	}

	//TODO:判断token的时效性
	//TODO:通过userName查询db中的用户的token
	//TODO:对比两个token是否一致
	return true
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		if len(userName) < 3 || !isTokenValid(userName, token) {
			c.Abort()

			c.JSON(http.StatusOK, gin.H{
				"code": int(common.StatusTokenInvalid),
				"msg":  "invalid token",
			})
			return
		}
		c.Next()
	}
}
