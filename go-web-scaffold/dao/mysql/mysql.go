package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
	"web-scaffold/settings"
)

var db *gorm.DB

// Init 初始化MySQL配置
func Init(cfg *settings.MySQLConfig) (err error) {
	// root:password@tcp(ip:port)/test?charset=utf8&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		zap.L().Error("open mysql failed", zap.Error(err))
	}

	// 定义数据库参数
	db.DB().SetMaxOpenConns(cfg.MaxOpenConns) // 设置数据库连接池最大连接数
	db.DB().SetMaxIdleConns(cfg.MaxIdleConns) // 连接池最大允许空闲连接数

	return
}

func Close() (err error) {
	err = db.Close()
	if err != nil {
		zap.L().Error("close mysql failed", zap.Error(err))
	}
	return
}
