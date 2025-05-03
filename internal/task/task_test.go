package task

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/configs"
	"context"
	"fmt"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestContainerCreateTask(t *testing.T) {
	// TODO: configs应该要能够接收测试环境的配置文件，不过目前至少可以读取默认配置
	configs.Init()
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

//func TestContainerCreateTask_InvalidPayload(t *testing.T) {
//	// 创建无效负载
//	invalidPayload := []byte("{invalid}")
//	task := asynq.NewTask(TypeContainerCreate, invalidPayload)
//
//	// 创建处理器
//	processor := &ContainerProcessor{
//		containerManager: container.NewManager(),
//		messageManager:          messageManager.NewChannelManager(),
//	}
//
//	// 测试处理无效负载
//	err := processor.handleContainerCreateTask(context.Background(), task)
//	if err == nil {
//		t.Error("Expected error for invalid payload, got nil")
//	}
//}
//
//func TestContainerCreateTask_ContainerCreationFailed(t *testing.T) {
//	// 准备测试数据
//	Template := model.ContainerTemplate{
//		Name:        "Invalid Template",
//		Description: "Will fail creation",
//		Image:       "invalid-image",
//	}
//	payload := ContainerCreatePayload{
//		Template: Template,
//		UserID:   1,
//	}
//
//	// 创建处理器
//	processor := &ContainerProcessor{
//		containerManager: container.NewManager(),
//		messageManager:          messageManager.NewChannelManager(),
//	}
//
//	// 创建测试任务
//	taskPayload, err := json.Marshal(payload)
//	if err != nil {
//		t.Fatal(err)
//	}
//	task := asynq.NewTask(TypeContainerCreate, taskPayload)
//
//	// 测试容器创建失败
//	err = processor.handleContainerCreateTask(context.Background(), task)
//	if err == nil {
//		t.Error("Expected error for container creation failure, got nil")
//	}
//}
