package container

import "awesomeProject/internal/model"

type Manager interface {
	CreateContainer(template *model.ContainerTemplate) (*model.ContainerInstance, error)
	StartContainer(instance *model.ContainerInstance) error
	StopContainer(instance *model.ContainerInstance) error
	RemoveContainer(instance *model.ContainerInstance) error
	Exists(containerName string) (bool, error)
	ExecCommand(instance *model.ContainerInstance, script *model.ContainerScript) (bool, error)
}

func NewManager() Manager {
	return newDockerEngine()
}
