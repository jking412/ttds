package app

import (
	"awesomeProject/internal/usecase"
)

var (
	courseService usecase.CourseService
)

// initCourseService 初始化课程服务
func initCourseService() {
	courseService = usecase.NewCourseService()
}
