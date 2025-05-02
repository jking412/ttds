package usecase

import (
	"awesomeProject/internal/model"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCourseRepository struct {
	mock.Mock
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCourseRepository) GetAllCourses() ([]model.Course, error) {
	args := m.Called()
	return args.Get(0).([]model.Course), args.Error(1)
}

func (m *MockCourseRepository) GetCourseByID(id uint) (*model.Course, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Course), args.Error(1)
}

func (m *MockCourseRepository) GetCourseReferencesByCourseID(courseID uint) ([]model.CourseReference, error) {
	args := m.Called(courseID)
	return args.Get(0).([]model.CourseReference), args.Error(1)
}

func (m *MockCourseRepository) GetCourseStatusByCourseID(userID, courseID uint) ([]model.UserSectionStatus, error) {
	args := m.Called(userID, courseID)
	return args.Get(0).([]model.UserSectionStatus), args.Error(1)
}

// TODO: 所有测试存在问题
func TestGetAllCoursesWithCache(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	mockCache := new(MockCache)
	service := &CourseServiceImpl{
		CourseRepository: mockRepo,
		cache:            mockCache,
	}

	// 测试缓存命中
	t.Run("Cache Hit", func(t *testing.T) {
		cachedCourses := []model.Course{{Title: "Cached Course"}}
		cachedData, _ := json.Marshal(cachedCourses)
		mockCache.On("Get", mock.Anything, "courses").Return(string(cachedData), nil)

		result, err := service.GetAllCourses()
		assert.NoError(t, err)
		assert.Equal(t, "Cached Course", result[0].Title)
		mockRepo.AssertNotCalled(t, "GetAllCourses")
	})

	// 测试缓存未命中
	t.Run("Cache Miss", func(t *testing.T) {
		dbCourses := []model.Course{{Title: "DB Course"}}
		mockCache.On("Get", mock.Anything, "courses").Return("", errors.New("not found"))
		mockCache.On("Set", mock.Anything, "courses", mock.Anything, time.Minute).Return(nil)
		mockRepo.On("GetAllCourses").Return(dbCourses, nil)

		result, err := service.GetAllCourses()
		assert.NoError(t, err)
		assert.Equal(t, "DB Course", result[0].Title)
		mockRepo.AssertCalled(t, "GetAllCourses")
		mockCache.AssertCalled(t, "Set", mock.Anything, "courses", mock.Anything, time.Minute)
	})

	// 测试缓存错误
	t.Run("Cache Error", func(t *testing.T) {
		mockCache.On("Get", mock.Anything, "courses").Return("", errors.New("cache error"))
		mockRepo.On("GetAllCourses").Return(nil, errors.New("db error"))

		_, err := service.GetAllCourses()
		assert.Error(t, err)
	})
}
