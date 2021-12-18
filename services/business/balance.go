package business

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"trade/message"
	"trade/model"
	"trade/utils"
)

// IncrBalance 添加余额
func IncrBalance(c *gin.Context) {
	req := new(message.AddBalanceReq)
	utils.SugarLogger.Info("IncrBalance Call", req)
	//校验参数
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		utils.SugarLogger.Error("Params Check Failed @IncrBalance", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数校验失败"})
	}
	// 写锁由MongoDB实现
	err := model.ChangeBalance(context.TODO(), req.TargetID, req.Value)
	if err != nil {
		utils.SugarLogger.Error("ChangeBalance Failed @IncrBalance", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "成功",
	})
}
