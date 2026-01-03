package controller

import (
	"goflylivechat/models"
	"goflylivechat/tools"
	"goflylivechat/ws"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostKefuAvator(c *gin.Context) {

	avator := c.PostForm("avator")
	if avator == "" {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "不能为空",
			"result": "",
		})
		return
	}
	kefuName, _ := c.Get("kefu_name")
	models.UpdateUserAvator(kefuName.(string), avator)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": "",
	})
}
func PostKefuPass(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	newPass := c.PostForm("new_pass")
	confirmNewPass := c.PostForm("confirm_new_pass")
	old_pass := c.PostForm("old_pass")
	if newPass != confirmNewPass {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "密码不一致",
			"result": "",
		})
		return
	}
	user := models.FindUser(kefuName.(string))
	if user.Password != tools.Md5(old_pass) {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "旧密码不正确",
			"result": "",
		})
		return
	}
	models.UpdateUserPass(kefuName.(string), tools.Md5(newPass))
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": "",
	})
}
func PostKefuClient(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	clientId := c.PostForm("client_id")

	if clientId == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "client_id不能为空",
		})
		return
	}
	models.CreateUserClient(kefuName.(string), clientId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": "",
	})
}
func GetKefuInfo(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	user := models.FindUser(kefuName.(string))
	info := make(map[string]interface{})
	info["avator"] = user.Avator
	info["username"] = user.Name
	info["nickname"] = user.Nickname
	info["role"] = user.Role
	info["enable2FA"] = user.OtpSecret != ""
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": info,
	})
}
func GetKefuInfoAll(c *gin.Context) {
	id, _ := c.Get("kefu_id")
	userinfo := models.FindUserRole("user.avator,user.name,user.id, role.name role_name", id)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "验证成功",
		"result": userinfo,
	})
}
func GetOtherKefuList(c *gin.Context) {
	idStr, _ := c.Get("kefu_id")
	id := idStr.(float64)
	result := make([]interface{}, 0)
	ws.SendPingToKefuClient()
	kefus := models.FindUsers()
	for _, kefu := range kefus {
		if uint(id) == kefu.ID {
			continue
		}

		item := make(map[string]interface{})
		item["name"] = kefu.Name
		item["nickname"] = kefu.Nickname
		item["avator"] = kefu.Avator
		item["status"] = "offline"
		kefu, ok := ws.KefuList[kefu.Name]
		if ok && kefu != nil {
			item["status"] = "online"
		}
		result = append(result, item)
	}
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": result,
	})
}
func PostTransKefu(c *gin.Context) {
	kefuId := c.Query("kefu_id")
	visitorId := c.Query("visitor_id")
	curKefuId, _ := c.Get("kefu_name")
	user := models.FindUser(kefuId)
	visitor := models.FindVisitorByVistorId(visitorId)
	if user.Name == "" || visitor.Name == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "访客或客服不存在",
		})
		return
	}
	models.UpdateVisitorKefu(visitorId, kefuId)
	ws.UpdateVisitorUser(visitorId, kefuId)
	go ws.VisitorOnline(kefuId, visitor)
	go ws.VisitorOffline(curKefuId.(string), visitor.VisitorId, visitor.Name)
	go ws.VisitorNotice(visitor.VisitorId, "客服转接到"+user.Nickname)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "转移成功",
	})
}
func GetKefuInfoSetting(c *gin.Context) {
	kefuId := c.Query("kefu_id")
	user := models.FindUserById(kefuId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": user,
	})
}
func PostKefuRegister(c *gin.Context) {
	name := c.PostForm("username")
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")
	// role := c.PostForm("role")
	avatar := "/static/images/4.jpg"

	kefu_role, ok := c.Get("kefu_role")
	if !ok || kefu_role != "1" {
		c.JSON(200, gin.H{
			"code":   403,
			"msg":    "没有权限",
			"result": nil,
		})
		return
	}

	if name == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":   400,
			"msg":    "All fields are required",
			"result": nil,
		})
		return
	}

	existingUser := models.FindUser(name)
	if existingUser.Name != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":   409,
			"msg":    "Username already exists",
			"result": nil,
		})
		return
	}

	userID := models.CreateUser(name, tools.Md5(password), avatar, nickname, "0")
	if userID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   500,
			"msg":    "Registration Failed",
			"result": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Registration successful",
		"result": gin.H{
			"user_id": userID,
		},
	})
}
func PostKefuInfo(c *gin.Context) {
	id := c.PostForm("id")
	name := c.PostForm("username")
	password := c.PostForm("password")
	avator := c.PostForm("avator")
	role := c.PostForm("role")
	nickname := c.PostForm("nickname")
	if password != "" {
		password = tools.Md5(password)
	}

	kf_id := c.GetString("kefu_id")

	kefu_role := c.GetString("kefu_role")
	if kefu_role != "1" && kf_id != id {
		c.JSON(200, gin.H{
			"code":   403,
			"msg":    "没有权限",
			"result": nil,
		})
		return
	}

	if name == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "客服账号不能为空",
		})
		return
	}
	models.UpdateUser(id, name, password, avator, nickname, role)

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": "",
	})
}
func GetKefuList(c *gin.Context) {
	kefu_role, ok := c.Get("kefu_role")
	if !ok || kefu_role != "1" {
		c.JSON(200, gin.H{
			"code":   403,
			"msg":    "没有权限",
			"result": nil,
		})
		return
	}
	users := models.FindUsers()
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "获取成功",
		"result": users,
	})
}
func DeleteKefuInfo(c *gin.Context) {
	kefuId := c.Query("id")

	kefu_role, ok := c.Get("kefu_role")
	if !ok || kefu_role != "1" {
		c.JSON(200, gin.H{
			"code":   403,
			"msg":    "没有权限",
			"result": nil,
		})
		return
	}

	kefu_id, ok := c.Get("kefu_id")
	if ok && kefu_id == kefuId {
		c.JSON(200, gin.H{
			"code":   403,
			"msg":    "不能删除当前登录账号",
			"result": nil,
		})
		return
	}

	models.DeleteUserById(kefuId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "删除成功",
		"result": "",
	})
}

func Enable2FA(c *gin.Context) {
	username := c.PostForm("username")
	enable2FA := c.PostForm("enable2FA")

	user := models.FindUser(username)

	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":   403,
			"msg":    "账号不存在",
			"result": nil,
		})
		return
	}

	if enable2FA != "true" {
		err := models.UpdateUserOtps(username, "")
		if err != nil {
			c.JSON(200, gin.H{
				"code":   500,
				"msg":    "停止2FA失败",
				"result": nil,
				"err":    err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"code":   200,
			"msg":    "关闭2FA成功",
			"result": "",
		})
		return
	}

	secret, image, err := tools.GenerateOpts(user.Name)

	if err != nil {
		c.JSON(200, gin.H{
			"code":   500,
			"msg":    "开启2FA失败",
			"result": nil,
			"err":    err.Error(),
		})
		return
	}

	err = models.UpdateUserOtps(username, secret)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   500,
			"msg":    "开启2FA失败",
			"result": nil,
			"err":    err.Error(),
		})
		return
	}

	base64Img, err := tools.ImageToBase64(image)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   500,
			"msg":    "开启2FA失败",
			"result": nil,
			"err":    err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "开启2FA成功",
		"result": base64Img,
	})
}
