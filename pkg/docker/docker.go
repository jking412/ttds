package docker

import (
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// 设置Go Docker Client变量
var DockerClient *client.Client

func init() {
	var err error
	DockerClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logrus.Fatalf("failed to create docker client: %v", err)
	}
}
