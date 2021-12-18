package userinfo

import (
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/model"
	"trade/utils"
)

// RetrieveUserInfo 查询用户信息
func RetrieveUserInfo(c *gin.Context) {
	utils.SugarLogger.Info("RetrieveUserInfo Call")
	var info model.Session
	session := sessions.Default(c)
	session.Get(info)
	userinfo, err := model.RetrieveUserInfo(info.User)
	if err != nil {
		utils.SugarLogger.Error("RetrieveUserInfo Failed @BuyProduct", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
		return
	}
	data, _ := json.Marshal(userinfo)
	c.JSON(http.StatusOK, gin.H{
		"message": "成功",
		"data":    data,
	})
}

// RetrieveUserOrders 查询用户订单
func RetrieveUserOrders(c *gin.Context) {
	utils.SugarLogger.Info("RetrieveUserOrders Call")
	var info model.Session
	session := sessions.Default(c)
	session.Get(info)
	orders, err := model.RetrieveOrder(info.User)
	if err != nil {
		utils.SugarLogger.Error("RetrieveOrder Failed @BuyProduct", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
		return
	}
	data, _ := json.Marshal(orders)
	c.JSON(http.StatusOK, gin.H{
		"message": "成功",
		"data":    data,
	})
}
