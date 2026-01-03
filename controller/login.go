package controller

import (
	"goflylivechat/models"
	"goflylivechat/tools"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginCheckPass(c *gin.Context) {
	token := c.PostForm("token")
	password := c.PostForm("password")
	username := c.PostForm("username")
	otpCode := c.PostForm("otpCode")

	ok, err := tools.SiteverifyWithLogin(token)
	if !ok {
		c.JSON(200, gin.H{
			"code":    401,
			"message": "google验证失败",
		})
		return
	}
	info := models.FindUser(username)
	// Authentication failed case
	if info.Name == "" || info.Password != tools.Md5(password) {
		c.JSON(200, gin.H{
			"code":    401,
			"message": "账号或密码错误",
		})
		return
	}

	if info.OtpSecret != "" {
		//需要验证2fa
		secret, err := tools.DecodeOtpsKey(info.OtpSecret)
		if err != nil {
			c.JSON(200, gin.H{
				"code":    401,
				"message": "2FA解析失败",
			})
			return
		}

		ok := tools.ValidateOtps(otpCode, secret)
		if !ok {
			c.JSON(200, gin.H{
				"code":    401,
				"message": "2FA认证失败",
			})
			return
		}
	}

	// Prepare user session data
	userinfo := map[string]interface{}{
		"kefu_name":   info.Name,
		"kefu_id":     info.ID,
		"kefu_role":   info.Role,
		"create_time": time.Now().Unix(),
	}

	// Token generation
	token, err = tools.MakeToken(userinfo)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "登录暂时不可用",
		})
		return
	}

	// Successful response
	c.JSON(200, gin.H{
		"code":    200,
		"message": "登录成功",
		"result": gin.H{
			"token":      token,
			"created_at": userinfo["create_time"],
		},
	})
}
