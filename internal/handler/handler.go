package handler

import (
	"awesomeProject/internal/handler/app"
	"awesomeProject/internal/handler/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {

	// 注册中间件
	r.Use(middleware.CORS())

	r.GET("/startContainer", app.StartContainer)
	r.GET("/check", app.CheckAnswer)

	// 课程相关路由
	courseGroup := r.Group("/courses")
	{
		courseGroup.GET("/:id", app.GetCourseByIDHandler)
		courseGroup.GET("", app.GetAllCoursesHandler)
	}

	// 章节相关路由
	chapterGroup := r.Group("/chapters")
	{
		chapterGroup.POST("", app.CreateChapterHandler)
	}
}
