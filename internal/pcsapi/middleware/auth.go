package middleware_auth

import (
	"crypto/md5"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine
var Router *gin.RouterGroup

// TODO: 添加登录中间件
func Init(auth bool, username string, passward string) {
	Engine = gin.Default()
	if auth {
		Router = Engine.Group("api", gin.BasicAuth(
			gin.Accounts{
				username: passward,
			},
		))
		Router.GET("/secrets", func(ctx *gin.Context) {
			user := ctx.MustGet(gin.AuthUserKey).(string)
			if user == username {
				ctx.JSON(http.StatusOK, gin.H{
					"user":   username,
					"secret": md5.Sum([]byte(username)),
				})
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"user":   username,
					"secret": "NO SECRET :(",
				})
			}
		})
	} else {
		Router = Engine.Group("api")
	}
}
