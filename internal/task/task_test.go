package task

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/configs"
	"awesomeProject/pkg/db"
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

func TestContainerCreateTask(t *testing.T) {
	// TODO: configs应该要能够接收测试环境的配置文件，不过目前至少可以读取默认配置
	configs.Init()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"localhost",
		3306,
		"ttds",
	)
	db.InitDB(dsn, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 👈 打印所有 SQL
	})
	InitTaskClient()
	go InitTaskServer()
	defer client.AsynqClient.Close()
	// sleep 1s，等待服务启动
	time.Sleep(time.Second)
	// 准备测试数据
	template := model.ContainerTemplate{
		Model: gorm.Model{
			ID: 1,
		},
		Name:        "test-container",
		Description: "测试容器",
		Image:       "os:test",
		Ports:       "3001:3000",
	}
	payload := ContainerCreatePayload{
		Template: template,
		UserID:   1,
	}

	err := client.EnqueueContainerCreateTask(payload)
	if err != nil {
		t.Errorf("failed to enqueue container create task: %v", err)
	}

	// 验证消息通道
	channelID := fmt.Sprintf("%d:%d", payload.UserID, payload.Template.ID)
	ch, err := processor.messageManager.GetChannel(channelID)
	if err != nil {
		t.Fatalf("GetChannel failed: %v", err)
	}

	// 测试消息接收，打印出来，最多接收10s
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			t.Fatal("Timeout waiting for messageManager")
		default:
			msg, ok := <-ch
			if !ok {
				t.Fatal("Channel closed unexpectedly")
			}
			t.Logf("Received message: %s", msg)
			// 接收到Running消息，说明容器创建成功
			if msg == "Running" {
				return
			}
		}
	}
}
