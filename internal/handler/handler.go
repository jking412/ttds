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

	r.GET("/hello", helloHandler)

	r.GET("/startContainer", app.StartContainer)
	r.GET("/check", app.CheckAnswer)

	// 用户认证相关路由
	r.POST("/api/v1/register", app.RegisterHandler)      // 用户注册
	r.POST("/api/v1/login", app.LoginHandler)          // 用户登录
	r.POST("/api/v1/logout", app.LogoutHandler)        // 用户注销
	
	// 需要认证的路由组
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.GET("/user", app.GetCurrentUserHandler)  // 获取当前用户信息
	}

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

func helloHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
