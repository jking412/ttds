package main

import (
	"awesomeProject/internal/handler"
	"awesomeProject/internal/task"
	"awesomeProject/pkg/configs"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	_ "net/http/pprof"
)

var (
	pprofEnabled bool
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the TTDS server",
	Long:  `Start the TTDS web server with optional pprof profiling`,
	Run: func(cmd *cobra.Command, args []string) {
		// 创建 Gin 引擎
		r := gin.Default()

		// 注册路由
		handler.RegisterRoutes(r)

		// 根据参数启动pprof
		if pprofEnabled {
			go func() {
				fmt.Println(http.ListenAndServe("localhost:6060", nil))
			}()
		}

		// 启动TaskServer和TaskClient
		// TODO: 存在一个严重的问题，未知原因会导致一个任务被执行多次，需要解决
		go task.InitTaskServer()
		client := task.InitTaskClient()
		defer client.AsynqClient.Close()

		// 启动服务器
		logrus.Infof("starting server on %s", configs.GetConfig().Server.Address)
		if err := r.Run(configs.GetConfig().Server.Address); err != nil {
			logrus.Fatalf("failed to start server: %v", err)
		}
	},
}

func init() {
	serverCmd.Flags().BoolVarP(&pprofEnabled, "pprof", "p", false, "Enable pprof profiling on port 6060")
	rootCmd.AddCommand(serverCmd)
}
