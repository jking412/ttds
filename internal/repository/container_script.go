package repository

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/configs"
	"gorm.io/gorm"
	"sync"
)

var (
	scriptRepositoryInstance ContainerScript
	scriptSyncOnce           sync.Once

	_ ContainerScript = (*ContainerScriptImpl)(nil)
)

type ContainerScript interface {
	GetScriptsByTemplateID(templateID uint) ([]*model.ContainerScript, error)
}

func NewContainerScript(db *gorm.DB) ContainerScript {
	scriptSyncOnce.Do(func() {
		if configs.GetConfig().Env == "dev" {
			scriptRepositoryInstance = &DevScriptRepositoryImpl{
				DB: db,
			}
		} else {
			scriptRepositoryInstance = &ContainerScriptImpl{
				DB: db,
			}
		}
	})
	return scriptRepositoryInstance
}

type ContainerScriptImpl struct {
	DB *gorm.DB
}

func (r *ContainerScriptImpl) GetScriptsByTemplateID(templateID uint) ([]*model.ContainerScript, error) {
	var scripts []*model.ContainerScript
	// 根据order升序排序
	// 从数据库中获取脚本
	result := r.DB.Where("template_id =?", templateID).Order("order ASC").Find(&scripts)
	return scripts, result.Error
}

type DevScriptRepositoryImpl struct {
	DB *gorm.DB
}

func (r *DevScriptRepositoryImpl) GetScriptsByTemplateID(templateID uint) ([]*model.ContainerScript, error) {
	// 返回测试用的脚本数据
	scripts := make([]*model.ContainerScript, 0)

	scripts = append(scripts, &model.ContainerScript{
		Order:      1,
		TemplateID: templateID,
		// sleep 5s
		Content: "sleep 5",
		Timeout: 10,
	})

	scripts = append(scripts, &model.ContainerScript{
		Order:      2,
		TemplateID: templateID,
		// sleep 5s
		Content: "sleep 5",
		Timeout: 10,
	})

	scripts = append(scripts, &model.ContainerScript{
		Order:      3,
		TemplateID: templateID,
		// sleep 5s
		Content: "sleep 5",
		Timeout: 10,
	})

	scripts = append(scripts, &model.ContainerScript{
		Order:      4,
		TemplateID: templateID,
		// exit 1
		Content: "exit 1",
		Timeout: 10,
	})

	return scripts, nil
}
