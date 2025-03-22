package repository

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/db"
)

// CreateChapter 创建一个新的章节
func CreateChapter(chapter *model.Chapter) error {
	return db.DB.Create(chapter).Error
}

// GetChapterByID 根据章节 ID 获取章节信息
func GetChapterByID(id uint) (*model.Chapter, error) {
	var chapter model.Chapter
	result := db.DB.Preload("Sections").First(&chapter, id)
	return &chapter, result.Error
}

// GetChaptersByCourseID 根据课程 ID 获取所有章节信息
func GetChaptersByCourseID(courseID uint) ([]model.Chapter, error) {
	var chapters []model.Chapter
	result := db.DB.Preload("Sections").Where("course_id = ?", courseID).Find(&chapters)
	return chapters, result.Error
}
