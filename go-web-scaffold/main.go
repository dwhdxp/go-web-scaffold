package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-scaffold/dao/mysql"
	"web-scaffold/dao/redis"
	"web-scaffold/logger"
	"web-scaffold/routers"
	"web-scaffold/settings"
)

func main() {
	// 1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}

	// 2.初始化
	// 日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync() // 刷新缓冲区

	// MySQL
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	// Redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	// 3.注册路由
	r := routers.SetUpRouter(settings.Conf.Mode)

	// 4.启动服务(优雅关闭)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	// 开启goroutine启动服务
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号优雅关闭服务
	quit := make(chan os.Signal, 1) // 创建接受信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号；Ctrl+C就是触发系统SIGINT信号
	// Notify将收到的通知转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞等待
	zap.L().Info("Shutdown Server ...")

	// 创建一个5秒超时的context：将未处理完的请求处理完再关闭服务，超过5s就超时退出
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown", zap.Error(err))
	}

}
