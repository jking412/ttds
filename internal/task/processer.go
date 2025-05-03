package task

import (
	"awesomeProject/pkg/container"
	"awesomeProject/pkg/message"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"sync"
	"time"
)

var processOnce sync.Once
var processor *ContainerProcessor

type ContainerProcessor struct {
	containerManager container.Manager
	messageManager   message.Manager
}

func newContainerProcessor() *ContainerProcessor {
	processOnce.Do(func() {
		processor = &ContainerProcessor{
			containerManager: container.NewManager(),
			messageManager:   message.NewChannelManager(),
		}
	})
	return processor
}

func (p *ContainerProcessor) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(TypeContainerCreate, p.handleContainerCreateTask)
}

func (p *ContainerProcessor) handleContainerCreateTask(ctx context.Context, t *asynq.Task) error {

	var payload ContainerCreatePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	channelID := fmt.Sprintf("%d:%d", payload.UserID, payload.Template.ID)

	defer func() {
		cancel := channelCancel[channelID]
		cancel()
		delete(channelCancel, channelID)
	}()

	instance, err := p.containerManager.CreateContainer(&payload.Template)
	if err != nil {
		return err
	}

	err = p.containerManager.StartContainer(instance)
	if err != nil {
		return err
	}

	ch, err := p.messageManager.GetChannel(channelID)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	defer close(ch)

	// 发送这个消息最多 3 秒
	timeout := time.After(3 * time.Second)

	for {
		select {
		case <-timeout:
			return nil
		case <-ticker.C:
			ch <- fmt.Sprintf("Running")
		}
	}

}
