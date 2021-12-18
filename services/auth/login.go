package auth

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"trade/message"
	"trade/model"
	"trade/utils"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	UserKey = "user"
)

// LoginRequired 检查session中间件
func LoginRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserKey)
	if user == nil {
		// 未登录，阻断请求
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// 继续
	c.Next()
}

// Login 用户登录
func Login(c *gin.Context) {
	logInfo := new(message.LoginReq)
	session := sessions.Default(c)
	utils.SugarLogger.Info("Login Call", logInfo)

	//校验参数
	if err := c.ShouldBindWith(&logInfo, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数校验失败"})
	}
	isPass, id, err := model.CheckUserInfo(logInfo.Username, logInfo.Password) // 从数据库检查用户信息
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
		return
	}
	if !isPass {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "认证失败"})
		return
	}

	session.Set(UserKey, id) // 保存用户session
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "认证成功"})
}
