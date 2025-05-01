package usecase

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"awesomeProject/pkg/db"
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

var (
	courseServiceInstance CourseService
	courseSyncOnce        sync.Once

	_ CourseService = (*CourseServiceImpl)(nil)
)

// CourseService 课程服务接口
type CourseService interface {
	GetAllCourses() ([]model.Course, error)
	GetCourseByID(id uint) (*model.Course, error)
	GetCourseReferences(courseID uint) ([]model.CourseReference, error)
	GetCourseStatus(userID, courseID uint) ([]model.UserSectionStatus, error)
}

// CourseServiceImpl 课程服务实现
type CourseServiceImpl struct {
	CourseRepository repository.CourseRepository
}

func NewCourseService() CourseService {
	courseSyncOnce.Do(func() {
		courseServiceInstance = &CourseServiceImpl{
			CourseRepository: repository.NewCourseRepository(db.DB),
		}
	})
	return courseServiceInstance
}

// GetCourseStatus 获取课程学习状态
func (s *CourseServiceImpl) GetCourseStatus(userID, courseID uint) ([]model.UserSectionStatus, error) {
	return s.CourseRepository.GetCourseStatusByCourseID(userID, courseID)
}

// GetAllCourses 获取所有课程
func (s *CourseServiceImpl) GetAllCourses() ([]model.Course, error) {
	// 尝试从缓存获取
	var courses []model.Course
	coursesStr, err := db.Cache.Get(context.Background(), "courses").Result()
	if err == nil {
		err := json.Unmarshal([]byte(coursesStr), &courses)
		if err == nil {
			return courses, nil
		}
	}

	// 从数据库获取
	courses, err = s.CourseRepository.GetAllCourses()
	if err != nil {
		return nil, err
	}

	// 设置缓存
	coursesJSON, err := json.Marshal(courses)
	if err == nil {
		db.Cache.Set(context.Background(), "courses", coursesJSON, time.Minute)
	}

	return courses, nil
}

// GetCourseByID 根据ID获取课程
func (s *CourseServiceImpl) GetCourseByID(id uint) (*model.Course, error) {
	course, err := s.CourseRepository.GetCourseByID(id)
	if err != nil {
		return nil, errors.New("课程不存在")
	}
	return course, nil
}

// GetCourseReferences 获取课程参考资料
func (s *CourseServiceImpl) GetCourseReferences(courseID uint) ([]model.CourseReference, error) {
	references, err := s.CourseRepository.GetCourseReferencesByCourseID(courseID)
	if err != nil {
		return nil, err
	}
	return references, nil
}
