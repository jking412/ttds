package usecase

import (
	"awesomeProject/internal/model"
	"context"
	"time"

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
