package task

import (
	"awesomeProject/pkg/configs"
	"awesomeProject/pkg/message"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var clientSyncOnce sync.Once
var client *Client
var channelCancel map[string]context.CancelFunc

type Client struct {
	AsynqClient *asynq.Client
	message     message.Manager
}

func GetTaskClient() *Client {
	return client
}

func InitTaskClient() *Client {
	redisAddr := configs.GetConfig().Redis.Host + ":" + configs.GetConfig().Redis.Port
	return newTaskClient(redisAddr)
}

func newTaskClient(redisAddr string) *Client {
	clientSyncOnce.Do(func() {
		channelCancel = make(map[string]context.CancelFunc)
		asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
		client = &Client{
			AsynqClient: asynqClient,
			message:     message.NewChannelManager(),
		}
	})
	return client
}

func (c *Client) EnqueueContainerCreateTask(p ContainerCreatePayload) error {
	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeContainerCreate, payload)
	_, err = c.AsynqClient.Enqueue(task, asynq.MaxRetry(1), asynq.Queue("default"))
	if err != nil {
		return err
	}

	channelID := fmt.Sprintf("%d:%d", p.UserID, p.Template.ID)
	// TODO: 考虑ch的泄露问题
	ch, err := c.message.CreateChannel(channelID)
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(time.Millisecond * 500)
		defer ticker.Stop()

		timeout := time.After(time.Minute)
		ctx, cancel := context.WithCancel(context.Background())
		channelCancel[channelID] = cancel

		for {
			select {
			case <-ctx.Done():
				return
			case <-timeout:
				err = c.message.RemoveChannel(channelID)
				if err != nil {
					logrus.Errorf("RemoveChannel failed: %v", err)
				}
				return
			case <-ticker.C:
				ch <- fmt.Sprintf("Pending")
			}
		}
	}()

	return nil
}
