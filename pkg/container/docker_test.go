package container

import (
	"awesomeProject/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 创建测试用的容器模板
func createTestTemplate() *model.ContainerTemplate {
	return &model.ContainerTemplate{
		Name:        "test-container",
		Description: "测试容器",
		Image:       "os:test",
		Ports:       "3001:3000",
	}
}

// 测试启动容器
func TestDockerEngine_StartContainer(t *testing.T) {
	// 跳过实际执行，除非明确要求进行集成测试
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 创建Docker引擎
	docker := newDockerEngine()

	// 创建测试模板和容器
	template := createTestTemplate()
	instance, err := docker.CreateContainer(template)
	assert.NoError(t, err, "创建容器应该成功")

	// 先停止容器
	err = docker.StopContainer(instance)
	assert.NoError(t, err, "停止容器应该成功")
	assert.Equal(t, "Stopped", instance.Status, "容器状态应为Stopped")

	// 启动容器
	err = docker.StartContainer(instance)

	// 断言
	assert.NoError(t, err, "启动容器应该成功")
	assert.Equal(t, "Running", instance.Status, "容器状态应为Running")
	assert.NotEmpty(t, instance.IPAddress, "容器IP地址不应为空")

	// 清理：移除测试容器
	defer docker.RemoveContainer(instance)
}

// 测试检查容器是否存在
func TestDockerEngine_Exists(t *testing.T) {

	// 创建Docker引擎
	docker := newDockerEngine()

	// 创建测试模板和容器
	template := createTestTemplate()
	instance, err := docker.CreateContainer(template)
	assert.NoError(t, err, "创建容器应该成功")

	// 检查容器是否存在
	exists, err := docker.Exists(instance.Name)

	// 断言
	assert.NoError(t, err, "检查容器存在性应该成功")
	assert.True(t, exists, "容器应该存在")

	// 清理：移除测试容器
	defer docker.RemoveContainer(instance)
}

// 测试在容器中执行命令
func TestDockerEngine_ExecCommand(t *testing.T) {

	// 创建Docker引擎
	docker := newDockerEngine()

	// 创建测试模板和容器
	template := createTestTemplate()
	instance, err := docker.CreateContainer(template)
	assert.NoError(t, err, "创建容器应该成功")

	// 启动容器
	err = docker.StartContainer(instance)
	assert.NoError(t, err, "启动容器应该成功")
	assert.Equal(t, "Running", instance.Status, "容器状态应为Running")

	// 创建测试脚本
	script := &model.ContainerScript{
		Content: "echo 'Hello, Docker!'",
		Timeout: 10,
	}

	// 执行命令
	err = docker.ExecCommand(instance, script)

	// 断言
	assert.NoError(t, err, "执行命令应该成功")

	// 测试超时情况
	timeoutScript := &model.ContainerScript{
		Content: "sleep 50",
		Timeout: 1, // 1秒超时，而命令需要5秒
	}

	// 执行超时命令
	err = docker.ExecCommand(instance, timeoutScript)

	// 断言应该超时
	assert.Error(t, err, "执行超时命令应该失败")
	assert.Contains(t, err.Error(), "timeout", "错误信息应该包含超时信息")

	// 清理：移除测试容器
	defer docker.RemoveContainer(instance)
}
