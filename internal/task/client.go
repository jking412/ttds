package task

import (
	"awesomeProject/pkg/configs"
	"awesomeProject/pkg/message"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"sync"
	"time"
)

// TODO: 修改所有的once，加上once的名称
var once sync.Once
var client *Client
var channelCancel map[string]context.CancelFunc

type Client struct {
	AsynqClient *asynq.Client
	message     message.Manager
}

func InitTaskClient() *Client {
	redisAddr := configs.GetConfig().Redis.Host + ":" + configs.GetConfig().Redis.Port
	return newTaskClient(redisAddr)
}

func newTaskClient(redisAddr string) *Client {
	once.Do(func() {
		channelCancel = make(map[string]context.CancelFunc)
		client = &Client{
			AsynqClient: asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr}),
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
	_, err = c.AsynqClient.Enqueue(task, asynq.MaxRetry(3), asynq.Queue("default"))
	if err != nil {
		return err
	}

	channelID := fmt.Sprintf("%d:%d", p.UserID, p.Template.ID)
	ch, err := c.message.CreateChannel(channelID)
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(time.Millisecond * 500)
		defer ticker.Stop()
		defer close(ch)

		timeout := time.After(time.Minute)
		ctx, cancel := context.WithCancel(context.Background())
		channelCancel[channelID] = cancel

		for {
			select {
			case <-ctx.Done():
				return
			case <-timeout:
				return
			case <-ticker.C:
				ch <- fmt.Sprintf("Pending")
			}
		}
	}()

	return nil
}
