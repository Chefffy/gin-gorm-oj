package middlewares

import (
	"gin-gorm-oj/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthAdminCheck
// 检查用户权限
func AuthAdminCheck() gin.HandlerFunc{
	return func(c *gin.Context) {
		//check if user is admin
		auth := c.GetHeader("Authorization")
		userClaim,err := helper.AnalyseToken(auth)
		if err != nil{
			c.Abort()
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusUnauthorized,
				"msg":"Unauthorized Authorization",
			})
			return
		}
		if userClaim == nil ||userClaim.IsAdmin != 1{
			c.Abort()
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusUnauthorized,
				"msg":"Unauthorized Admin",
			})
			return
		}
		c.Next()
	}
}
