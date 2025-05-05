package usecase

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/task"
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/message"
	"errors"
	"sync"
)

const (
	ContainerCreate = iota
	ContainerExec
)

var (
	containerServiceInstance ContainerService
	containerSyncOnce        sync.Once
	_                        ContainerService = (*ContainerServiceImpl)(nil)
)

type ContainerService interface {
	CreateContainer(userID, templateID uint) error
	GetContainer(userID, templateID uint) (*model.ContainerInstance, error)
	GetChannel(userID, templateID uint, typ int) (chan string, error)
	CheckContainer(userID, templateID uint) error
}

type ContainerServiceImpl struct {
	instanceRepo   repository.InstanceRepository
	templateRepo   repository.TemplateRepository
	scriptRepo     repository.ContainerScript
	taskClient     *task.Client
	messageManager message.Manager
}

func NewContainerService() ContainerService {
	containerSyncOnce.Do(func() {
		containerServiceInstance = &ContainerServiceImpl{
			instanceRepo:   repository.NewInstanceRepository(db.DB),
			templateRepo:   repository.NewTemplateRepository(db.DB),
			scriptRepo:     repository.NewContainerScript(db.DB),
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

func (s *ContainerServiceImpl) GetChannel(userID, templateID uint, typ int) (chan string, error) {
	if typ == ContainerCreate {
		return s.messageManager.GetChannel(task.ContainerCreateChannelName(userID, templateID))
	} else if typ == ContainerExec {
		return s.messageManager.GetChannel(task.ContainerExecChannelName(userID, templateID))
	}
	return nil, errors.New("invalid type")
}

func (s *ContainerServiceImpl) CheckContainer(userID, templateID uint) error {
	instance, err := s.instanceRepo.GetInstanceByUserIDAndTemplateID(templateID, userID)
	if err != nil {
		return err
	}

	scripts, err := s.scriptRepo.GetScriptsByTemplateID(templateID)
	if err != nil {
		return err
	}

	scriptSlice := make([]model.ContainerScript, 0)
	for _, script := range scripts {
		scriptSlice = append(scriptSlice, *script)
	}

	payload := task.ContainerExecPayload{
		Instance:   *instance,
		Scripts:    scriptSlice,
		UserID:     userID,
		TemplateID: templateID,
	}

	return s.taskClient.EnqueueContainerExecTask(payload)
}
