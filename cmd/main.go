package main

import (
	"awesomeProject/internal/handler"
	"awesomeProject/pkg/configs"
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 加载配置
	if err := configs.InitConfig(); err != nil {
		logrus.Fatalf("failed to read config: %v", err)
	}

	// 设置日志级别
	log.InitLog()

	// 初始化数据库
	db.InitDB(db.GenerateDsnFromConfig())
	db.InitRedis()

	// 插入模拟数据
	//utils.InsertMockData()

	// 创建 Gin 引擎
	r := gin.Default()

	// 注册路由
	handler.RegisterRoutes(r)

	// 启动服务器
	logrus.Infof("starting server on %s", configs.GetConfig().Server.Address)
	if err := r.Run(configs.GetConfig().Server.Address); err != nil {
		logrus.Fatalf("failed to start server: %v", err)
	}
	// https 启动
	//if err := r.RunTLS(viper.GetString("server.address"), "cert\\my.crt", "cert\\my.key"); err != nil {
	//	logrus.Fatalf("failed to start server: %v", err)
	//}
}
