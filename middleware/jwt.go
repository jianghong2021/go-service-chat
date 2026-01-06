package middleware

import (
	"goflylivechat/tools"
	"time"

	"github.com/gin-gonic/gin"
)

func JwtPageMiddleware(c *gin.Context) {
	//暂时不处理
	//token := c.Query("token")
	//userinfo := tools.ParseToken(token)
	//if userinfo == nil {
	//	c.Redirect(302,"/login")
	//	c.Abort()
	//}
}
func JwtApiMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	if token == "" {
		token = c.Query("token")
	}
	userinfo := tools.ParseToken(token)
	if userinfo == nil || userinfo["kefu_name"] == nil || userinfo["create_time"] == nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "验证失败",
		})
		c.Abort()
		return
	}
	createTime := int64(userinfo["create_time"].(float64))
	var expire int64 = 24 * 60 * 60
	nowTime := time.Now().Unix()
	if (nowTime - createTime) >= expire {
		c.JSON(200, gin.H{
			"code": 401,
			"msg":  "token失效",
		})
		c.Abort()
	}
	c.Set("kefu_id", userinfo["kefu_id"])
	c.Set("kefu_name", userinfo["kefu_name"])
	c.Set("kefu_role", userinfo["kefu_role"])

}

func JwtApiVisitorMiddleware(c *gin.Context) {
	token := c.GetHeader("access")
	if token == "" {
		token = c.Query("access")
	}
	userinfo := tools.ParseToken(token)
	if userinfo == nil || userinfo["visitor_id"] == nil || userinfo["create_time"] == nil {
		c.JSON(200, gin.H{
			"code": 401,
			"msg":  "验证失败",
		})
		c.Abort()
		return
	}
	createTime := int64(userinfo["create_time"].(float64))
	var expire int64 = 24 * 60 * 60
	nowTime := time.Now().Unix()
	if (nowTime - createTime) >= expire {
		c.JSON(200, gin.H{
			"code": 401,
			"msg":  "token失效",
		})
		c.Abort()
	}
	c.Set("visitor_id", userinfo["visitor_id"])
	c.Set("name", userinfo["name"])

}
