package initialize

import (
	"github.com/gin-gonic/gin"
	"myshop_api/order-web/middlewares"
	"myshop_api/order-web/router"
	"net/http"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	//配置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/o/v1")
	router.InitOrderRouter(ApiGroup)
	router.InitShopCartRouter(ApiGroup)

	return Router
}
