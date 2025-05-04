package repository

import (
	"awesomeProject/internal/model"
	"gorm.io/gorm"
	"sync"
)

var (
	courseRepositoryInstance CourseRepository
	courseSyncOnce           sync.Once

	_ CourseRepository = (*CourseRepositoryImpl)(nil)
)

// CourseRepository 定义课程仓库接口
type CourseRepository interface {
	GetCourseByID(id uint) (*model.Course, error)
	GetAllCourses() ([]model.Course, error)
	GetCourseReferencesByCourseID(courseID uint) ([]model.CourseReference, error)
	GetCourseStatusByCourseID(userID, courseID uint) ([]model.UserSectionStatus, error)
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	courseSyncOnce.Do(func() {
		courseRepositoryInstance = &CourseRepositoryImpl{
			DB: db,
		}
	})
	return courseRepositoryInstance
}

type CourseRepositoryImpl struct {
	DB *gorm.DB
}

// GetCourseByID 根据课程 ID 获取课程信息
func (r *CourseRepositoryImpl) GetCourseByID(id uint) (*model.Course, error) {
	var course model.Course
	result := r.DB.Preload("Chapters.Sections").Preload("References").First(&course, id)
	return &course, result.Error
}

// GetAllCourses 获取所有课程信息
func (r *CourseRepositoryImpl) GetAllCourses() ([]model.Course, error) {
	var courses []model.Course
	result := r.DB.Find(&courses)
	return courses, result.Error
}

// GetCourseReferencesByCourseID 根据课程 ID 获取所有参考资料信息
func (r *CourseRepositoryImpl) GetCourseReferencesByCourseID(courseID uint) ([]model.CourseReference, error) {
	var references []model.CourseReference
	result := r.DB.Where("course_id = ?", courseID).Find(&references)
	return references, result.Error
}

// GetCourseStatusByCourseID 根据课程ID获取用户学习状态
func (r *CourseRepositoryImpl) GetCourseStatusByCourseID(userID, courseID uint) ([]model.UserSectionStatus, error) {
	// 1. 查询该课程下所有章节的小节ID
	var sectionIDs []uint
	err := r.DB.Table("sections").
		Select("sections.id").
		Joins("JOIN chapters ON sections.chapter_id = chapters.id").
		Where("chapters.course_id = ?", courseID).
		Scan(&sectionIDs).Error
	if err != nil {
		return nil, err
	}
	if len(sectionIDs) == 0 {
		return []model.UserSectionStatus{}, nil
	}

	// 2. 查询user_section_status表中对应的学习状态
	var statuses []model.UserSectionStatus
	err = r.DB.Where("user_id = ? AND section_id IN ?", userID, sectionIDs).Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}
