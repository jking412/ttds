package repository

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/db"
)

// CreateCourse 创建一个新的课程
func CreateCourse(course *model.Course) error {
	return db.DB.Create(course).Error
}

// GetCourseByID 根据课程 ID 获取课程信息
func GetCourseByID(id uint) (*model.Course, error) {
	var course model.Course
	result := db.DB.Preload("Chapters").Preload("ReferenceBooks").First(&course, id)
	return &course, result.Error
}

// GetAllCourses 获取所有课程信息
func GetAllCourses() ([]model.Course, error) {
	var courses []model.Course
	result := db.DB.Preload("Chapters").Preload("ReferenceBooks").Find(&courses)
	return courses, result.Error
}

// CreateCourseReferenceBook 创建课程参考书籍记录
func CreateCourseReferenceBook(book *model.CourseReferenceBook) error {
	return db.DB.Create(book).Error
}

// GetCourseReferenceBooksByCourseID 根据课程 ID 获取所有参考书籍信息
func GetCourseReferenceBooksByCourseID(courseID uint) ([]model.CourseReferenceBook, error) {
	var books []model.CourseReferenceBook
	result := db.DB.Where("course_id = ?", courseID).Find(&books)
	return books, result.Error
}
