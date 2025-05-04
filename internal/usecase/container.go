package usecase

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/task"
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/message"
	"fmt"
	"sync"
)

var (
	containerServiceInstance ContainerService
	containerSyncOnce        sync.Once
	_                        ContainerService = (*ContainerServiceImpl)(nil)
)

type ContainerService interface {
	CreateContainer(userID, templateID uint) error
	GetContainer(userID, templateID uint) (*model.ContainerInstance, error)
	GetChannel(userID, templateID uint) (chan string, error)
}

type ContainerServiceImpl struct {
	instanceRepo   repository.InstanceRepository
	templateRepo   repository.TemplateRepository
	taskClient     *task.Client
	messageManager message.Manager
}

func NewContainerService() ContainerService {
	containerSyncOnce.Do(func() {
		containerServiceInstance = &ContainerServiceImpl{
			instanceRepo:   repository.NewInstanceRepository(db.DB),
			templateRepo:   repository.NewTemplateRepository(db.DB),
			taskClient:     task.GetTaskClient(),
			messageManager: message.NewChannelManager(),
		}
	})

	return containerServiceInstance
}

func (s *ContainerServiceImpl) CreateContainer(userID, templateID uint) error {

	// 获取模板信息
	template, err := s.templateRepo.GetTemplateByID(templateID)
	if err != nil {
		return err
	}

	// 创建异步任务
	payload := task.ContainerCreatePayload{
		UserID:   userID,
		Template: *template,
	}

	return s.taskClient.EnqueueContainerCreateTask(payload)
}

func (s *ContainerServiceImpl) GetContainer(userID, templateID uint) (*model.ContainerInstance, error) {
	return s.instanceRepo.GetInstanceByUserIDAndTemplateID(templateID, userID)
}

func (s *ContainerServiceImpl) GetChannel(userID, templateID uint) (chan string, error) {
	channelID := fmt.Sprintf("%d:%d", userID, templateID)
	return s.messageManager.GetChannel(channelID)
}
