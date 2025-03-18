package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web-scaffold/logger"
)

func SetUpRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置为开发者模式
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "GET ping pong")
	})
	r.POST("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "POST ping pong")
	})

	// 处理其他路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
