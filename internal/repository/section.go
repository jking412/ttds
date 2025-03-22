package repository

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/db"
)

// CreateSection 创建一个新的小节
func CreateSection(section *model.Section) error {
	return db.DB.Create(section).Error
}

// GetSectionByID 根据小节 ID 获取小节信息
func GetSectionByID(id uint) (*model.Section, error) {
	var section model.Section
	result := db.DB.Preload("UserStatus").First(&section, id)
	return &section, result.Error
}

// GetSectionsByChapterID 根据章节 ID 获取所有小节信息
func GetSectionsByChapterID(chapterID uint) ([]model.Section, error) {
	var sections []model.Section
	result := db.DB.Preload("UserStatus").Where("chapter_id = ?", chapterID).Find(&sections)
	return sections, result.Error
}
