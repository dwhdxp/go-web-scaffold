package controllers

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"friend-ranking/models"
)

// 返回的用户信息
type UserInfo struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

// 通过结构体对象调用方法，避免一个包中同名函数的冲突
type UserController struct{}

// 注册业务
func (u UserController) Register(c *gin.Context) {
	// 获取form参数
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")
	confirmPassword := c.DefaultPostForm("confirmPassword", "")

	// 参数校验
	// 1.判断是否有未输入项
	if username == "" || password == "" || confirmPassword == "" {
		ReturnError(c, 4001, "注册失败：请输入完整信息")
		return
	}
	// 2.判断输入两次密码是否一致
	if password != confirmPassword {
		ReturnError(c, 4001, "注册失败：两次输入密码不一致")
		return
	}

	// 3.判断用户是否已注册
	user, err := models.GetUserInfoByUsername(username)
	if user.Id != 0 {
		ReturnError(c, 4001, "注册失败：用户已存在")
		return
	}

	// 加密密码并存储到mysql中
	_, err = models.AddUser(username, EncryMd5(password))
	if err != nil {
		ReturnError(c, 4002, "注册失败，请重试")
		return
	}

	ReturnSuccess(c, 0, "register success", "", 1)
}

// 登录业务
func (u UserController) Login(c *gin.Context) {
	// 获取form参数
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")

	// 参数校验
	// 1.判断是否有未输入项
	if username == "" || password == "" {
		ReturnError(c, 4001, "登录失败：请输入完整信息")
		return
	}

	// 2.检查用户名和密码
	user, _ := models.GetUserInfoByUsername(username)
	if user.Id == 0 || user.Username != username || user.Password != EncryMd5(password) {
		ReturnError(c, 4001, "登录失败：用户名或密码错误")
		return
	}

	// 将UserInfo保存到session中
	session := sessions.Default(c)
	session.Set("login:"+strconv.Itoa(user.Id), user.Id)
	// session.Options(sessions.Options{MaxAge: 86400}) 设置session有效时间
	session.Save() // 保存

	data := UserInfo{Id: user.Id, Username: user.Username}
	ReturnSuccess(c, 0, "success", data, 1)
}
