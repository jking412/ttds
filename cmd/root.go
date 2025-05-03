package main

import (
	"awesomeProject/pkg/configs"
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var rootCmd = &cobra.Command{
	Use:   "ttds",
	Short: "TTDS application",
	Long:  `TTDS is a teaching platform for data science`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 加载配置
		if err := configs.Init(); err != nil {
			logrus.Fatalf("failed to read config: %v", err)
		}

		// 设置日志级别
		log.InitLog()

		// 初始化数据库
		db.InitDB(db.GenerateDsnFromConfig(), &gorm.Config{})
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
