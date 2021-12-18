package auth

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Logout 用户登出
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserKey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "系统错误"})
		return
	}
	session.Delete(UserKey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已成功登出"})
}
