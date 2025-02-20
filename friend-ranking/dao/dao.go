package dao

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"friend-ranking/config"
	"friend-ranking/pkg/logger"
)

var (
	Db  *gorm.DB
	err error
)

// 连接mysql，并设置连接池
func init() {
	// 初始化数据库连接不能使用短变量声明，因为会创建一个新的局部变量，而不是给全局变量赋值
	// Db, err := gorm.Open("mysql", config.Mysqldb)
	Db, err = gorm.Open("mysql", config.Mysqldb)
	if err != nil {
		logger.Error(map[string]interface{}{"mysql connect error": err.Error()})
	}
	if Db.Error != nil {
		logger.Error(map[string]interface{}{"database error": Db.Error})
	}

	// 设置连接池
	// 设置空闲池中连接的最大数量
	Db.DB().SetMaxIdleConns(10)

	// 设置打开数据库连接的最大数量
	Db.DB().SetMaxOpenConns(100)

	// 设置连接可复用的最大时间
	Db.DB().SetConnMaxLifetime(time.Hour)
}
