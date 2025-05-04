package repository

import (
	"awesomeProject/internal/model"
	"gorm.io/gorm"
	"sync"
)

var (
	instanceRepositoryInstance InstanceRepository
	instanceSyncOnce           sync.Once

	_ InstanceRepository = (*InstanceRepositoryImpl)(nil)
)

type InstanceRepository interface {
	CreateInstance(*model.ContainerInstance) error
	GetInstanceByUserIDAndTemplateID(userID, templateID uint) (*model.ContainerInstance, error)
}

func NewInstanceRepository(db *gorm.DB) InstanceRepository {
	instanceSyncOnce.Do(func() {
		instanceRepositoryInstance = &InstanceRepositoryImpl{
			DB: db,
		}
	})
	return instanceRepositoryInstance
}

type InstanceRepositoryImpl struct {
	DB *gorm.DB
}

func (r *InstanceRepositoryImpl) GetInstanceByUserIDAndTemplateID(userID, templateID uint) (*model.ContainerInstance, error) {
	var instance model.ContainerInstance
	result := r.DB.Where("user_id = ? AND template_id = ?", userID, templateID).First(&instance)
	return &instance, result.Error
}

func (r *InstanceRepositoryImpl) CreateInstance(instance *model.ContainerInstance) error {
	return r.DB.Create(instance).Error
}
