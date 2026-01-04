package controller

import (
	"encoding/json"
	"fmt"
	"goflylivechat/models"
	"goflylivechat/tools"
	"goflylivechat/ws"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// 客服登录
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

// 用户登录
func VisitorLogin(c *gin.Context) {
	token := c.PostForm("token")
	avator := ""
	userAgent := c.GetHeader("User-Agent")
	if tools.IsMobile(userAgent) {
		avator = "/static/images/1.png"
	} else {
		avator = "/static/images/2.png"
	}

	//google验证
	ok, err := tools.SiteverifyWithLogin(token)
	if !ok {
		c.JSON(200, gin.H{
			"code":    401,
			"message": "google验证失败",
		})
		return
	}

	toId := c.PostForm("to_id")
	visitor_name := c.PostForm("visitor_name")
	id := tools.Uuid()
	refer := c.PostForm("refer")
	name := "Guest"
	vn, err := url.QueryUnescape(visitor_name)
	if vn != "" && err == nil {
		name = vn
	}
	city := ""
	countryname, cityname := tools.GetCity("./config/GeoLite2-City.mmdb", c.ClientIP())
	if countryname != "" || cityname != "" {
		city = fmt.Sprintf("%s %s", countryname, cityname)
		name = fmt.Sprintf("%s Guest", city)
	}

	client_ip := c.ClientIP()
	extra := c.PostForm("extra")
	extraJson := tools.Base64Decode(extra)
	if extraJson != "" {
		var extraObj VisitorExtra
		err := json.Unmarshal([]byte(extraJson), &extraObj)
		if err == nil {
			if extraObj.VisitorName != "" {
				name = extraObj.VisitorName
			}
			if extraObj.VisitorAvatar != "" {
				avator = extraObj.VisitorAvatar
			}
		}
	}
	//log.Println(name,avator,c.ClientIP(),toId,id,refer,city,client_ip)
	if name == "" || avator == "" || toId == "" || id == "" || refer == "" || client_ip == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "error",
		})
		return
	}
	kefuInfo := models.FindUser(toId)
	if kefuInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "The customer service account does not exist",
		})
		return
	}
	visitor := models.FindVisitorByVistorId(id)
	if visitor.Name != "" {
		avator = visitor.Avator
		//更新状态上线
		models.UpdateVisitor(name, visitor.Avator, id, 1, c.ClientIP(), toId, c.ClientIP(), refer, extra)
	} else {
		models.CreateVisitor(name, avator, c.ClientIP(), toId, id, refer, city, client_ip, extra)
	}
	visitor.Name = name
	visitor.Avator = avator
	visitor.ToId = toId
	visitor.ClientIp = c.ClientIP()
	visitor.VisitorId = id

	//各种通知
	go SendNoticeEmail(visitor.Name, " incoming!")
	//go SendAppGetuiPush(kefuInfo.Name, visitor.Name, visitor.Name+" incoming!")
	go SendVisitorLoginNotice(kefuInfo.Name, visitor.Name, visitor.Avator, visitor.Name+" incoming!", visitor.VisitorId)
	go ws.VisitorOnline(kefuInfo.Name, visitor)
	//go SendServerJiang(visitor.Name, "来了", c.Request.Host)

	// Token generation
	info := map[string]interface{}{
		"name":        visitor.Name,
		"visitor_id":  visitor.VisitorId,
		"create_time": time.Now().Unix(),
	}
	authToken, err := tools.MakeToken(info)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "登录暂时不可用",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": authToken,
	})
}
