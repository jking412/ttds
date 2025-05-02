package app

import (
	"awesomeProject/pkg/docker"
)

var cli = docker.DockerClient

//const (
//	systemCallDir = "/home/skynesser/code/system_call"
//)
//
//func StartContainer(c *gin.Context) {
//
//	exists := db.Cache.Exists(context.Background(), "container:ssh")
//	val, err := exists.Result()
//	if val == 1 {
//		c.JSON(200, gin.H{"message": "Container created and started successfully(cache)"})
//		return
//	}
//
//	options := filters.NewArgs()
//	options.Add("name", "os")
//
//	// Check if a container named "os" already exists
//	containers, err := cli.ContainerList(context.Background(), container.ListOptions{Filters: options, All: true})
//
//	if err != nil {
//		// 打印err
//		logrus.Errorf("failed to list containers: %v", err)
//		c.JSON(500, gin.H{"error": err.Error()})
//		return
//	}
//
//	if len(containers) == 0 {
//		// Set up the container configuration
//		// sudo docker run --name os -v $(pwd):/root -p 8000:8080 -d os:base
//		config := &container.Config{
//			Image: "os:ssh",
//		}
//
//		// Set up the container host configuration
//		hostConfig := &container.HostConfig{
//			PortBindings: nat.PortMap{
//				"8080/tcp": []nat.PortBinding{
//					{
//						HostIP:   "0.0.0.0",
//						HostPort: "8000",
//					},
//					{
//						HostIP:   "::",
//						HostPort: "8000",
//					},
//				},
//				"22/tcp": []nat.PortBinding{
//					{
//						HostIP:   "0.0.0.0",
//						HostPort: "22",
//					},
//					{
//
//						HostIP:   "::",
//						HostPort: "22",
//					},
//				},
//			},
//			Mounts: []mount.Mount{
//				{
//					Type:   mount.TypeBind,
//					Source: systemCallDir,
//					Target: "/root",
//				},
//			},
//		}
//
//		// Create the container
//		resp, err := cli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, "os")
//		if err != nil {
//			// 打印resp和err
//			logrus.Errorf("failed to create container: %v, resp: %v", err, resp)
//			c.JSON(500, gin.H{"error": err.Error()})
//			return
//		}
//
//		err = cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
//
//		if err != nil {
//			// 打印err
//			logrus.Errorf("failed to start container: %v", err)
//			c.JSON(500, gin.H{"error": err.Error()})
//			return
//		}
//	} else {
//		// Container already exists, start it
//		err = cli.ContainerStart(context.Background(), containers[0].ID, container.StartOptions{})
//		if err != nil {
//			// 打印err
//			logrus.Errorf("failed to start container: %v", err)
//			c.JSON(500, gin.H{"error": err.Error()})
//			return
//		}
//	}
//
//	cmd := db.Cache.Set(context.Background(), "container:ssh", 1, time.Minute)
//	if cmd.Err() != nil {
//		logrus.Errorf("failed to set cache: %v", cmd.Err())
//		c.JSON(500, gin.H{"error": cmd.Err()})
//		return
//	}
//
//	c.JSON(200, gin.H{"message": "Container created and started successfully"})
//}
