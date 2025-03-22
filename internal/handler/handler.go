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
		courseGroup.POST("", app.CreateCourseHandler)
		courseGroup.GET("/:id", app.GetCourseByIDHandler)
		courseGroup.GET("", app.GetAllCoursesHandler)
	}

	// 章节相关路由
	chapterGroup := r.Group("/chapters")
	{
		chapterGroup.POST("", app.CreateChapterHandler)
		chapterGroup.GET("/:id", app.GetChapterByIDHandler)
		chapterGroup.GET("/course/:courseID", app.GetChaptersByCourseIDHandler)
	}

	// 小节相关路由
	sectionGroup := r.Group("/sections")
	{
		sectionGroup.POST("", app.CreateSectionHandler)
		sectionGroup.GET("/:id", app.GetSectionByIDHandler)
		sectionGroup.GET("/chapter/:chapterID", app.GetSectionsByChapterIDHandler)
	}

	// 用户相关路由
	userGroup := r.Group("/users")
	{
		userGroup.POST("/", app.CreateUserHandler)
		userGroup.GET("/:id", app.GetUserByIDHandler)
	}

	// 用户小节状态相关路由
	userSectionStatusGroup := r.Group("/user-section-status")
	{
		userSectionStatusGroup.POST("", app.CreateUserSectionStatusHandler)
		userSectionStatusGroup.GET("/user/:userID/section/:sectionID", app.GetUserSectionStatusByUserAndSectionIDHandler)
		userSectionStatusGroup.PUT("/user/:userID/section/:sectionID", app.UpdateUserSectionStatusHandler)
	}

	// 课程参考书籍相关路由
	courseReferenceBookGroup := r.Group("/course-reference-books")
	{
		courseReferenceBookGroup.POST("", app.CreateCourseReferenceBookHandler)
		courseReferenceBookGroup.GET("/course/:courseID", app.GetCourseReferenceBooksByCourseIDHandler)
	}
}
