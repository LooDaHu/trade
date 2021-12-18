package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	auth2 "trade/services/auth"
	business2 "trade/services/business"
	"trade/services/userinfo"
	"trade/utils"
)

func main() {
	utils.InitLogger()
	r := engine()
	r.Use(gin.Logger())
	if err := engine().Run(":8080"); err != nil {
		utils.SugarLogger.Fatal("Unable to start", err)
	}
}

func engine() *gin.Engine {
	r := gin.New()
	r.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))
	r.POST("/login", auth2.Login)  // 登入
	r.GET("/logout", auth2.Logout) // 登出

	rBusiness := r.Group("/business")
	rBusiness.Use(auth2.LoginRequired)
	{
		rBusiness.POST("/top-up", business2.IncrBalance)                       // 充值
		rBusiness.GET("/product-list", business2.RetrieveAvailableProductList) // 获取商品列表
		rBusiness.POST("/buy", business2.BuyProduct)                           // 购买商品
	}

	rUserInfo := r.Group("/user")
	rUserInfo.Use(auth2.LoginRequired)
	{
		rUserInfo.GET("/info", userinfo.RetrieveUserInfo)     // 获取用户信息
		rUserInfo.GET("/orders", userinfo.RetrieveUserOrders) // 获取用户订单
	}

	return r
}
