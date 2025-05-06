package repository

import (
	"awesomeProject/internal/model"
	"gorm.io/gorm"
	"sync"
)

var (
	templateRepositoryInstance TemplateRepository
	templateSyncOnce           sync.Once

	_ TemplateRepository = (*TemplateRepositoryImpl)(nil)
)

type TemplateRepository interface {
	GetTemplateByID(id uint) (*model.ContainerTemplate, error)
}

func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	templateSyncOnce.Do(func() {
		templateRepositoryInstance = &TemplateRepositoryImpl{
			DB: db,
		}
	})
	return templateRepositoryInstance
}

type TemplateRepositoryImpl struct {
	DB *gorm.DB
}

func (r *TemplateRepositoryImpl) GetTemplateByID(id uint) (*model.ContainerTemplate, error) {
	var template model.ContainerTemplate
	result := r.DB.First(&template, id)
	return &template, result.Error
}

type DevRepositoryImpl struct {
	DB *gorm.DB
}

func (r *DevRepositoryImpl) GetTemplateByID(id uint) (*model.ContainerTemplate, error) {
	// 不管id，直接返回一个测试用的模板
	return &model.ContainerTemplate{
		Model: gorm.Model{
			ID: id,
		},
		Name:        "test-container",
		Description: "测试容器",
		Image:       "os:test",
		Ports:       "3001:3000",
	}, nil
}
