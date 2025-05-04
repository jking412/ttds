package usecase

import (
	"awesomeProject/internal/task"
	"awesomeProject/pkg/configs"
	"awesomeProject/pkg/db"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

func TestContainerService(t *testing.T) {

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

	// 启动TaskServer和TaskClient
	go task.InitTaskServer()
	client := task.InitTaskClient()
	defer client.AsynqClient.Close()

	service := NewContainerService()

	// 测试数据
	userID := uint(1)
	templateID := uint(1)

	// 1. 创建容器
	err := service.CreateContainer(userID, templateID)
	if err != nil {
		t.Fatalf("创建容器失败: %v", err)
	}

	// 2. 等待几秒
	time.Sleep(3 * time.Second)

	// 3. 获取容器信息并打印
	instance, err := service.GetContainer(userID, templateID)
	if err != nil {
		t.Fatalf("获取容器信息失败: %v", err)
	}

	t.Logf("容器信息: %+v", instance)
}
