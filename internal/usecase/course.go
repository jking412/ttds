package usecase

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/oss"
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
	GetCourseReferencesDownloadURL(referenceID uint) (string, error)
	GetCourseStatus(userID, courseID uint) ([]model.UserSectionStatus, error)
}

// CourseServiceImpl 课程服务实现
type CourseServiceImpl struct {
	CourseRepository repository.CourseRepository
	cache            db.Cache
	ossManager       oss.Manager
}

func NewCourseService() CourseService {
	courseSyncOnce.Do(func() {
		courseServiceInstance = &CourseServiceImpl{
			CourseRepository: repository.NewCourseRepository(db.DB),
			cache:            db.NewCache(),
			ossManager:       oss.NewOssClient(),
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
	coursesStr, err := s.cache.Get(context.Background(), "courses")
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
		if err = s.cache.Set(context.Background(), "courses", coursesJSON, time.Minute); err != nil {
			return nil, err
		}
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

func (s *CourseServiceImpl) GetCourseReferencesDownloadURL(referenceID uint) (string, error) {
	reference, err := s.CourseRepository.GetCourseReferenceByID(referenceID)
	if err != nil {
		return "", err
	}

	defaultExpiredTime := 10 * time.Minute

	// 生成预签名URL
	downloadURL, err := s.ossManager.GetObjectUrl(reference.Title, int64(defaultExpiredTime/time.Second))
	if err != nil {
		return "", err
	}

	return downloadURL, nil
}
