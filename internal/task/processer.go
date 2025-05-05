package task

import (
	"awesomeProject/internal/repository"
	"awesomeProject/pkg/container"
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/message"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var processOnce sync.Once
var processor *ContainerProcessor

type ContainerProcessor struct {
	containerManager   container.Manager
	messageManager     message.Manager
	instanceRepository repository.InstanceRepository
}

func newContainerProcessor() *ContainerProcessor {
	processOnce.Do(func() {
		processor = &ContainerProcessor{
			containerManager:   container.NewManager(),
			messageManager:     message.NewChannelManager(),
			instanceRepository: repository.NewInstanceRepository(db.DB),
		}
	})
	return processor
}

func (p *ContainerProcessor) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(TypeContainerCreate, p.handleContainerCreateTask)
	mux.HandleFunc(TypeContainerExec, p.handleContainerExecTask)
}

func (p *ContainerProcessor) handleContainerCreateTask(ctx context.Context, t *asynq.Task) error {

	var payload ContainerCreatePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	channelID := ContainerCreateChannelName(payload.UserID, payload.Template.ID)
	defer func() {
		cancel := channelCancel[channelID]
		cancel()
		delete(channelCancel, channelID)
		err := p.messageManager.RemoveChannel(channelID)
		if err != nil {
			logrus.Warnf("messageManager.RemoveChannel failed: %v", err)
		}
	}()

	instance, err := p.containerManager.CreateContainer(&payload.Template)
	if err != nil {
		logrus.Warnf("containerManager.CreateContainer failed: %v", err)
		return err
	}

	err = p.containerManager.StartContainer(instance)
	if err != nil {
		logrus.Warnf("containerManager.StartContainer failed: %v", err)
		return err
	}

	instance.UserID = payload.UserID
	instance.TemplateID = payload.Template.ID

	err = p.instanceRepository.CreateInstance(instance)
	if err != nil {
		logrus.Warnf("instanceRepository.CreateInstance failed: %v", err)
		return err
	}

	ch, err := p.messageManager.GetChannel(channelID)
	if err != nil {
		return err
	}

	cancel := channelCancel[channelID]
	cancel()

	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	// 发送这个消息最多 3 秒
	timeout := time.After(3 * time.Second)

	for {
		select {
		case <-timeout:
			return nil
		case <-ticker.C:
			ch <- fmt.Sprintf(runningMessage)
		}
	}

}

func (p *ContainerProcessor) handleContainerExecTask(ctx context.Context, t *asynq.Task) error {
	var payload ContainerExecPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	channelID := ContainerExecChannelName(payload.UserID, payload.TemplateID)
	ch, err := p.messageManager.CreateChannel(channelID)
	if err != nil {
		return err
	}

	defer func() {
		err = p.messageManager.RemoveChannel(channelID)
		if err != nil {
			logrus.Warnf("messageManager.RemoveChannel failed: %v", err)
		}
	}()

	for order, script := range payload.Scripts {
		ok, err := p.containerManager.ExecCommand(&payload.Instance, &script)
		if err != nil {
			logrus.Warnf("containerManager.ExecCommand failed: %d,%v", order, err)
		}
		if !ok {
			ch <- fmt.Sprintf(failMessage)
		} else {
			ch <- fmt.Sprintf(passMessage)
		}
	}

	return nil
}
