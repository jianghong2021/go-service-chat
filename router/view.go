package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitViewRouter(engine *gin.Engine) {
	engine.GET("/", PageIndex)
	engine.GET("/login", PageLogin)
	engine.GET("/pannel", PagePannel)
	engine.GET("/livechat", PageChat)
	engine.GET("/main", PageMain)
	engine.GET("/chat_main", PageChatMain)
	engine.GET("/setting", PageSetting)
	engine.GET("/users", PageUsers)
}

func PageIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "error.html", gin.H{
		"code": 404,
		"msg":  "不存在此页面",
	})
}

// Login page
func PageLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Dashboard
func PagePannel(c *gin.Context) {
	c.HTML(http.StatusOK, "pannel.html", nil)
}

// Admin console
func PageMain(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", nil)
}

// Customer chat interface
func PageChat(c *gin.Context) {
	referralSource := c.Query("refer") // More clear variable name

	if referralSource == "" {
		referralSource = c.Request.Referer()
	}
	if referralSource == "" {
		referralSource = "Direct access" // More natural English
	}

	token := c.Query("Authorization")
	kefu_id := c.Query("kefu_id")
	if token == "" || kefu_id == "" {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"code": 403,
			"msg":  "禁止访问此页面",
		})
		return
	}

	c.HTML(http.StatusOK, "chat_page.html", gin.H{
		"Refer": referralSource, // Keeping original template variable name
	})
}

// Support agent console
func PageChatMain(c *gin.Context) {
	c.HTML(http.StatusOK, "chat_main.html", nil)
}

// Settings
func PageSetting(c *gin.Context) {
	c.HTML(http.StatusOK, "setting.html", nil)
}

// Users
func PageUsers(c *gin.Context) {
	c.HTML(http.StatusOK, "users.html", nil)
}
