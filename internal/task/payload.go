package task

import "awesomeProject/internal/model"

type ContainerCreatePayload struct {
	Template model.ContainerTemplate
	UserID   uint
}

type ContainerExecPayload struct {
	Instance   model.ContainerInstance
	Scripts    []model.ContainerScript
	UserID     uint
	TemplateID uint
}
