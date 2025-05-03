package task

import "awesomeProject/internal/model"

type ContainerCreatePayload struct {
	Template model.ContainerTemplate
	UserID   uint
}
