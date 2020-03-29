package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go_Ocean/common"
	"go_Ocean/model"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 獲得Authrization header
		tokenString := c.GetHeader("Authorization")

		// 驗證 token farmat
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "權限不足"})
			c.Abort() // 拋棄這次請求
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "權限不足"})
			c.Abort() // 拋棄這次請求
			return
		}

		// 通過驗證後，獲取claims中的userId
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 用戶
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "權限不足"})
			c.Abort() // 拋棄這次請求
			return
		}

		// 用戶存在 將user的訊息寫入上下文
		c.Set("user", user)
		c.Next()

	}
}
