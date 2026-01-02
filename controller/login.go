package controller

import (
	"goflylivechat/models"
	"goflylivechat/tools"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginCheckPass(c *gin.Context) {
	token := c.PostForm("token")
	password := c.PostForm("password")
	username := c.PostForm("username")
	info := models.FindUser(username)

	ok, err := tools.SiteverifyWithLogin(token)
	if !ok {
		log.Println("验证:", err)
		c.JSON(200, gin.H{
			"code":    401,
			"message": "google验证失败",
		})
		return
	}
	// Authentication failed case
	if info.Name == "" || info.Password != tools.Md5(password) {
		c.JSON(200, gin.H{
			"code":    401,
			"message": "账号密码必填",
		})
		return
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
