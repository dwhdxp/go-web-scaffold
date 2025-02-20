package models

import (
	"friend-ranking/pkg/logger"
	"time"

	"friend-ranking/dao"
)

type User struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	AddTime    int64  `json:"addTime"`
	UpdateTime int64  `json:"updataTime"`
}

func (User) TableName() string { return "user" }

// 根据用户名查找用户
func GetUserInfoByUsername(username string) (User, error) {
	if dao.Db == nil {
		logger.Error(map[string]interface{}{"Database connection is not initialized": dao.Db.Error})
	}
	var user User
	err := dao.Db.Where("username = ?", username).First(&user).Error
	return user, err
}

// 添加用户到数据库
func AddUser(username string, password string) (int, error) {
	user := User{Username: username, Password: password, AddTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}
	err := dao.Db.Create(&user).Error
	return user.Id, err
}
