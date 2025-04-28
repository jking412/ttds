package main

import (
	"awesomeProject/internal/handler"
	"awesomeProject/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// 加载配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("failed to read config: %v", err)
	}

	// 初始化数据库
	db.InitDB()
	db.InitRedis()

	// 插入模拟数据
	//utils.InsertMockData()

	// 创建 Gin 引擎
	r := gin.Default()

	// 注册路由
	handler.RegisterRoutes(r)

	// 启动服务器
	logrus.Infof("starting server on %s", viper.GetString("server.address"))
	//if err := r.Run(viper.GetString("server.address")); err != nil {
	//	logrus.Fatalf("failed to start server: %v", err)
	//}
	// https 启动
	if err := r.RunTLS(viper.GetString("server.address"), "cert\\my.crt", "cert\\my.key"); err != nil {
		logrus.Fatalf("failed to start server: %v", err)
	}
}
