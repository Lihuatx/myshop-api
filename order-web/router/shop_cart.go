package router

import (
	"github.com/gin-gonic/gin"
	"myshop_api/order-web/api/shop_cart"
	"myshop_api/order-web/middlewares"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())
	{
		GoodsRouter.GET("", shop_cart.List)          //购物车列表
		GoodsRouter.DELETE("/:id", shop_cart.Delete) //删除条目
		GoodsRouter.POST("", shop_cart.New)          //添加商品到购物车
		GoodsRouter.PATCH("/:id", shop_cart.Update)  //修改条目
	}
}
