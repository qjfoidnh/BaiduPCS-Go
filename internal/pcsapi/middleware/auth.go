package middleware_auth

import "github.com/gin-gonic/gin"

var Engine *gin.Engine
var Router *gin.RouterGroup

// TODO: 添加登录中间件
func init() {
	Engine = gin.Default()
	Router = Engine.Group("api")
}
