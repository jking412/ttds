package container

import (
	"awesomeProject/internal/model"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

// 设置Go Docker Client变量
var cli *client.Client
var _ Manager = (*DockerEngine)(nil)
var once sync.Once

func newDockerEngine() *DockerEngine {
	once.Do(func() {
		var err error
		cli, err = client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			logrus.Fatalf("failed to create docker client: %v", err)
		}
	})
	return &DockerEngine{
		cli: cli,
	}
}

type DockerEngine struct {
	cli *client.Client
}

func (d *DockerEngine) CreateContainer(template *model.ContainerTemplate) (*model.ContainerInstance, error) {
	// 创建容器配置
	config := &container.Config{
		Image: template.Image,
		Env:   []string{},
	}

	// 创建主机配置
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{},
		Mounts:       []mount.Mount{},
	}

	// 解析并设置环境变量 (格式: key=value;)
	if template.Envs != "" {
		envPairs := strings.Split(template.Envs, ";")
		for _, pair := range envPairs {
			pair = strings.TrimSpace(pair)
			if pair != "" {
				config.Env = append(config.Env, pair)
			}
		}
	}

	// 设置固定的sudo密码
	config.Env = append(config.Env, "SUDO_PASSWORD=123456")

	// 生成随机token CONNECTION_TOKEN
	token := generateRandomToken()
	config.Env = append(config.Env, "CONNECTION_TOKEN="+token)

	// 解析并设置端口映射 (格式: key:value;)
	if template.Ports != "" {
		portPairs := strings.Split(template.Ports, ";")
		for _, pair := range portPairs {
			pair = strings.TrimSpace(pair)
			if pair != "" {
				parts := strings.Split(pair, ":")
				if len(parts) == 2 {
					hostPort := strings.TrimSpace(parts[0])

					// 创建端口映射
					containerPort, err := nat.NewPort("tcp", parts[1])
					if err != nil {
						return nil, fmt.Errorf("invalid containerPort mapping: %s, error: %v", pair, err)
					}

					hostConfig.PortBindings[containerPort] = []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: hostPort,
						},
					}
				}
			}
		}
	}

	// 解析并设置卷挂载 (格式: key:value;)
	if template.Volumes != "" {
		volumePairs := strings.Split(template.Volumes, ";")
		for _, pair := range volumePairs {
			pair = strings.TrimSpace(pair)
			if pair != "" {
				parts := strings.Split(pair, ":")
				if len(parts) == 2 {
					hostPath := strings.TrimSpace(parts[0])
					containerPath := strings.TrimSpace(parts[1])

					// 创建卷挂载
					hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
						Type:   mount.TypeBind,
						Source: hostPath,
						Target: containerPath,
					})
				}
			}
		}
	}

	// 生成随机容器名称
	containerName := fmt.Sprintf("%s-%s", template.Name, generateRandomString(8))

	// 创建容器
	resp, err := d.cli.ContainerCreate(
		context.Background(),
		config,
		hostConfig,
		&network.NetworkingConfig{},
		nil,
		containerName,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create container: %v", err)
	}

	// 获得ip
	containerInfo, err := d.cli.ContainerInspect(context.Background(), resp.ID)
	if err != nil {
		logrus.Warnf("failed to inspect container: %v", err)
		return nil, err
	}

	// 获取容器IP地址
	var ipAddress string
	for _, net := range containerInfo.NetworkSettings.Networks {
		if net.IPAddress != "" {
			ipAddress = net.IPAddress
			break
		}
	}

	// 创建容器实例记录
	instance := &model.ContainerInstance{
		TemplateID:  template.ID,
		ContainerID: resp.ID,
		Name:        containerName,
		Status:      "Pending", // 初始状态为Pending
		StartAt:     time.Now(),
		Token:       token,
		IPAddress:   ipAddress,
	}

	return instance, nil
}

// 生成随机字符串
func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		logrus.Errorf("failed to generate random string: %v", err)
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// 生成随机token
func generateRandomToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	// 生成 32 字节的随机数
	if _, err := rand.Read(b); err != nil {
		logrus.Errorf("failed to generate 32-char random string: %v", err)
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	for i := range b {
		// 将随机字节映射到字符集
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

func (d *DockerEngine) StartContainer(instance *model.ContainerInstance) error {
	// 启动容器
	err := d.cli.ContainerStart(context.Background(), instance.ContainerID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	// 更新容器状态
	instance.Status = "Running"

	// 获取容器信息以获取IP地址
	containerInfo, err := d.cli.ContainerInspect(context.Background(), instance.ContainerID)
	if err != nil {
		logrus.Warnf("failed to inspect container: %v", err)
	} else {
		// 获取容器IP地址
		for _, net := range containerInfo.NetworkSettings.Networks {
			instance.IPAddress = net.IPAddress
			break
		}
	}

	return nil
}

func (d *DockerEngine) StopContainer(instance *model.ContainerInstance) error {
	// 设置超时时间（10秒）
	timeout := 10

	// 停止容器
	err := d.cli.ContainerStop(context.Background(), instance.ContainerID, container.StopOptions{
		Timeout: &timeout,
	})

	if err != nil {
		return fmt.Errorf("failed to stop container: %v", err)
	}

	// 更新容器状态
	instance.Status = "Stopped"
	instance.EndAt = time.Now()

	return nil
}

func (d *DockerEngine) RemoveContainer(instance *model.ContainerInstance) error {
	// 设置移除选项
	options := container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	// 移除容器
	err := d.cli.ContainerRemove(context.Background(), instance.ContainerID, options)
	if err != nil {
		return fmt.Errorf("failed to remove container: %v", err)
	}

	// 更新容器状态
	instance.Status = "Removed"
	instance.EndAt = time.Now()

	return nil
}

func (d *DockerEngine) Exists(containerName string) (bool, error) {
	// 设置过滤器
	filter := filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})

	// 列出符合条件的容器
	containers, err := d.cli.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filter,
	})

	if err != nil {
		return false, fmt.Errorf("failed to list containers: %v", err)
	}

	// 检查是否存在匹配的容器
	return len(containers) > 0, nil
}

func (d *DockerEngine) ExecCommand(instance *model.ContainerInstance, script *model.ContainerScript) error {
	// 创建执行配置
	execConfig := container.ExecOptions{
		Cmd:          []string{"/bin/sh", "-c", script.Content},
		AttachStdout: true,
		AttachStderr: true,
	}

	// 创建执行实例
	execID, err := d.cli.ContainerExecCreate(context.Background(), instance.ContainerID, execConfig)
	if err != nil {
		return fmt.Errorf("failed to create exec instance: %v", err)
	}

	// 启动执行实例
	err = d.cli.ContainerExecStart(context.Background(), execID.ID, container.ExecStartOptions{
		Detach: true,
	})
	if err != nil {
		return fmt.Errorf("failed to start exec instance: %v", err)
	}

	// 设置整体超时
	timeout := time.Duration(script.Timeout) * time.Second
	deadline := time.After(timeout)

	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-deadline:
			return fmt.Errorf("execution timeout after %s", timeout)
		case <-tick.C:
			inspect, err := d.cli.ContainerExecInspect(context.Background(), execID.ID)
			if err != nil {
				return fmt.Errorf("failed to inspect exec instance: %v", err)
			}

			if !inspect.Running {
				// 脚本执行完成
				if inspect.ExitCode != 0 {
					return fmt.Errorf("script failed with exit code %d", inspect.ExitCode)
				}
				return nil
			}
		}
	}

}
